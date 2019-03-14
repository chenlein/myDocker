package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chenlei/myDocker/cgroups"
	"github.com/chenlei/myDocker/cgroups/subsystems"
	"github.com/chenlei/myDocker/container"
	"github.com/sirupsen/logrus"
)

func Run(tty bool, command []string, resourceConfig *subsystems.ResourceConfig, volume string, containerName string) {
	containerId := randStringBytes(10)

	parent, writePipe := container.NewParentProcess(containerId, tty, volume)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	containerName, err := recordContainerInformation(containerId, parent.Process.Pid, command, containerName)
	if err != nil {
		logrus.Errorf("record container information error: %v", err)
		return
	}

	cgroupManager := cgroups.NewCgroupManager("docker-cgroup")
	defer cgroupManager.Destory()
	cgroupManager.Set(resourceConfig)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(command, writePipe)

	if tty {
		parent.Wait()

		mntURL := fmt.Sprintf("/var/run/myDocker/%s/rootfs/", containerId)
		rootURL := fmt.Sprintf("/var/run/myDocker/%s/", containerId)
		container.DeleteWorkspace(rootURL, mntURL, volume)

		deleteContainerInformation(containerId)
	}

	os.Exit(1)
}

func sendInitCommand(commandArray []string, writePipe *os.File) {
	command := strings.Join(commandArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func recordContainerInformation(containerId string, containerPID int, commandArray []string, containerName string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	if containerName == "" {
		containerName = containerId
	}
	containerInformation := &container.ContainerInfo{
		Id:         containerId,
		Pid:        strconv.Itoa(containerPID),
		Command:    command,
		CreateTime: createTime,
		Status:     container.RUNNING,
		Name:       containerName,
	}
	jsonBytes, err := json.Marshal(containerInformation)
	if err != nil {
		logrus.Errorf("record container info error: %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirURL := fmt.Sprintf(container.InformationLocation, containerId)
	if err := os.MkdirAll(dirURL, 0622); err != nil {
		logrus.Errorf("make dir %s error: %v", dirURL, err)
		return "", err
	}
	fileName := dirURL + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		logrus.Errorf("create file %s error: %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		logrus.Errorf("file write error: %v", err)
		return "", err
	}
	return containerName, nil
}

func deleteContainerInformation(containerId string) {
	dirURL := fmt.Sprintf(container.InformationLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		logrus.Errorf("remove dir %s error: %v", dirURL, err)
	}
}
