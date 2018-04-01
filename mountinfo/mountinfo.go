package mountinfo
// Copyright 2018 Raymond Barbiero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"github.com/visualphoenix/disk-go/fs"
	"github.com/visualphoenix/disk-go/lvm"
	"github.com/visualphoenix/disk-go/lsblk"
)

// MountInfo is a handy representation of mount info
type MountInfo struct {
	Mountpoint string
	FilesystemType string
	BlockDevice string
	BlockDeviceType string
	PhysicalDevices []string
}

// GetMountInfoFrom returns a list of MountInfo from lsblk struct info
func GetMountInfoFrom(l lsblk.Lsblk) []MountInfo {
	var result []MountInfo;
	mountpointToDisks := make(map[string][]string)
	for _, d := range l.Disks {
		if d.Disk.Mountpoint != "" {
			mountpointToDisks[d.Disk.Mountpoint] =  append(mountpointToDisks[d.Disk.Mountpoint],d.Disk.Device)
			m := MountInfo {
				Mountpoint: d.Disk.Mountpoint,
				FilesystemType: d.Disk.Fstype,
				BlockDevice: d.Disk.Device,
				BlockDeviceType: d.Disk.Dtype,
			}
			result = append(result, m)

		}
		for _, p := range d.Parts {
			if p.Mountpoint != "" {
				mountpointToDisks[p.Mountpoint] =  append(mountpointToDisks[p.Mountpoint],d.Disk.Device)
				m := MountInfo {
					Mountpoint: p.Mountpoint,
					FilesystemType: p.Fstype,
					BlockDevice: p.Device,
					BlockDeviceType: p.Dtype,
				}
				result = append(result, m)
			}
		}
	}
	for i := range result {
		disks := mountpointToDisks[result[i].Mountpoint]
		result[i].PhysicalDevices = disks
	}
	return result
}

// GetMountInfo returns a list of MountInfo
func GetMountInfo() ([]MountInfo, error) {
	raw, err := lsblk.ExecLsblk()
	if err != nil {
		return []MountInfo{}, fmt.Errorf("lsblk exec error: %s", err)
	}
	l, err := lsblk.ParseRawLsblk(raw)
	if err != nil {
		return []MountInfo{}, fmt.Errorf("parse error: %s", err)
	}
	mi := GetMountInfoFrom(l)
	return mi, nil
}

// Suspend writes to the device/partition given the type of the mount
func (mi MountInfo) Suspend() {
	switch mi.BlockDeviceType {
	case "disk":
		fallthrough
	case "part":
		fs.Freeze(mi.Mountpoint)
	case "lvm":
		lvm.Suspend(mi.BlockDevice)
	default:
	}
}

// Resume writes to the device/partition given the type of the mount
func (mi MountInfo) Resume() {
	switch mi.BlockDeviceType {
	case "disk":
		fallthrough
	case "part":
		fs.Freeze(mi.Mountpoint)
	case "lvm":
		lvm.Suspend(mi.BlockDevice)
	default:
	}
}
