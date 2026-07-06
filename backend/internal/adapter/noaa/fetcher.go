package noaa

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

const (
	fipsURL       = "https://www.ncei.noaa.gov/erddap/convert/fipscounty.csv"
	eventCodesURL = "https://www.weather.gov/nwr/eventcodes"
)

// Fetcher retrieves SAME reference data from NOAA sources.
type Fetcher struct {
	client *http.Client
}

func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// FetchFIPS downloads county FIPS codes from NOAA NCEI.
// The ERDDAP CSV endpoint returns rows with FIPS code and county name.
func (f *Fetcher) FetchFIPS(ctx context.Context) ([]entity.FIPSCode, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fipsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching FIPS CSV: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from NOAA FIPS endpoint", resp.StatusCode)
	}

	return parseFIPSCSV(resp.Body)
}

// parseFIPSCSV parses the NOAA NCEI FIPS county CSV.
// Expected columns: FIPS (5-digit), Name (e.g. "Autauga County, Alabama")
func parseFIPSCSV(r io.Reader) ([]entity.FIPSCode, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	// County names contain commas (e.g. "Autauga County, Alabama"), so allow variable fields.
	reader.FieldsPerRecord = -1

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("reading CSV header: %w", err)
	}

	now := time.Now()
	var codes []entity.FIPSCode
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading CSV row: %w", err)
		}
		if len(record) < 2 {
			continue
		}

		fips := strings.TrimSpace(record[0])
		if len(fips) != 5 {
			continue
		}

		// Rejoin remaining fields to reconstruct "County Name, State Name"
		name := strings.TrimSpace(strings.Join(record[1:], ", "))
		countyName, stateName := splitCountyState(name)

		codes = append(codes, entity.FIPSCode{
			StateCode:  fips[:2],
			CountyCode: fips[2:],
			StateName:  stateName,
			CountyName: countyName,
			UpdatedAt:  now,
		})
	}

	if len(codes) == 0 {
		return nil, fmt.Errorf("no FIPS codes parsed from response")
	}
	return codes, nil
}

// splitCountyState splits "Autauga County, Alabama" into ("Autauga County", "Alabama").
func splitCountyState(name string) (county, state string) {
	idx := strings.LastIndex(name, ", ")
	if idx == -1 {
		return name, ""
	}
	return strings.TrimSpace(name[:idx]), strings.TrimSpace(name[idx+2:])
}

