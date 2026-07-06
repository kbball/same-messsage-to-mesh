package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"
	_ "github.com/lib/pq"

	httphandler "github.com/kbball/same-message-to-mesh/backend/internal/adapter/http/handler"
	mqttadapter "github.com/kbball/same-message-to-mesh/backend/internal/adapter/mqtt"
	"github.com/kbball/same-message-to-mesh/backend/internal/adapter/noaa"
	"github.com/kbball/same-message-to-mesh/backend/internal/adapter/repository"
	"github.com/kbball/same-message-to-mesh/backend/internal/adapter/sdr"
	sseadapter "github.com/kbball/same-message-to-mesh/backend/internal/adapter/sse"
	"github.com/kbball/same-message-to-mesh/backend/internal/application/service"
	"github.com/kbball/same-message-to-mesh/backend/internal/config"
	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	setupLogger(cfg)

	db, err := connectDB(cfg.DB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database connection", "error", err)
		}
	}()

	if err := runMigrations(db); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Repositories
	alertRepo := repository.NewAlertRepo(db)
	ecRepo := repository.NewEventCodeRepo(db)
	fipsRepo := repository.NewFIPSRepo(db)
	filterRepo := repository.NewFilterRepo(db)
	sdrCfgRepo := repository.NewSDRConfigRepo(db)
	mqttCfgRepo := repository.NewMQTTConfigRepo(db)

	// SSE broker
	broker := sseadapter.NewBroker()

	// NOAA fetcher
	noaaFetcher := noaa.New()

	// MQTT publisher — start if enabled in DB config
	var mqttPub *mqttadapter.Publisher
	mqttCfg, err := mqttCfgRepo.Get(context.Background())
	if err != nil {
		slog.Warn("could not load MQTT config from DB", "error", err)
	} else if mqttCfg.Enabled {
		if p, err := mqttadapter.New(mqttCfg); err != nil {
			slog.Warn("MQTT publisher could not start", "error", err)
		} else {
			mqttPub = p
			defer mqttPub.Close()
		}
	}

	// Application services
	alertSvc := service.NewAlertService(alertRepo, filterRepo, fipsRepo, ecRepo, mqttPub)
	filterSvc := service.NewFilterService(filterRepo, sdrCfgRepo, mqttCfgRepo)
	refDataSvc := service.NewReferenceDataService(fipsRepo, ecRepo, noaaFetcher)

	// SDR pipeline — load persisted config from DB
	sdrCfg, err := sdrCfgRepo.Get(context.Background())
	if err != nil {
		slog.Warn("could not load SDR config from DB, using defaults from env",
			"device", cfg.SDR.DevicePath,
			"frequency", cfg.SDR.Frequency,
		)
		sdrCfg = entity.SDRDeviceConfig{
			DevicePath: cfg.SDR.DevicePath,
			Frequency:  cfg.SDR.Frequency,
		}
	}

	alertCh := make(chan entity.SAMEAlert, 32)
	sdrAdapter := sdr.New(sdrCfg.DevicePath, sdrCfg.Frequency)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := sdrAdapter.Start(ctx, alertCh); err != nil {
		slog.Warn("SDR pipeline could not start — alerts will not be decoded", "error", err)
	} else {
		defer sdrAdapter.Stop()
		go func() {
			for alert := range alertCh {
				saved, err := alertSvc.Handle(context.Background(), alert)
				if err != nil {
					slog.Error("failed to handle alert", "error", err)
					continue
				}
				if saved.ID > 0 {
					broker.Publish("alert", saved)
				}
			}
		}()
	}

	reconnectMQTT := func(cfg entity.MQTTConfig) error {
		if mqttPub != nil {
			mqttPub.Close()
		}
		if !cfg.Enabled {
			mqttPub = nil
			alertSvc.SetPublisher(nil)
			return nil
		}
		p, err := mqttadapter.New(cfg)
		if err != nil {
			mqttPub = nil
			alertSvc.SetPublisher(nil)
			return err
		}
		mqttPub = p
		alertSvc.SetPublisher(p)
		return nil
	}

	h := httphandler.New(alertSvc, filterSvc, refDataSvc, broker).WithMQTT(mqttPub, reconnectMQTT)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.Handle("GET /api/stream", broker)
	h.Register(mux)
	mux.Handle("/", serveSPA("frontend/dist"))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      httphandler.LoggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server started", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}

func setupLogger(cfg *config.Config) {
	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}

func connectDB(cfg config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("opening connection: %w", err)
	}
	for attempt := range 10 {
		if err = db.Ping(); err == nil {
			slog.Info("database connected")
			return db, nil
		}
		slog.Info("waiting for database", "attempt", attempt+1, "error", err)
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("database not ready after 10 attempts: %w", err)
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(repository.MigrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("setting goose dialect: %w", err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}
	slog.Info("migrations applied")
	return nil
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func serveSPA(dir string) http.Handler {
	root := os.DirFS(dir)
	fileServer := http.FileServerFS(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Stat(root, r.URL.Path[1:])
		if err != nil {
			http.ServeFileFS(w, r, root, "index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
