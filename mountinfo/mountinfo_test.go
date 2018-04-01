// Copyright 2018 Raymond Barbiero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package mountinfo

import (
	"testing"
	"github.com/visualphoenix/disk-go/lsblk"
)

func TestMounts(t *testing.T) {
	raw := `disk   xvda
part / xfs xvda1
disk   xvdb
part  LVM2_member xvdb1
lvm /home ext4 test-home_vol
lvm /usr/local ext4 test-local_vol
lvm /opt ext4 test-opt_vol
lvm  ext4 test-tmp_vol
lvm /var ext4 test-var_vol
disk   xvdk
part  LVM2_member xvdk1
lvm /var ext4 test-var_vol
`
    res, err := lsblk.ParseRawLsblk(raw)
	if err != nil {
		t.Errorf("parse error: %s", err)
	}
	got := GetMountInfoFrom(res)
	want := []MountInfo{
		MountInfo{
			Mountpoint: "/",
			FilesystemType: "xfs",
			BlockDevice: "xvda1",
			BlockDeviceType: "part",
			PhysicalDevices: []string{"xvda"},
		},
		MountInfo{
			Mountpoint: "/home",
			FilesystemType: "ext4",
			BlockDevice: "test-home_vol",
			BlockDeviceType: "lvm",
			PhysicalDevices: []string{"xvdb"},
		},
		MountInfo{
			Mountpoint: "/usr/local",
			FilesystemType: "ext4",
			BlockDevice: "test-local_vol",
			BlockDeviceType: "lvm",
			PhysicalDevices: []string{"xvdb"},
		},
		MountInfo{
			Mountpoint: "/opt",
			FilesystemType: "ext4",
			BlockDevice: "test-opt_vol",
			BlockDeviceType: "lvm",
			PhysicalDevices: []string{"xvdb"},
		},
		MountInfo{
			Mountpoint: "/var",
			FilesystemType: "ext4",
			BlockDevice: "test-var_vol",
			BlockDeviceType: "lvm",
			PhysicalDevices: []string{"xvdb", "xvdk"},
		},
		MountInfo{
			Mountpoint: "/var",
			FilesystemType: "ext4",
			BlockDevice: "test-var_vol",
			BlockDeviceType: "lvm",
			PhysicalDevices: []string{"xvdb", "xvdk"},
		},
	}
	if len(got) != len(want) {
		t.Errorf("got: %+v\nwant: %+v\n", got, want)
	}
}
