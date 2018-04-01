// Copyright 2018 Raymond Barbiero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package lsblk

import (
	"fmt"
	"testing"
)

func deepEquals(l Lsblk, r Lsblk) error {
	if l.Raw != r.Raw {
		return fmt.Errorf("l: %s, r: %s", l.Raw, r.Raw)
	}
	if len(l.Disks) != len(r.Disks) {
		return fmt.Errorf("l: %d, r %d", len(l.Disks), len(r.Disks))
	}
	for i := range l.Disks {
		if l.Disks[i].Disk != r.Disks[i].Disk {
			return fmt.Errorf("l: %+v\nr: %+v", l, r)
		}
		glen, wlen := len(l.Disks[i].Parts), len(r.Disks[i].Parts)
		if glen != wlen {
			return fmt.Errorf("l: %d, r: %d", glen, wlen)
		}
		for j := range l.Disks[i].Parts {
			if l.Disks[i].Parts[j] != r.Disks[i].Parts[j] {
				return fmt.Errorf("l: %+v\nr: %+v", l.Disks[i].Parts[j], r.Disks[i].Parts[j])
			}
		}
	}
	return nil
}

func equal(l Lsblk, r Lsblk) bool {
	res := deepEquals(l, r)
	return res == nil
}

func TestDisk_LsblkRawLineToNode(t *testing.T) {
	got, _ := RawLineToNode("disk   xvda")
	want := Node{Dtype: "disk", Device: "xvda"}
	if got != want {
		t.Errorf("got: %+v\nwant: %+v\n", got, want)
	}
}

func TestDiskLVM2_Member_LsblkLsblkRawLineToNode(t *testing.T) {
	got, _ := RawLineToNode("disk  LVM2_member xvdd")
	want := Node{Dtype: "disk", Fstype: "LVM2_member", Device: "xvdd"}
	if got != want {
		t.Errorf("got: %+v\nwant: %+v\n", got, want)
	}
}

func TestPartMounted_LsblkLsblkRawLineToNode(t *testing.T) {
	got, _ := RawLineToNode("part / xfs xvda1")
	want := Node{Dtype: "part", Mountpoint: "/", Fstype: "xfs", Device: "xvda1"}

	if got != want {
		t.Errorf("got: %+v\nwant: %+v\n", got, want)
	}
}

func TestPartLVM2_Member_LsblkLsblkRawLineToNode(t *testing.T) {
	got, _ := RawLineToNode("part  LVM2_member xvdb1")
	want := Node{Dtype: "part", Mountpoint: "", Fstype: "LVM2_member", Device: "xvdb1"}

	if got != want {
		t.Errorf("got: %+v\nwant: %+v\n", got, want)
	}
}

func TestLvmMounted_LsblkLsblkRawLineToNode(t *testing.T) {
	got, _ := RawLineToNode("lvm /usr/local ext4 e-local_vol")
	want := Node{Dtype: "lvm", Mountpoint: "/usr/local", Fstype: "ext4", Device: "e-local_vol"}

	if got != want {
		t.Errorf("got: %+v\nwant: %+v\n", got, want)
	}
}

func TestLvm_LsblkLsblkRawLineToNode(t *testing.T) {
	got, _ := RawLineToNode("lvm   docker-thinpool_tmeta")
	want := Node{Dtype: "lvm", Mountpoint: "", Fstype: "", Device: "docker-thinpool_tmeta"}

	if got != want {
		t.Errorf("got: %+v\nwant: %+v\n", got, want)
	}
}
func TestParseComplex(t *testing.T) {
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
	want := Lsblk{
		Raw: raw,
		Disks: []Disk{
			Disk{
				Disk: Node{Dtype: "disk", Device: "xvda"},
				Parts: []Node{
					Node{Dtype: "part", Device: "xvda1", Fstype: "xfs", Mountpoint: "/"},
				},
			},
			Disk{
				Disk: Node{Dtype: "disk", Device: "xvdb"},
				Parts: []Node{
					Node{Dtype: "part", Device: "xvdb1", Fstype: "LVM2_member", Mountpoint: ""},
					Node{Dtype: "lvm", Device: "test-home_vol", Fstype: "ext4", Mountpoint: "/home"},
					Node{Dtype: "lvm", Device: "test-local_vol", Fstype: "ext4", Mountpoint: "/usr/local"},
					Node{Dtype: "lvm", Device: "test-opt_vol", Fstype: "ext4", Mountpoint: "/opt"},
					Node{Dtype: "lvm", Device: "test-tmp_vol", Fstype: "ext4", Mountpoint: ""},
					Node{Dtype: "lvm", Device: "test-var_vol", Fstype: "ext4", Mountpoint: "/var"},
				},
			},
			Disk{
				Disk: Node{Dtype: "disk", Device: "xvdk"},
				Parts: []Node{
					Node{Dtype: "part", Device: "xvdk1", Fstype: "LVM2_member", Mountpoint: ""},
					Node{Dtype: "lvm", Device: "test-var_vol", Fstype: "ext4", Mountpoint: "/var"},
				},
			},
		},
	}
	got, err := ParseRawLsblk(raw)
	if err != nil {
		t.Errorf("parse error: %s", err)
	}
	if !equal(want, got) {
		err := deepEquals(want, got)
		t.Errorf("%s", err)
	}
}

func TestParse(t *testing.T) {
	raw := `disk   xvda
part / xfs xvda1
`
	want := Lsblk{
		Raw: raw,
		Disks: []Disk{
			Disk{
				Disk: Node{Dtype: "disk", Device: "xvda"},
				Parts: []Node{
					Node{Dtype: "part", Device: "xvda1", Fstype: "xfs", Mountpoint: "/"},
				},
			},
		},
	}
	got, err := ParseRawLsblk(raw)
	if err != nil {
		t.Errorf("parse error: %s", err)
	}
	if !equal(want, got) {
		err := deepEquals(want, got)
		t.Errorf("%s", err)
	}
}
