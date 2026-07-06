package noaa

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFIPSCSV(t *testing.T) {
	csv := `FIPS,Name
01001,Autauga County, Alabama
13121,Fulton County, Georgia
13067,Cobb County, Georgia
`
	codes, err := parseFIPSCSV(strings.NewReader(csv))
	require.NoError(t, err)
	assert.Len(t, codes, 3)
	assert.Equal(t, "13", codes[1].StateCode)
	assert.Equal(t, "121", codes[1].CountyCode)
	assert.Equal(t, "Georgia", codes[1].StateName)
	assert.Equal(t, "Fulton County", codes[1].CountyName)
}

func TestParseFIPSCSV_SkipsInvalidRows(t *testing.T) {
	csv := `FIPS,Name
bad,Incomplete
13121,Fulton County, Georgia
`
	codes, err := parseFIPSCSV(strings.NewReader(csv))
	require.NoError(t, err)
	assert.Len(t, codes, 1)
}

func TestFetchEventCodes(t *testing.T) {
	f := New()
	codes, err := f.FetchEventCodes(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, codes)

	// Verify required test codes are present
	found := map[string]bool{}
	for _, c := range codes {
		found[c.Code] = true
	}
	assert.True(t, found["RWT"], "RWT (Required Weekly Test) should be present")
	assert.True(t, found["TOR"], "TOR (Tornado Warning) should be present")
}

func TestSplitCountyState(t *testing.T) {
	county, state := splitCountyState("Autauga County, Alabama")
	assert.Equal(t, "Autauga County", county)
	assert.Equal(t, "Alabama", state)

	county, state = splitCountyState("No comma here")
	assert.Equal(t, "No comma here", county)
	assert.Equal(t, "", state)
}
