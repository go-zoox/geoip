package geoip

import (
	"testing"
)

func TestGeoIP(t *testing.T) {
	t.Log(GetAddress("216.8.252.36"))
	t.Log(GetAddress("216.8.252.36", "en-US"))
}
