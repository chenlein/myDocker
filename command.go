package main

import (
	"fmt"
	"os"

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
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
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

		containerName := ctx.String("name")

		Run(tty, commandArray, resourceConfig, volume, containerName)
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

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list container",
	Action: func(ctx *cli.Context) error {
		ListContainers()
		return nil
	},
}

var logsCommand = cli.Command{
	Name:  "logs",
	Usage: "show container logs",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) != 1 {
			logrus.Errorf("missing container name")
		}
		showContainerLogs(ctx.Args().Get(0))
		return nil
	},
}

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "execute a command into container",
	Action: func(ctx *cli.Context) error {
		if os.Getenv(ENV_EXEC_PID) != "" {
			logrus.Infof("execute callback pid: %s", os.Getpid())
			logrus.Infof("exec pid: %s", os.Getenv(ENV_EXEC_PID))
			logrus.Infof("exec command: %s", os.Getenv(ENV_EXEC_CMD))
			return nil
		}
		if len(ctx.Args()) < 2 {
			return fmt.Errorf("missing ontainer name or command")
		}
		var commandArray []string
		for _, arg := range ctx.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		executeIntoContainer(ctx.Args().Get(0), commandArray)
		return nil
	},
}
