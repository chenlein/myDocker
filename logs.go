package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chenlei/myDocker/container"
	"github.com/sirupsen/logrus"
)

func showContainerLogs(containerName string) {
	rootURL := fmt.Sprintf(container.InformationLocation, containerName)
	logFileLocation := rootURL + container.LogFileName
	file, err := os.Open(logFileLocation)
	defer file.Close()
	if err != nil {
		logrus.Errorf("open file %s error: %v", logFileLocation, err)
	}
	if logs, err := ioutil.ReadAll(file); err == nil {
		fmt.Fprint(os.Stdout, string(logs))
	} else {
		logrus.Error("read file %s error: %v", logFileLocation, err)
	}
}
