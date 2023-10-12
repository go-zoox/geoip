package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/geoip"
	"github.com/go-zoox/geoip/cmd/geoip/commands"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "geoip",
		Usage:   "geoip is a portable ip database/server/client",
		Version: geoip.Version,
	})

	// server
	commands.RegistryServer(app)
	// client
	commands.RegistryClient(app)

	app.Run()
}
