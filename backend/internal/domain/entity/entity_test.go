package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFIPSCode_FIPS(t *testing.T) {
	f := FIPSCode{StateCode: "13", CountyCode: "121"}
	assert.Equal(t, "13121", f.FIPS())
}

func TestAlertFilter_EmptyMeansAll(t *testing.T) {
	f := AlertFilter{}
	assert.Empty(t, f.StateCodes)
	assert.Empty(t, f.FIPSCodes)
	assert.Empty(t, f.EventCodes)
}
