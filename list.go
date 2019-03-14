package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/chenlei/myDocker/container"
	"github.com/sirupsen/logrus"
)

func ListContainers() {
	dirURL := fmt.Sprintf(container.InformationLocation, "")
	dirURL = dirURL[:len(dirURL)-1]
	files, err := ioutil.ReadDir(dirURL)
	if err != nil {
		logrus.Errorf("read dir %s error: %v", dirURL, err)
		return
	}

	var containers []*container.ContainerInfo
	for _, file := range files {
		information, err := getContainerInformation(file)
		if err != nil {
			logrus.Errorf("get container information error: %v", err)
			continue
		}
		containers = append(containers, information)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", item.Id, item.Name, item.Pid, item.Status, item.Command, item.CreateTime)
	}
	if err := w.Flush(); err != nil {
		logrus.Errorf("flush error: %v", err)
		return
	}
}

func getContainerInformation(file os.FileInfo) (*container.ContainerInfo, error) {
	containerName := file.Name()
	configFileDir := fmt.Sprintf(container.InformationLocation, containerName)
	configFileDir = configFileDir + container.ConfigName
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		logrus.Errorf("read file %s error: %v", configFileDir, err)
		return nil, err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		logrus.Errorf("json unmarshal error: %v", err)
		return nil, err
	}
	return &containerInfo, nil
}
