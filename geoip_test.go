package geoip

import (
	"os"
	"testing"
)

func TestGeoIP(t *testing.T) {
	currentDir, _ := os.Getwd()
	geoip := New(&Config{
		DatabaseFilePath: currentDir + "/GeoLite2-City.mmdb",
	})
	geoip.Load()
	defer geoip.Destroy()
	t.Log(geoip.GetAddress("216.8.252.36"))
	t.Log(geoip.GetAddress("216.8.252.36", "en-US"))
}
