package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	RUNNING             string = "running"
	STOP                string = "stoped"
	EXIT                string = "exited"
	InformationLocation string = "/var/run/myDocker/%s/"
	ConfigName          string = "config.json"
	LogFileName         string = "container.log"
)

type ContainerInfo struct {
	Pid        string `json:"pid"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	CreateTime string `json:"createTime"`
	Status     string `json:"status"`
}

func NewParentProcess(containerId string, tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	cmd.ExtraFiles = []*os.File{readPipe}

	mntURL := fmt.Sprintf("/var/run/myDocker/%s/rootfs/", containerId)
	rootURL := fmt.Sprintf("/var/run/myDocker/%s/", containerId)

	NewWorkspace(rootURL, mntURL, volume)

	cmd.Dir = mntURL

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		if file, err := os.Create(rootURL + LogFileName); err == nil {
			cmd.Stdout = file
		} else {
			return nil, nil
		}
	}

	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func NewWorkspace(rootURL string, mntURL string, volume string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateWorkLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)

	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(rootURL, mntURL, volumeURLs)
			logrus.Infof("%q", volumeURLs)
		} else {
			logrus.Infof("volume parameter input is not correct.")
		}
	}
}

func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := "/root/busybox.tar"

	exist, err := PathExists(busyboxURL)
	if err != nil {
		logrus.Infof("fail to judge whether dir %s exists: %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.MkdirAll(busyboxURL, 0777); err != nil {
			logrus.Errorf("make dir %s error: %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			logrus.Errorf("untar dir %s error: %v", busyboxURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		logrus.Errorf("make dir %s error: %v", writeURL, err)
	}
}

func CreateWorkLayer(rootURL string) {
	workURL := rootURL + "workLayer"
	if err := os.Mkdir(workURL, 0777); err != nil {
		logrus.Errorf("make dir %s error: %v", workURL, err)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		logrus.Errorf("make dir %s error: %v", mntURL, err)
	}
	dirs := "lowerdir=" + rootURL + "busybox,upperdir=" + rootURL + "writeLayer,workdir=" + rootURL + "workLayer"
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount error: %v", err)
	}
}

func MountVolume(rootURL string, mntURL string, volumeURLs []string) {
	parentURL := volumeURLs[0]
	if err := os.Mkdir(parentURL, 0777); err != nil {
		logrus.Infof("make dir %s error: %v", parentURL, err)
	}
	contrainerURL := volumeURLs[1]
	contrainerVolumeURL := mntURL + contrainerURL
	if err := os.Mkdir(contrainerVolumeURL, 0777); err != nil {
		logrus.Infof("make dir %s error: %v", contrainerVolumeURL, err)
	}
	cmd := exec.Command("mount", "--bind", parentURL, contrainerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount volume failed: %v", err)
	}
}

func DeleteWorkspace(rootURL string, mntURL string, volume string) {
	if volume != "" {
		DeleteVolume(mntURL, strings.Split(volume, ":"))
		logrus.Infof("success for umount volume.")
	}
	DeleteMountPoint(rootURL, mntURL)
	DeleteWriteLayer(rootURL)
	DeleteWorkLayer(rootURL)
}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("umount error: %v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("remove dir %s error: %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer/"
	if err := os.RemoveAll(writeURL); err != nil {
		logrus.Errorf("remove dir %s error: %v", writeURL, err)
	}
}

func DeleteWorkLayer(rootURL string) {
	workURL := rootURL + "workLayer"
	if err := os.RemoveAll(workURL); err != nil {
		logrus.Errorf("remove dir %s error: %v", workURL, err)
	}
}

func DeleteVolume(mntURL string, volumeURLs []string) {
	contrainerURL := mntURL + volumeURLs[1]
	cmd := exec.Command("umount", contrainerURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("umount volume failed: %v", err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
