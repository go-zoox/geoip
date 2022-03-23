package geoip

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type Address struct {
	Country     string
	Province    string
	City        string
	CountryCode string
	TimeZone    string
	Coordinates []float64
}

type GeoIP struct {
	db     *geoip2.Reader
	config *Config
}

type Config struct {
	DatabaseFilePath string
}

func New(c *Config) *GeoIP {
	return &GeoIP{
		config: c,
	}
}

func (g *GeoIP) Load() error {
	db, err := geoip2.Open(g.config.DatabaseFilePath)
	if err != nil {
		return err
	}
	g.db = db
	return nil
}

func (g *GeoIP) Destroy() error {
	return g.db.Close()
}

func (g *GeoIP) GetAddress(ip string, language ...interface{}) (address *Address, addressErr error) {
	defer func() {
		if err := recover(); err != nil {
			// log.Println(err)
			addressErr = errors.New("invalid ip address")
		}
	}()

	if g.db == nil {
		return nil, errors.New("geoip database not initialized, you can download database file from https://github.com/go-zoox/geoip/releases/download/v0.0.3/GeoLite2-City.mmdb")
	}

	langugaes := []string{"zh-CN", "en-US"}
	lang := "zh-CN"
	if len(language) > 0 {
		lang = language[0].(string)
	}

	// If you are using strings that may be invalid, check that ip is not nil
	_ip := net.ParseIP(ip)
	record, err := g.db.City(_ip)
	if err != nil {
		return nil, err
	}

	switch lang {
	case "zh-CN":
		return &Address{
			Province:    record.Subdivisions[0].Names["zh-CN"],
			City:        record.City.Names["zh-CN"],
			Country:     record.Country.Names["zh-CN"],
			CountryCode: record.Country.IsoCode,
			TimeZone:    record.Location.TimeZone,
			Coordinates: []float64{
				record.Location.Latitude,
				record.Location.Longitude,
			},
		}, nil
	case "en-US":
		return &Address{
			Province:    record.Subdivisions[0].Names["en"],
			City:        record.City.Names["en"],
			Country:     record.Country.Names["en"],
			CountryCode: record.Country.IsoCode,
			TimeZone:    record.Location.TimeZone,
			Coordinates: []float64{
				record.Location.Latitude,
				record.Location.Longitude,
			},
		}, nil
	default:
		return nil, fmt.Errorf("language %s not support in (%s)", lang, strings.Join(langugaes, ","))
	}
}
