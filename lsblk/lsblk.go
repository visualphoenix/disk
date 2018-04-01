package lsblk
// Copyright 2018 Raymond Barbiero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os/exec"
	"strings"
)

// Node represents an entry in lsblk output
type Node struct {
	Dtype      string
	Device     string
	Fstype     string
	Mountpoint string
}

// Disk represents a base block device made of parts
type Disk struct {
	Disk  Node
	Parts []Node
}

// Lsblk is a collection of disk info
type Lsblk struct {
	Raw   string
	Disks []Disk
}

// ExecLsblk runs the lsblk command and returns the output
func ExecLsblk() (string, error) {
	//lsblk -o TYPE,MOUNTPOINT,FSTYPE,NAME -r -n -i
	cmd := exec.Command("lsblk", "-o", "TYPE,MOUNTPOINT,FSTYPE,NAME", "-r", "-i", "-n")
	out, err := cmd.Output()
	if err != nil {
		exitErr := err.(*exec.ExitError)
		return "", fmt.Errorf("Lsblk: %v, stderr: %s", err, exitErr.Stderr)
	}
	return string(out), nil
}

// RawLineToNode converts a string to a Node
func RawLineToNode(line string) (Node, error) {
	n := Node{}
	line = strings.TrimSpace(line)
	r := strings.Split(line, " ")
	if len(r) != 4 {
		return n, fmt.Errorf("unexpected number of results. Expected 4, got %d", len(r))
	}
	n.Dtype, n.Mountpoint = r[0], r[1]
	n.Fstype, n.Device = r[2], r[3]
	return n, nil
}

// ParseRawLsblk parses raw lsblk output and creates an Lsblk struct
func ParseRawLsblk(raw string) (Lsblk, error) {
	l := Lsblk{}
	l.Raw = raw
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		node, err := RawLineToNode(line)
		if err != nil {
			return l, fmt.Errorf("parse error. %s", err)
		}

		switch node.Dtype {
		case "disk":
			disk := Disk{}
			disk.Disk = node
			l.Disks = append(l.Disks, disk)
		default:
			last := len(l.Disks) - 1
			l.Disks[last].Parts = append(l.Disks[last].Parts, node)
		}
	}
	return l, nil
}

// GetLsblkInfo returns an Lsblk struct based on the current raw parsed lsblk output
func GetLsblkInfo() (Lsblk, error) {
	raw, err := ExecLsblk()
	if err != nil {
		return Lsblk{}, err
	}
	result, err := ParseRawLsblk(raw)
	if err != nil {
		return Lsblk{}, err
	}
	return result, nil
}
