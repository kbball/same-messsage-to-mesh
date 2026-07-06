package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

// stubNOAAFetcher returns configurable FIPS/EventCode slices or errors.
type stubNOAAFetcher struct {
	fips      []entity.FIPSCode
	events    []entity.EventCode
	fipsErr   error
	eventsErr error
}

func (f *stubNOAAFetcher) FetchFIPS(_ context.Context) ([]entity.FIPSCode, error) {
	return f.fips, f.fipsErr
}
func (f *stubNOAAFetcher) FetchEventCodes(_ context.Context) ([]entity.EventCode, error) {
	return f.events, f.eventsErr
}

// errFIPSRepo returns an error from Upsert.
type errFIPSRepo struct{ stubFIPSRepo }

func (e *errFIPSRepo) Upsert(_ context.Context, _ []entity.FIPSCode) error {
	return errors.New("db write error")
}

// errECRepo returns an error from Upsert.
type errECRepo struct{ stubECRepo }

func (e *errECRepo) Upsert(_ context.Context, _ []entity.EventCode) error {
	return errors.New("db write error")
}

func TestReferenceDataService_RefreshFIPS(t *testing.T) {
	fips := []entity.FIPSCode{
		{StateCode: "13", CountyCode: "121", StateName: "Georgia", CountyName: "Fulton County"},
		{StateCode: "12", CountyCode: "086", StateName: "Florida", CountyName: "Miami-Dade County"},
	}
	svc := NewReferenceDataService(&stubFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{fips: fips})
	n, err := svc.RefreshFIPS(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, n)
}

func TestReferenceDataService_RefreshFIPS_FetchError(t *testing.T) {
	svc := NewReferenceDataService(&stubFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{fipsErr: errors.New("timeout")})
	_, err := svc.RefreshFIPS(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fetching FIPS codes")
}

func TestReferenceDataService_RefreshFIPS_UpsertError(t *testing.T) {
	fips := []entity.FIPSCode{{StateCode: "13", CountyCode: "121"}}
	svc := NewReferenceDataService(&errFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{fips: fips})
	_, err := svc.RefreshFIPS(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "upserting FIPS codes")
}

func TestReferenceDataService_RefreshEventCodes(t *testing.T) {
	codes := []entity.EventCode{
		{Code: "TOR", Description: "Tornado Warning"},
		{Code: "RWT", Description: "Required Weekly Test"},
	}
	svc := NewReferenceDataService(&stubFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{events: codes})
	n, err := svc.RefreshEventCodes(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, n)
}

func TestReferenceDataService_RefreshEventCodes_FetchError(t *testing.T) {
	svc := NewReferenceDataService(&stubFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{eventsErr: errors.New("network")})
	_, err := svc.RefreshEventCodes(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fetching event codes")
}

func TestReferenceDataService_RefreshEventCodes_UpsertError(t *testing.T) {
	codes := []entity.EventCode{{Code: "TOR"}}
	svc := NewReferenceDataService(&stubFIPSRepo{}, &errECRepo{}, &stubNOAAFetcher{events: codes})
	_, err := svc.RefreshEventCodes(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "upserting event codes")
}

func TestReferenceDataService_ListFIPSByState(t *testing.T) {
	svc := NewReferenceDataService(&stubFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{})
	_, err := svc.ListFIPSByState(context.Background(), "13")
	require.NoError(t, err)
}

func TestReferenceDataService_ListStates(t *testing.T) {
	svc := NewReferenceDataService(&stubFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{})
	_, err := svc.ListStates(context.Background())
	require.NoError(t, err)
}

func TestReferenceDataService_ListEventCodes(t *testing.T) {
	svc := NewReferenceDataService(&stubFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{})
	_, err := svc.ListEventCodes(context.Background())
	require.NoError(t, err)
}

func TestReferenceDataService_FIPSCount(t *testing.T) {
	svc := NewReferenceDataService(&stubFIPSRepo{}, &stubECRepo{}, &stubNOAAFetcher{})
	n, err := svc.FIPSCount(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, n)
}
