package commands

import (
	"net/http"

	"time"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/fs"
	"github.com/go-zoox/geoip"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

var GeoIPDatabaseHomeDir = fs.JoinHomeDir(".cache/geoip")
var GeoIPDatabaseFilePath = fs.JoinPath(GeoIPDatabaseHomeDir, "GeoLite2-City.mmdb")
var GeoIPDatabaseDownloadURL = "https://github.com/go-zoox/geoip/releases/download/v0.0.3/GeoLite2-City.mmdb"

func RegistryServer(app *cli.MultipleProgram) {
	app.Register("server", &cli.Command{
		Name:  "server",
		Usage: "geoip server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "server port",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
				Value:   8080,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			return server(ctx.Int("port"))
		},
	})
}

func server(port int) error {
	if ok := fs.IsExist(GeoIPDatabaseHomeDir); !ok {
		if err := fs.Mkdirp(GeoIPDatabaseHomeDir); err != nil {
			return err
		}
	}

	if ok := fs.IsExist(GeoIPDatabaseFilePath); !ok {
		logger.Infof("downloading geoip database from %s to %s", GeoIPDatabaseDownloadURL, GeoIPDatabaseFilePath)
		response, err := fetch.Download(GeoIPDatabaseDownloadURL, GeoIPDatabaseFilePath)
		if err != nil {
			return err
		}

		if !response.Ok() {
			return fmt.Errorf("download geoip database failed, status code: %d", response.StatusCode)
		}
	}

	app := defaults.Default()

	gps := geoip.New(&geoip.Config{
		DatabaseFilePath: GeoIPDatabaseFilePath,
	})

	if err := gps.Load(); err != nil {
		return err
	}
	defer gps.Destroy()

	app.Get("/:ip", func(ctx *zoox.Context) {
		ip := ctx.Param().Get("ip").String()
		if ip == "" {
			ctx.Error(http.StatusBadRequest, "ip is required")
			return
		}

		// check is valid ip address
		if !geoip.IsIPv4(ip) && !geoip.IsIPv6(ip) {
			ctx.Error(http.StatusBadRequest, fmt.Sprintf("invalid ip address: %s", ip))
			return
		}

		address := &geoip.Address{}

		if err := ctx.Cache().Get(ip, address); err != nil {
			address, err = gps.GetAddress(ip)
			if err != nil {
				ctx.Logger.Infof("invalid ip address: %s (err: %s)", ip, err)
				ctx.Error(http.StatusBadRequest, fmt.Sprintf("invalid ip address: %s", ip))
				return
			}

			if err := ctx.Cache().Set(ip, address, 5*time.Minute); err != nil {
				ctx.Logger.Infof("cache ip address: %s failed, err: %s", ip, err)
			}
		}

		ctx.JSON(http.StatusOK, address)
	})

	return app.Run(fmt.Sprintf(":%d", port))
}
