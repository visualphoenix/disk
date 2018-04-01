package lvm

import "os/exec"

// Suspend sends a dmsetup suspend for a given block device
func Suspend(blockDevice string) error {
	_, err := exec.Command("dmsetup", "suspend", blockDevice).Output()
	return err
}

// Resume sends a dmsetup resume for a given block device
func Resume(blockDevice string) error {
	_, err := exec.Command("dmsetup", "resume", blockDevice).Output()
	return err
}

// IsSuspended queries dmsetup for the suspended status of a block device
func IsSuspended(blockDevice string) (bool, error) {
	out, err := exec.Command("dmsetup", "info", "--noheadings", "-Co", "suspended", blockDevice).Output()
	if err != nil {
		return false, err
	}
	return string(out) == "Suspended", nil
}
