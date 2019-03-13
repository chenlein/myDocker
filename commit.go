package main

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func CommitContainer(imageName string) {
	if _, err := exec.Command("tar", "-czf", "/root/"+imageName+".tar", "-C", "/root/mnt", ".").CombinedOutput(); err != nil {
		logrus.Errorf("tar folder /root/mnt error: %v", err)
	}
}
