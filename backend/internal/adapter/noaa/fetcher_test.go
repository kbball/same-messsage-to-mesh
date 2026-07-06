package noaa

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFIPSPipe(t *testing.T) {
	data := `STATE|STATEFP|COUNTYFP|COUNTYNS|COUNTYNAME
AL|01|001|00161526|Autauga County
GA|13|121|00351257|Fulton County
GA|13|067|00351224|Cobb County
`
	codes, err := parseFIPSPipe(strings.NewReader(data))
	require.NoError(t, err)
	assert.Len(t, codes, 3)
	assert.Equal(t, "13", codes[1].StateCode)
	assert.Equal(t, "121", codes[1].CountyCode)
	assert.Equal(t, "Georgia", codes[1].StateName)
	assert.Equal(t, "Fulton County", codes[1].CountyName)
}

func TestParseFIPSPipe_SkipsInvalidRows(t *testing.T) {
	data := `STATE|STATEFP|COUNTYFP|COUNTYNS|COUNTYNAME
XX|bad|00|12345|Incomplete
GA|13|121|00351257|Fulton County
`
	codes, err := parseFIPSPipe(strings.NewReader(data))
	require.NoError(t, err)
	assert.Len(t, codes, 1)
}

func TestFetchEventCodes(t *testing.T) {
	f := New()
	codes, err := f.FetchEventCodes(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, codes)

	found := map[string]bool{}
	for _, c := range codes {
		found[c.Code] = true
	}
	assert.True(t, found["RWT"], "RWT (Required Weekly Test) should be present")
	assert.True(t, found["TOR"], "TOR (Tornado Warning) should be present")
}

func TestStateFIPSName(t *testing.T) {
	assert.Equal(t, "Georgia", stateFIPSName["13"])
	assert.Equal(t, "Alabama", stateFIPSName["01"])
	assert.Equal(t, "Puerto Rico", stateFIPSName["72"])
}
