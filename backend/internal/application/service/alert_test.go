package service

import (
	"context"
	"errors"
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
func (s *stubFilterRepo) Update(_ context.Context, f entity.AlertFilter) error {
	s.filter = f
	return nil
}

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
			name:   "state match",
			alert:  entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{StateCodes: []string{"13"}},
			want:   true,
		},
		{
			name:   "state no match",
			alert:  entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{StateCodes: []string{"12"}},
			want:   false,
		},
		{
			name:   "fips match",
			alert:  entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{FIPSCodes: []string{"13121"}},
			want:   true,
		},
		{
			name:   "event code match",
			alert:  entity.SAMEAlert{EventCode: "TOR", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{EventCodes: []string{"TOR"}},
			want:   true,
		},
		{
			name:   "event code no match",
			alert:  entity.SAMEAlert{EventCode: "RWT", FIPSCodes: []string{"013121"}},
			filter: entity.AlertFilter{EventCodes: []string{"TOR"}},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.matchesFilter(tt.alert, tt.filter)
			assert.Equal(t, tt.want, got)
		})
	}
}

// stubPublisher records calls to Publish.
type stubPublisher struct {
	called  int
	lastMsg string
	err     error
}

func (p *stubPublisher) Publish(_ context.Context, _ entity.SAMEAlert, msg string) error {
	p.called++
	p.lastMsg = msg
	return p.err
}

func TestAlertService_List(t *testing.T) {
	alertRepo := &stubAlertRepo{}
	svc := NewAlertService(alertRepo, &stubFilterRepo{}, &stubFIPSRepo{}, &stubECRepo{}, nil)

	results, err := svc.List(context.Background(), 10)
	require.NoError(t, err)
	assert.Empty(t, results)

	alertRepo.created = []entity.SAMEAlert{{ID: 1, EventCode: "TOR"}}
	results, err = svc.List(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestAlertService_SetPublisher(t *testing.T) {
	svc := NewAlertService(&stubAlertRepo{}, &stubFilterRepo{}, &stubFIPSRepo{}, &stubECRepo{}, nil)
	assert.Nil(t, svc.publisher)
	pub := &stubPublisher{}
	svc.SetPublisher(pub)
	assert.Equal(t, pub, svc.publisher)
}

func TestAlertService_Handle_WithPublisher(t *testing.T) {
	alertRepo := &stubAlertRepo{}
	pub := &stubPublisher{}
	svc := NewAlertService(alertRepo, &stubFilterRepo{}, &stubFIPSRepo{}, &stubECRepo{}, pub)

	saved, err := svc.Handle(context.Background(), entity.SAMEAlert{
		EventCode: "TOR", FIPSCodes: []string{"013121"},
	})
	require.NoError(t, err)
	assert.True(t, saved.Published)
	assert.Equal(t, 1, pub.called)
	assert.Contains(t, pub.lastMsg, "[TOR]")
}

func TestAlertService_Handle_PublishError(t *testing.T) {
	alertRepo := &stubAlertRepo{}
	pub := &stubPublisher{err: errors.New("broker down")}
	svc := NewAlertService(alertRepo, &stubFilterRepo{}, &stubFIPSRepo{}, &stubECRepo{}, pub)

	saved, err := svc.Handle(context.Background(), entity.SAMEAlert{
		EventCode: "TOR", FIPSCodes: []string{"013121"},
	})
	require.NoError(t, err)
	assert.False(t, saved.Published) // publish failed, should not be marked
	assert.Len(t, alertRepo.created, 1)
}

func TestIsStateWideCode(t *testing.T) {
	assert.True(t, isStateWideCode("13000"))
	assert.False(t, isStateWideCode("13121"))
	assert.False(t, isStateWideCode("13"))
}

func TestStripFIPSPrefix(t *testing.T) {
	assert.Equal(t, "13121", stripFIPSPrefix("013121"))
	assert.Equal(t, "13121", stripFIPSPrefix("13121")) // no-op for 5-char
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
