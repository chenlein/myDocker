package cgroups

import (
	"github.com/chenlei/myDocker/cgroups/subsystems"
	"github.com/sirupsen/logrus"
)

type CgroupManager struct {
	Path     string
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subsystemInstance := range subsystems.SubsystemInstances {
		subsystemInstance.Apply(c.Path, pid)
	}
	return nil
}

func (c *CgroupManager) Set(config *subsystems.ResourceConfig) error {
	for _, subsystemInstance := range subsystems.SubsystemInstances {
		subsystemInstance.Set(c.Path, config)
	}
	return nil
}

func (c *CgroupManager) Destory() error {
	for _, subsystemInstance := range subsystems.SubsystemInstances {
		if err := subsystemInstance.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
