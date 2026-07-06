package sdr

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portsvc "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/service"
)

// Adapter wraps the rtl_fm | multimon-ng pipeline to decode SAME/EAS messages.
type Adapter struct {
	devicePath string
	frequency  int64
	cancel     context.CancelFunc
}

func New(devicePath string, frequency int64) *Adapter {
	return &Adapter{
		devicePath: devicePath,
		frequency:  frequency,
	}
}

var _ portsvc.SAMEDecoder = (*Adapter)(nil)

// Start launches the rtl_fm → multimon-ng pipeline in the background.
// Decoded SAME alerts are sent to ch. The caller drains ch and calls Stop.
func (a *Adapter) Start(ctx context.Context, ch chan<- entity.SAMEAlert) error {
	pipeCtx, cancel := context.WithCancel(ctx)
	a.cancel = cancel

	// rtl_fm demodulates the NOAA Weather Radio FM signal to raw audio.
	rtlCmd := exec.CommandContext(pipeCtx,
		"rtl_fm",
		"-f", fmt.Sprintf("%d", a.frequency),
		"-M", "fm",
		"-s", "22050",
		"-g", "100",
		"-",
	)

	// multimon-ng decodes the EAS/SAME protocol from the audio stream.
	mmCmd := exec.CommandContext(pipeCtx,
		"multimon-ng",
		"-a", "EAS",
		"-t", "raw",
		"-",
	)

	pipe, err := rtlCmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("creating rtl_fm stdout pipe: %w", err)
	}
	mmCmd.Stdin = pipe

	mmOut, err := mmCmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("creating multimon-ng stdout pipe: %w", err)
	}

	if err := rtlCmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("starting rtl_fm: %w", err)
	}
	if err := mmCmd.Start(); err != nil {
		cancel()
		if killErr := rtlCmd.Process.Kill(); killErr != nil {
			slog.Error("failed to kill rtl_fm after multimon-ng start error", "error", killErr)
		}
		return fmt.Errorf("starting multimon-ng: %w", err)
	}

	go func() {
		defer cancel()
		if err := a.readOutput(mmOut, ch); err != nil && pipeCtx.Err() == nil {
			slog.Error("SDR pipeline error", "error", err)
		}
		if err := rtlCmd.Wait(); err != nil && pipeCtx.Err() == nil {
			slog.Error("rtl_fm exited with error", "error", err)
		}
		if err := mmCmd.Wait(); err != nil && pipeCtx.Err() == nil {
			slog.Error("multimon-ng exited with error", "error", err)
		}
	}()

	slog.Info("SDR pipeline started",
		"frequency_hz", a.frequency,
		"device", a.devicePath,
	)
	return nil
}

// Stop terminates the SDR pipeline.
func (a *Adapter) Stop() {
	if a.cancel != nil {
		a.cancel()
	}
}

// readOutput scans multimon-ng output lines and parses EAS SAME messages.
func (a *Adapter) readOutput(r io.Reader, ch chan<- entity.SAMEAlert) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "EAS: ZCZC-") {
			continue
		}
		raw := strings.TrimPrefix(line, "EAS: ")
		alert, err := ParseSAME(raw)
		if err != nil {
			slog.Warn("failed to parse SAME message", "raw", raw, "error", err)
			continue
		}
		slog.Info("decoded SAME alert",
			"event_code", alert.EventCode,
			"originator", alert.Originator,
			"fips_count", len(alert.FIPSCodes),
		)
		ch <- alert
	}
	return scanner.Err()
}

// ParseSAME parses a raw SAME header string into a SAMEAlert.
// Input format: ZCZC-ORG-EEE-PSSCCC+TTTT-JJJHHMM-CCCCCCCC-
// Example: ZCZC-WXR-TOR-037121+0030-1820218-KRAH/NWS-
func ParseSAME(raw string) (entity.SAMEAlert, error) {
	// Strip ZCZC- prefix and trailing -
	s := strings.TrimPrefix(raw, "ZCZC-")
	s = strings.TrimSuffix(s, "-")

	parts := strings.Split(s, "-")
	if len(parts) < 5 {
		return entity.SAMEAlert{}, fmt.Errorf("too few fields: %d", len(parts))
	}

	alert := entity.SAMEAlert{
		Originator: parts[0],
		EventCode:  parts[1],
		RawMessage: raw,
	}

	// Fields starting at index 2 are FIPS codes. The last FIPS field has the purge
	// time embedded as "+TTTT" (e.g. "037121+0030"). Split on "+" to extract both.
	idx := 2
	for idx < len(parts) {
		p := parts[idx]
		if plus := strings.Index(p, "+"); plus != -1 {
			// This is the last FIPS+purge field.
			alert.FIPSCodes = append(alert.FIPSCodes, p[:plus])
			alert.PurgeTime = p[plus+1:]
			idx++
			break
		}
		alert.FIPSCodes = append(alert.FIPSCodes, p)
		idx++
	}

	if alert.PurgeTime == "" {
		return entity.SAMEAlert{}, fmt.Errorf("missing purge time field")
	}

	if idx >= len(parts) {
		return entity.SAMEAlert{}, fmt.Errorf("missing issue time field")
	}
	alert.IssueTime = parts[idx]
	idx++

	if idx < len(parts) {
		alert.CallSign = parts[idx]
	}

	if len(alert.FIPSCodes) == 0 {
		return entity.SAMEAlert{}, fmt.Errorf("no FIPS codes found")
	}

	return alert, nil
}
