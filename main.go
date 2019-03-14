package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `myDocker is a simple container runtime implementaion`

func main() {
	app := cli.NewApp()
	app.Name = "myDocker"
	app.Usage = usage
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
		listCommand,
		logsCommand,
		execCommand,
	}

	app.Before = func(ctx *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
