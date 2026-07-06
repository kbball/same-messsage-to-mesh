package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

// stubFilterRepo returns a fixed AlertFilter.
type stubFilterRepo struct {
	filter entity.AlertFilter
}

func (s *stubFilterRepo) Get(_ context.Context) (entity.AlertFilter, error) { return s.filter, nil }
func (s *stubFilterRepo) Update(_ context.Context, _ entity.AlertFilter) error { return nil }

// stubAlertRepo records created alerts.
type stubAlertRepo struct {
	created []entity.SAMEAlert
}

func (s *stubAlertRepo) Create(_ context.Context, a entity.SAMEAlert) (entity.SAMEAlert, error) {
	a.ID = int64(len(s.created) + 1)
	a.ReceivedAt = time.Now()
	s.created = append(s.created, a)
	return a, nil
}
func (s *stubAlertRepo) List(_ context.Context, _ int) ([]entity.SAMEAlert, error) {
	return s.created, nil
}
func (s *stubAlertRepo) MarkPublished(_ context.Context, _ int64) error { return nil }

// stubFIPSRepo returns a fixed FIPSCode.
type stubFIPSRepo struct{}

func (s *stubFIPSRepo) ListStates(_ context.Context) ([]entity.FIPSCode, error) { return nil, nil }
func (s *stubFIPSRepo) ListByState(_ context.Context, _ string) ([]entity.FIPSCode, error) {
	return nil, nil
}
func (s *stubFIPSRepo) GetByFIPS(_ context.Context, fips string) (entity.FIPSCode, error) {
	return entity.FIPSCode{
		StateCode: fips[:2], CountyCode: fips[2:],
		StateName: "Georgia", CountyName: "Fulton County",
	}, nil
}
func (s *stubFIPSRepo) Upsert(_ context.Context, _ []entity.FIPSCode) error { return nil }
func (s *stubFIPSRepo) Count(_ context.Context) (int, error)                { return 0, nil }

// stubECRepo returns a fixed EventCode.
type stubECRepo struct{}

func (s *stubECRepo) List(_ context.Context) ([]entity.EventCode, error) { return nil, nil }
func (s *stubECRepo) Get(_ context.Context, code string) (entity.EventCode, error) {
	return entity.EventCode{Code: code, Description: "Tornado Warning"}, nil
}
func (s *stubECRepo) Upsert(_ context.Context, _ []entity.EventCode) error { return nil }

func TestAlertService_Handle_NoFilter(t *testing.T) {
	alertRepo := &stubAlertRepo{}
	filterRepo := &stubFilterRepo{filter: entity.AlertFilter{}}
	svc := NewAlertService(alertRepo, filterRepo, &stubFIPSRepo{}, &stubECRepo{}, nil)

	alert := entity.SAMEAlert{
		Originator: "WXR",
		EventCode:  "TOR",
		FIPSCodes:  []string{"013121"},
		PurgeTime:  "0030",
		IssueTime:  "1820218",
		CallSign:   "KFFC/NWS",
		RawMessage: "ZCZC-WXR-TOR-013121+0030-1820218-KFFC/NWS-",
	}

	saved, err := svc.Handle(context.Background(), alert)
	require.NoError(t, err)
	assert.Equal(t, int64(1), saved.ID)
	assert.Equal(t, "TOR", saved.EventCode)
	assert.Len(t, alertRepo.created, 1)
}

func TestAlertService_Handle_FilteredOut(t *testing.T) {
	alertRepo := &stubAlertRepo{}
	filterRepo := &stubFilterRepo{
		filter: entity.AlertFilter{
			StateCodes: []string{"12"}, // Florida only
		},
	}
	svc := NewAlertService(alertRepo, filterRepo, &stubFIPSRepo{}, &stubECRepo{}, nil)

	alert := entity.SAMEAlert{
		EventCode: "TOR",
		FIPSCodes: []string{"013121"}, // Georgia (state 13)
	}

	saved, err := svc.Handle(context.Background(), alert)
	require.NoError(t, err)
	assert.Equal(t, int64(0), saved.ID)
	assert.Empty(t, alertRepo.created)
}

func TestAlertService_Handle_EventCodeFilter(t *testing.T) {
	alertRepo := &stubAlertRepo{}
	filterRepo := &stubFilterRepo{
		filter: entity.AlertFilter{
			EventCodes: []string{"TOR"},
		},
	}
	svc := NewAlertService(alertRepo, filterRepo, &stubFIPSRepo{}, &stubECRepo{}, nil)

	rwtAlert := entity.SAMEAlert{EventCode: "RWT", FIPSCodes: []string{"013121"}}
	saved, err := svc.Handle(context.Background(), rwtAlert)
	require.NoError(t, err)
	assert.Equal(t, int64(0), saved.ID)

	torAlert := entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}}
	saved, err = svc.Handle(context.Background(), torAlert)
	require.NoError(t, err)
	assert.Equal(t, int64(1), saved.ID)
}

func TestAlertService_MatchesFilter(t *testing.T) {
	svc := &AlertService{}

	tests := []struct {
		name   string
		alert  entity.SAMEAlert
		filter entity.AlertFilter
		want   bool
	}{
		{
			name:   "empty filter matches all",
			alert:  entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{},
			want:   true,
		},
		{
			name:  "state match",
			alert: entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{StateCodes: []string{"13"}},
			want:  true,
		},
		{
			name:  "state no match",
			alert: entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{StateCodes: []string{"12"}},
			want:  false,
		},
		{
			name:  "fips match",
			alert: entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{FIPSCodes: []string{"13121"}},
			want:  true,
		},
		{
			name:  "event code match",
			alert: entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{EventCodes: []string{"TOR"}},
			want:  true,
		},
		{
			name:  "event code no match",
			alert: entity.SAMEAlert{EventCode: "RWT", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{EventCodes: []string{"TOR"}},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.matchesFilter(tt.alert, tt.filter)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatMessage(t *testing.T) {
	svc := &AlertService{
		fipsRepo: &stubFIPSRepo{},
		ecRepo:   &stubECRepo{},
	}

	alert := entity.SAMEAlert{
		EventCode: "TOR",
		FIPSCodes: []string{"013121"},
		CallSign:  "KFFC/NWS",
	}

	msg := svc.formatMessage(context.Background(), alert)
	assert.Contains(t, msg, "[TOR]")
	assert.Contains(t, msg, "Tornado Warning")
	assert.Contains(t, msg, "Fulton County")
	assert.Contains(t, msg, "KFFC/NWS")
}
