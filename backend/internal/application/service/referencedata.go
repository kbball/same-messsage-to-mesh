package service

import (
	"context"
	"fmt"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portrepo "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/repository"
)

// NOAAFetcher is the port for fetching reference data from NOAA.
type NOAAFetcher interface {
	FetchFIPS(ctx context.Context) ([]entity.FIPSCode, error)
	FetchEventCodes(ctx context.Context) ([]entity.EventCode, error)
}

type ReferenceDataService struct {
	fipsRepo portrepo.FIPSRepository
	ecRepo   portrepo.EventCodeRepository
	fetcher  NOAAFetcher
}

func NewReferenceDataService(
	fipsRepo portrepo.FIPSRepository,
	ecRepo portrepo.EventCodeRepository,
	fetcher NOAAFetcher,
) *ReferenceDataService {
	return &ReferenceDataService{
		fipsRepo: fipsRepo,
		ecRepo:   ecRepo,
		fetcher:  fetcher,
	}
}

func (s *ReferenceDataService) RefreshFIPS(ctx context.Context) (int, error) {
	codes, err := s.fetcher.FetchFIPS(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetching FIPS codes: %w", err)
	}
	if err := s.fipsRepo.Upsert(ctx, codes); err != nil {
		return 0, fmt.Errorf("upserting FIPS codes: %w", err)
	}
	return len(codes), nil
}

func (s *ReferenceDataService) RefreshEventCodes(ctx context.Context) (int, error) {
	codes, err := s.fetcher.FetchEventCodes(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetching event codes: %w", err)
	}
	if err := s.ecRepo.Upsert(ctx, codes); err != nil {
		return 0, fmt.Errorf("upserting event codes: %w", err)
	}
	return len(codes), nil
}

func (s *ReferenceDataService) ListFIPSByState(ctx context.Context, stateCode string) ([]entity.FIPSCode, error) {
	return s.fipsRepo.ListByState(ctx, stateCode)
}

func (s *ReferenceDataService) ListStates(ctx context.Context) ([]entity.FIPSCode, error) {
	return s.fipsRepo.ListStates(ctx)
}

func (s *ReferenceDataService) ListEventCodes(ctx context.Context) ([]entity.EventCode, error) {
	return s.ecRepo.List(ctx)
}

func (s *ReferenceDataService) FIPSCount(ctx context.Context) (int, error) {
	return s.fipsRepo.Count(ctx)
}
