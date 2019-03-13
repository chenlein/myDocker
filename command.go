package main

import (
	"fmt"

	"github.com/chenlei/myDocker/cgroups/subsystems"
	"github.com/chenlei/myDocker/container"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process.",
	Action: func(ctx *cli.Context) error {
		logrus.Infof("Init come on")
		cmd := ctx.Args().Get(0)
		logrus.Infof("command %s", cmd)
		err := container.RunContainerInitProcess()
		return err
	},
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.StringFlag{
			Name:  "d",
			Usage: "detach",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var commandArray []string
		for _, arg := range ctx.Args() {
			commandArray = append(commandArray, arg)
		}
		tty := ctx.Bool("ti")
		detach := ctx.Bool("d")
		if tty && detach {
			return fmt.Errorf("ti and d parameter can not both provided.")
		}

		resourceConfig := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
		}
		volume := ctx.String("v")
		Run(tty, commandArray, resourceConfig, volume)
		return nil
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container into image",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}
		imageName := ctx.Args().Get(0)
		CommitContainer(imageName)
		return nil
	},
}
