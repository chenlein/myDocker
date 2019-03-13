package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

type MemorySubsystem struct {
}

func (s *MemorySubsystem) Name() string {
	return "memory"
}

func (s *MemorySubsystem) Set(cgroupPath string, config *ResourceConfig) error {
	if subsystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		logrus.Infof("subsystem cgroup path %s", subsystemCgroupPath)
		if config.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(subsystemCgroupPath, "memory.limit_in_bytes"), []byte(config.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *MemorySubsystem) Remove(cgroupPath string) error {
	if subsystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.Remove(subsystemCgroupPath)
	} else {
		return err
	}
}

func (s *MemorySubsystem) Apply(cgroupPath string, pid int) error {
	if subsystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subsystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail%v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}
