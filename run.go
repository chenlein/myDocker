package main

import (
	"os"
	"strings"

	"github.com/chenlei/myDocker/cgroups"
	"github.com/chenlei/myDocker/cgroups/subsystems"
	"github.com/chenlei/myDocker/container"
	"github.com/sirupsen/logrus"
)

func Run(tty bool, command []string, resourceConfig *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	cgroupManager := cgroups.NewCgroupManager("docker-cgroup")
	defer cgroupManager.Destory()
	cgroupManager.Set(resourceConfig)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(command, writePipe)
	parent.Wait()

	mntURL := "/root/mnt/"
	rootURL := "/root/"
	container.DeleteWorkspace(rootURL, mntURL, volume)

	os.Exit(1)
}

func sendInitCommand(commandArray []string, writePipe *os.File) {
	command := strings.Join(commandArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
