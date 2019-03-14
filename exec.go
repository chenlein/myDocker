package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/chenlei/myDocker/container"
	_ "github.com/chenlei/myDocker/nsenter"
	"github.com/sirupsen/logrus"
)

const ENV_EXEC_PID = "myDocker_pid"
const ENV_EXEC_CMD = "myDocker_cmd"

func executeIntoContainer(containerName string, commandArray []string) {
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		logrus.Errorf("get container %s pid error: %v", containerName, err)
		return
	}
	logrus.Infof("container pid: %s", pid)
	logrus.Infof("command: %s", strings.Join(commandArray, " "))

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, strings.Join(commandArray, " "))
	cmd.Env = append(os.Environ(), getEnvByPid(pid)...)

	if err := cmd.Run(); err != nil {
		logrus.Errorf("execute container %s error: %v", containerName, err)
	}
}

func getContainerPidByName(containerName string) (string, error) {
	rootURL := fmt.Sprintf(container.InformationLocation, containerName)

	bytes, err := ioutil.ReadFile(rootURL + container.ConfigName)

	if err != nil {
		return "", err
	}
	var containerInformation container.ContainerInfo
	if err := json.Unmarshal(bytes, &containerInformation); err != nil {
		return "", err
	}
	return containerInformation.Pid, nil
}

func getEnvByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	if bytes, err := ioutil.ReadFile(path); err != nil {
		logrus.Errorf("read file %s error: %v", path, err)
		return nil
	} else {
		return strings.Split(string(bytes), "\u0000")
	}
}
