package main

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"bufio"
	"strings"
	"path"
	"syscall"
	"fmt"
)

type subsystem struct {
	mountInfo string
	name string
	underPath string
    cgroupRelativeFile string
	thresholdValue string
	cgroupPath string
}
//type ErrorF func() string
//func (self ErrorF) Error() string { return self() }

func (u *subsystem) set() error{
	file, err := os.Open(u.mountInfo)
	if(err != nil){
		log.Error(err.Error())
		return err
	}
	scan := bufio.NewScanner(file)
	for scan.Scan(){
		text := scan.Text()
		line := strings.Split(text, ",")
		if line[len(line)-1] == u.name {
			preLine :=strings.Split(line[0]," ")
			 cg := path.Join(preLine[4],u.underPath)
			 if _ ,err :=  os.Stat(cg); err != nil||os.IsNotExist(err) {
				 if err := os.Mkdir(cg,755); err ==nil {
					 fileName := path.Join(cg , u.cgroupRelativeFile)
					 u.cgroupPath = cg
					 log.Error("11" + path.Join(u.cgroupPath, "tasks")+":" + u.cgroupPath)
					 return ioutil.WriteFile(fileName, []byte(u.thresholdValue), 775)
				 }else{
				 	return err;
				 }
			 }else{
				 u.cgroupPath = cg
				 log.Error("11" + path.Join(u.cgroupPath, "tasks")+":" + u.cgroupPath)
			 }
		}
	}
	return nil;
}

func (u *subsystem) apply(pid int) error{
	if u.cgroupPath == "" {
		log.Error(path.Join(u.cgroupPath, "tasks")+":" + u.cgroupPath)
		return fmt.Errorf("%s is null", u.cgroupPath);
	}
	log.Info(path.Join(u.cgroupPath, "tasks"),fmt.Sprintf("%d",pid))
	return ioutil.WriteFile(path.Join(u.cgroupPath, "tasks"), []byte(fmt.Sprintf("%d",pid)), 775)
}

func (u *subsystem) delete() error{
	if len(u.cgroupPath) < 1 {
		return fmt.Errorf("%s is null", u.cgroupPath);
	}
	return syscall.Rmdir(u.cgroupPath)
}

