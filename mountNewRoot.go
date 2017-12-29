package main

import (
	"syscall"
	"os"
	"path/filepath"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"io/ioutil"
)

type  mountRoot struct{
	oldPath string
}


func pivotRoot(rootfs, pivotBaseDir string) error {
	if pivotBaseDir == "" {
		pivotBaseDir = "/"
	}
	tmpDir := filepath.Join(rootfs, pivotBaseDir)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("can't create tmp dir %s, error %v", tmpDir, err)
	}
	pivotDir, err := ioutil.TempDir(tmpDir, ".pivot_root")
	if err != nil {
		return fmt.Errorf("can't create pivot_root dir %s, error %v", pivotDir, err)
	}
	if err := syscall.PivotRoot(rootfs, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %s", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %s", err)
	}
	// path to pivot dir now changed, update
	pivotDir = filepath.Join(pivotBaseDir, filepath.Base(pivotDir))

	// Make pivotDir rprivate to make sure any of the unmounts don't
	// propagate to parent.
	if err := syscall.Mount("", pivotDir, "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return err
	}

	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %s", err)
	}
	return os.Remove(pivotDir)
}


func (self *mountRoot) mountNewRoot() error {

	d,e := os.Getwd()

	if e != nil {
		return e;
	}

	log.Info("mount new root current dir "+d)

	if err := syscall.Mount(d, d,"bind", syscall.MS_BIND| syscall.MS_REC,"");err != nil {
		return err
	}

	pivotRoot(d,".pivotRoot")

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	if err:= syscall.Mount("tmpfs","/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "");err != nil {
		log.Error("Mount tmpfs" + err.Error())
		return err
	}

	return nil
}