// FetchEventCodes returns the hardcoded list of SAME event codes.
// The NOAA event codes page (weather.gov/nwr/eventcodes) is HTML-only and
// changes rarely. We maintain the authoritative list here and let users
// trigger a refresh to pick up any additions.
func (f *Fetcher) FetchEventCodes(_ context.Context) ([]entity.EventCode, error) {
	now := time.Now()
	codes := []entity.EventCode{
		{Code: "BZW", Description: "Blizzard Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "CFW", Description: "Coastal Flood Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "DSW", Description: "Dust Storm Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "EWW", Description: "Extreme Wind Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "FFW", Description: "Flash Flood Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "FLW", Description: "Flood Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "FRW", Description: "Fire Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "FSW", Description: "Flash Freeze Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "FZW", Description: "Freeze Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "GLW", Description: "Gale Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "HFW", Description: "Hurricane Force Wind Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "HUW", Description: "Hurricane Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "HWW", Description: "High Wind Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "ISW", Description: "Ice Storm Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "LAW", Description: "Landslide Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "LEW", Description: "Law Enforcement Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "LSW", Description: "Land Slide Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "MAW", Description: "Marine Weather Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "NAW", Description: "National Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "NMW", Description: "Nuclear Power Plant Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "RHW", Description: "Radiological Hazard Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "SMW", Description: "Special Marine Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "SPW", Description: "Shelter in Place Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "SQW", Description: "Snow Squall Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "SSW", Description: "Storm Surge Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "SVW", Description: "Severe Thunderstorm Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "TOE", Description: "911 Telephone Outage Emergency", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "TOR", Description: "Tornado Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "TRW", Description: "Tropical Storm Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "TSW", Description: "Tsunami Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "TYW", Description: "Typhoon Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "VOW", Description: "Volcano Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "WSW", Description: "Winter Storm Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "AVW", Description: "Avalanche Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "DBW", Description: "Dam Break Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "EVW", Description: "Evacuation Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		{Code: "FWW", Description: "Fire Weather Warning", Category: "Warning", IsWarning: true, UpdatedAt: now},
		// Watches
		{Code: "CFA", Description: "Coastal Flood Watch", Category: "Watch", UpdatedAt: now},
		{Code: "FFA", Description: "Flash Flood Watch", Category: "Watch", UpdatedAt: now},
		{Code: "FLA", Description: "Flood Watch", Category: "Watch", UpdatedAt: now},
		{Code: "HUA", Description: "Hurricane Watch", Category: "Watch", UpdatedAt: now},
		{Code: "SVA", Description: "Severe Thunderstorm Watch", Category: "Watch", UpdatedAt: now},
		{Code: "TOA", Description: "Tornado Watch", Category: "Watch", UpdatedAt: now},
		{Code: "TRA", Description: "Tropical Storm Watch", Category: "Watch", UpdatedAt: now},
		{Code: "TSA", Description: "Tsunami Watch", Category: "Watch", UpdatedAt: now},
		{Code: "AVA", Description: "Avalanche Watch", Category: "Watch", UpdatedAt: now},
		{Code: "DBA", Description: "Dam Watch", Category: "Watch", UpdatedAt: now},
		{Code: "EVA", Description: "Evacuation Watch", Category: "Watch", UpdatedAt: now},
		{Code: "STA", Description: "Storm Surge Watch", Category: "Watch", UpdatedAt: now},
		{Code: "WFA", Description: "Wild Fire Watch", Category: "Watch", UpdatedAt: now},
		// Advisories / Statements
		{Code: "ADR", Description: "Administrative Message", Category: "Statement", UpdatedAt: now},
		{Code: "BHS", Description: "Beach Hazards Statement", Category: "Statement", UpdatedAt: now},
		{Code: "CDA", Description: "Child Abduction Emergency", Category: "Advisory", UpdatedAt: now},
		{Code: "FLS", Description: "Flood Statement", Category: "Statement", UpdatedAt: now},
		{Code: "HLS", Description: "Hurricane Local Statement", Category: "Statement", UpdatedAt: now},
		{Code: "MHW", Description: "Marine Weather Statement", Category: "Advisory", UpdatedAt: now},
		{Code: "MWS", Description: "Marine Weather Statement", Category: "Statement", UpdatedAt: now},
		{Code: "NOW", Description: "Short Term Forecast", Category: "Advisory", UpdatedAt: now},
		{Code: "SCY", Description: "Hazardous Seas", Category: "Advisory", UpdatedAt: now},
		{Code: "SDS", Description: "Special Drought Statement", Category: "Statement", UpdatedAt: now},
		{Code: "SPS", Description: "Special Weather Statement", Category: "Statement", UpdatedAt: now},
		{Code: "SRF", Description: "Surf Advisory", Category: "Advisory", UpdatedAt: now},
		// Tests
		{Code: "EAN", Description: "Emergency Action Notification", Category: "Test", UpdatedAt: now},
		{Code: "EAT", Description: "Emergency Action Termination", Category: "Test", UpdatedAt: now},
		{Code: "EVI", Description: "Immediate Evacuation Warning", Category: "Test", UpdatedAt: now},
		{Code: "NPT", Description: "National Periodic Test", Category: "Test", UpdatedAt: now},
		{Code: "RMT", Description: "Required Monthly Test", Category: "Test", UpdatedAt: now},
		{Code: "RWT", Description: "Required Weekly Test", Category: "Test", UpdatedAt: now},
	}
	return codes, nil
}
