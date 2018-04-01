package fs

import (
	"syscall"
	"os"
	"unsafe"
	"C"
)

const (
	// From https://android.googlesource.com/platform/system/sepolicy/+/nougat-dev/ioctl_defines
	fiFreeze = uintptr(0xc0045877)
	fiThaw   = uintptr(0xc0045878)
)

func ioctl(fd, cmd, ptr uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if e != 0 {
		return e
	}
	return nil
}

// Freeze sends the fifreeze ioctl to the fd of mountpoint
func Freeze(mountpoint string) error {
	f, err := os.Open(mountpoint)
	if err != nil {
		return err
	}
	err = fifreeze(f.Fd())
	if err != nil {
		return err
	}
	_ = f.Close()
	return nil
}


// Unfreeze sends the fithaw ioctl to the fd of mountpoint
func Unfreeze(mountpoint string) error {
	f, err := os.Open(mountpoint)
	if err != nil {
		return err
	}
	err = fithaw(f.Fd())
	if err != nil {
		return err
	}
	_ = f.Close()
	return nil
}

func fifreeze(fd uintptr) error {
	var n C.uint
	return ioctl(fd, fiFreeze, uintptr(unsafe.Pointer(&n)))
}

func fithaw(fd uintptr) error {
	var n C.uint
	return ioctl(fd, fiThaw, uintptr(unsafe.Pointer(&n)))
}
