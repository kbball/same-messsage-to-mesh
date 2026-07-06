package sdr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSAME(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		wantCode  string
		wantOrg   string
		wantFIPS  []string
		wantPurge string
		wantIssue string
		wantCall  string
		wantErr   bool
	}{
		{
			name:      "tornado warning",
			raw:       "ZCZC-WXR-TOR-037121+0030-1820218-KRAH/NWS-",
			wantCode:  "TOR",
			wantOrg:   "WXR",
			wantFIPS:  []string{"037121"},
			wantPurge: "0030",
			wantIssue: "1820218",
			wantCall:  "KRAH/NWS",
		},
		{
			name:      "required weekly test with multiple counties",
			raw:       "ZCZC-WXR-RWT-013121-013067-013135+0030-1820218-KFFC/NWS-",
			wantCode:  "RWT",
			wantOrg:   "WXR",
			wantFIPS:  []string{"013121", "013067", "013135"},
			wantPurge: "0030",
			wantIssue: "1820218",
			wantCall:  "KFFC/NWS",
		},
		{
			name:     "too few fields",
			raw:      "ZCZC-WXR-",
			wantErr:  true,
		},
		{
			name:     "no FIPS codes",
			raw:      "ZCZC-WXR-TOR+0030-1820218-KRAH/NWS-",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alert, err := ParseSAME(tt.raw)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantOrg, alert.Originator)
			assert.Equal(t, tt.wantCode, alert.EventCode)
			assert.Equal(t, tt.wantFIPS, alert.FIPSCodes)
			assert.Equal(t, tt.wantPurge, alert.PurgeTime)
			assert.Equal(t, tt.wantIssue, alert.IssueTime)
			assert.Equal(t, tt.wantCall, alert.CallSign)
			assert.Equal(t, tt.raw, alert.RawMessage)
		})
	}
}

