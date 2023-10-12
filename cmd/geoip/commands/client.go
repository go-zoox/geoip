package commands

import (
	"github.com/go-zoox/cli"
)

func RegistryClient(app *cli.MultipleProgram) {
	app.Register("client", &cli.Command{
		Name:  "client",
		Usage: "geoip client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "server",
				Usage:    "server url",
				Aliases:  []string{"s"},
				EnvVars:  []string{"SERVER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "username",
				Usage:   "Username for Basic Auth",
				EnvVars: []string{"USERNAME"},
			},
			&cli.StringFlag{
				Name:    "password",
				Usage:   "Password for Basic Auth",
				EnvVars: []string{"PASSWORD"},
			},
			&cli.StringFlag{
				Name:    "ip",
				Usage:   "specify search ip",
				Aliases: []string{"i"},
				EnvVars: []string{"IP"},
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			return client(
				ctx.String("ip"),
				ctx.String("server"),
				ctx.String("username"),
				ctx.String("password"),
			)
		},
	})
}

func client(ip, server, username, password string) error {
	return nil
}
