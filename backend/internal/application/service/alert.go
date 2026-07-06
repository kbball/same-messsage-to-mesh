package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portrepo "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/repository"
	portsvc "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/service"
)

// AlertService receives decoded SAME alerts, applies filter logic, persists matches,
// and optionally publishes them.
type AlertService struct {
	alertRepo portrepo.AlertRepository
	filterRepo portrepo.FilterRepository
	fipsRepo   portrepo.FIPSRepository
	ecRepo     portrepo.EventCodeRepository
	publisher  portsvc.AlertPublisher // may be nil when MQTT disabled
}

func NewAlertService(
	alertRepo portrepo.AlertRepository,
	filterRepo portrepo.FilterRepository,
	fipsRepo portrepo.FIPSRepository,
	ecRepo portrepo.EventCodeRepository,
	publisher portsvc.AlertPublisher,
) *AlertService {
	return &AlertService{
		alertRepo:  alertRepo,
		filterRepo: filterRepo,
		fipsRepo:   fipsRepo,
		ecRepo:     ecRepo,
		publisher:  publisher,
	}
}

// SetPublisher swaps the publisher at runtime (used when MQTT config is updated via UI).
func (s *AlertService) SetPublisher(pub portsvc.AlertPublisher) {
	s.publisher = pub
}

// Handle processes a decoded SAME alert: checks the filter, persists it, and publishes if enabled.
func (s *AlertService) Handle(ctx context.Context, alert entity.SAMEAlert) (entity.SAMEAlert, error) {
	filter, err := s.filterRepo.Get(ctx)
	if err != nil {
		return entity.SAMEAlert{}, fmt.Errorf("getting filter: %w", err)
	}

	if !s.matchesFilter(alert, filter) {
		slog.Info("alert filtered out",
			"event_code", alert.EventCode,
			"fips_codes", alert.FIPSCodes,
		)
		return entity.SAMEAlert{}, nil
	}

	saved, err := s.alertRepo.Create(ctx, alert)
	if err != nil {
		return entity.SAMEAlert{}, fmt.Errorf("saving alert: %w", err)
	}

	if s.publisher != nil {
		msg := s.formatMessage(ctx, saved)
		if err := s.publisher.Publish(ctx, saved, msg); err != nil {
			slog.Warn("failed to publish alert", "id", saved.ID, "error", err)
		} else {
			_ = s.alertRepo.MarkPublished(ctx, saved.ID)
			saved.Published = true
		}
	}

	return saved, nil
}

// List returns recent alerts up to limit.
func (s *AlertService) List(ctx context.Context, limit int) ([]entity.SAMEAlert, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.alertRepo.List(ctx, limit)
}

// matchesFilter returns true if the alert should be acted on given the current filter config.
// Empty filter dimensions mean "match all".
func (s *AlertService) matchesFilter(alert entity.SAMEAlert, filter entity.AlertFilter) bool {
	// Event code filter
	if len(filter.EventCodes) > 0 && !contains(filter.EventCodes, alert.EventCode) {
		return false
	}

	// Geographic filter: check if any of the alert's FIPS codes overlap with the filter.
	// A FIPS code in SAME is "PSSCCC" (1 digit prefix + 2 state + 3 county).
	// Strip the prefix digit to get the 5-character "SSCCC" form.
	if len(filter.StateCodes) == 0 && len(filter.FIPSCodes) == 0 {
		return true
	}

	for _, rawFIPS := range alert.FIPSCodes {
		fips := stripFIPSPrefix(rawFIPS) // normalize to "SSCCC"
		stateCode := fips[:2]

		if len(filter.FIPSCodes) > 0 && contains(filter.FIPSCodes, fips) {
			return true
		}
		if len(filter.StateCodes) > 0 && contains(filter.StateCodes, stateCode) &&
			(len(filter.FIPSCodes) == 0 || isStateWideCode(fips)) {
			return true
		}
	}
	return false
}

// formatMessage builds the plain-text MQTT message for mesh broadcast.
func (s *AlertService) formatMessage(ctx context.Context, alert entity.SAMEAlert) string {
	description := alert.EventCode
	if ec, err := s.ecRepo.Get(ctx, alert.EventCode); err == nil {
		description = ec.Description
	}

	var locations []string
	for _, rawFIPS := range alert.FIPSCodes {
		fips := stripFIPSPrefix(rawFIPS)
		if f, err := s.fipsRepo.GetByFIPS(ctx, fips); err == nil {
			locations = append(locations, fmt.Sprintf("%s %s", f.CountyName, f.StateName))
		} else {
			locations = append(locations, fips)
		}
	}

	msg := fmt.Sprintf("[%s] %s", alert.EventCode, description)
	if len(locations) > 0 {
		msg += " - " + strings.Join(locations, ", ")
	}
	if alert.CallSign != "" {
		msg += fmt.Sprintf(" (%s)", alert.CallSign)
	}
	return msg
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// stripFIPSPrefix removes the leading subdivision digit from a SAME FIPS code.
// SAME encodes areas as "PSSCCC" where P=0 means whole county.
func stripFIPSPrefix(fips string) string {
	if len(fips) == 6 {
		return fips[1:] // "PSSCCC" → "SSCCC"
	}
	return fips
}

// isStateWideCode returns true if the county portion is "000" (whole state).
func isStateWideCode(fips string) bool {
	return len(fips) >= 5 && fips[2:] == "000"
}
