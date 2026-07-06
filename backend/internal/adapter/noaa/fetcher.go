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
	// Census Bureau national county file — pipe-delimited, very reliable.
	// Columns: STATE|STATEFP|COUNTYFP|COUNTYNS|COUNTYNAME
	fipsURL       = "https://www2.census.gov/geo/docs/reference/codes2020/national_county2020.txt"
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

// FetchFIPS downloads county FIPS codes from the Census Bureau national county file.
// Format: pipe-delimited with header STATE|STATEFP|COUNTYFP|COUNTYNS|COUNTYNAME
func (f *Fetcher) FetchFIPS(ctx context.Context) ([]entity.FIPSCode, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fipsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching FIPS data: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from Census FIPS endpoint", resp.StatusCode)
	}

	return parseFIPSPipe(resp.Body)
}

// parseFIPSPipe parses the Census Bureau national_county2020.txt pipe-delimited file.
// Columns: STATE|STATEFP|COUNTYFP|COUNTYNS|COUNTYNAME
func parseFIPSPipe(r io.Reader) ([]entity.FIPSCode, error) {
	reader := csv.NewReader(r)
	reader.Comma = '|'
	reader.TrimLeadingSpace = true

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	now := time.Now()
	var codes []entity.FIPSCode
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading row: %w", err)
		}
		// STATE|STATEFP|COUNTYFP|COUNTYNS|COUNTYNAME
		if len(record) < 5 {
			continue
		}
		stateCode := strings.TrimSpace(record[1])
		countyCode := strings.TrimSpace(record[2])
		countyName := strings.TrimSpace(record[4])
		if len(stateCode) != 2 || len(countyCode) != 3 || countyName == "" {
			continue
		}
		stateName := stateFIPSName[stateCode]

		codes = append(codes, entity.FIPSCode{
			StateCode:  stateCode,
			CountyCode: countyCode,
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

// stateFIPSName maps 2-digit state FIPS codes to state names.
var stateFIPSName = map[string]string{
	"01": "Alabama", "02": "Alaska", "04": "Arizona", "05": "Arkansas",
	"06": "California", "08": "Colorado", "09": "Connecticut", "10": "Delaware",
	"11": "District of Columbia", "12": "Florida", "13": "Georgia", "15": "Hawaii",
	"16": "Idaho", "17": "Illinois", "18": "Indiana", "19": "Iowa",
	"20": "Kansas", "21": "Kentucky", "22": "Louisiana", "23": "Maine",
	"24": "Maryland", "25": "Massachusetts", "26": "Michigan", "27": "Minnesota",
	"28": "Mississippi", "29": "Missouri", "30": "Montana", "31": "Nebraska",
	"32": "Nevada", "33": "New Hampshire", "34": "New Jersey", "35": "New Mexico",
	"36": "New York", "37": "North Carolina", "38": "North Dakota", "39": "Ohio",
	"40": "Oklahoma", "41": "Oregon", "42": "Pennsylvania", "44": "Rhode Island",
	"45": "South Carolina", "46": "South Dakota", "47": "Tennessee", "48": "Texas",
	"49": "Utah", "50": "Vermont", "51": "Virginia", "53": "Washington",
	"54": "West Virginia", "55": "Wisconsin", "56": "Wyoming",
	"60": "American Samoa", "66": "Guam", "69": "Northern Mariana Islands",
	"72": "Puerto Rico", "78": "U.S. Virgin Islands",
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
