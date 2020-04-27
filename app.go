package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
)

func newApp() *cli.App {
	config := NewDefaultConfig()
	app := cli.NewApp()
	app.Name = prog
	app.Usage = description
	app.Version = version
	app.Flags = getCommandLineOptions()
	app.UsageText = "app-name [options]"

	app.Action = func(ctx *cli.Context) error {
		if err := parseCLIOptions(ctx, config); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		server, err := createNewServer(config)
		if err != nil {
			return err
		}

		server.Run()

		defer server.logger.Sync()

		// Setup the termination signals
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-signalChannel

		return nil
	}

	return app
}
