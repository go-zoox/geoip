package geoip

import (
	"errors"
	"fmt"
	"log"
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

func GetAddress(ip string, language ...interface{}) (address *Address, addressErr error) {
	defer func() {
		if err := recover(); err != nil {
			// log.Println(err)
			addressErr = errors.New("invalid ip address")
		}
	}()

	langugaes := []string{"zh-CN", "en-US"}
	lang := "zh-CN"
	if len(language) > 0 {
		lang = language[0].(string)
	}

	db, err := geoip2.Open("data/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	_ip := net.ParseIP(ip)
	record, err := db.City(_ip)
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
