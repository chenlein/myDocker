package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	commandArray := readUserCommand()
	if commandArray == nil || len(commandArray) == 0 {
		return fmt.Errorf("container get user command error, command array is nil")
	}

	//defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	//syscall.Mount("", "/", "", uintptr(syscall.MS_PRIVATE|syscall.MS_REC|defaultMountFlags), "")
	//syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	setUpMount()

	path, err := exec.LookPath(commandArray[0])
	if err != nil {
		logrus.Errorf("exec look path error %v", err)
		return err
	}
	logrus.Infof("find path %s", path)
	if err := syscall.Exec(path, commandArray[:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	message, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	return strings.Split(string(message), " ")
}

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0755); err != nil {
		fmt.Println("--------")
		return err
	}
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot root: %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot root dir %v", err)
	}
	return os.Remove(pivotDir)
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("get current location error %v", err)
		return
	}
	logrus.Infof("current location is %s", pwd)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("", "/", "", uintptr(syscall.MS_PRIVATE|syscall.MS_REC|defaultMountFlags), "")

	pivotRoot(pwd)

	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}
