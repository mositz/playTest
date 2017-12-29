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
	"strconv"
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
			 u.cgroupPath = cg
			 fileName := path.Join(cg , u.cgroupRelativeFile)
			 _ ,err :=  os.Stat(cg)
			 if err != nil||os.IsNotExist(err) {
				 if err := os.Mkdir(cg,0755); err ==nil {
					 log.Info(fmt.Sprintf("create new  file %s and write thresholdValue %s", fileName, u.thresholdValue))
					 return ioutil.WriteFile(fileName, []byte(u.thresholdValue), 0775)
				 }else{
				 	return err;
				 }
			 }
			 if err == nil {
				 log.Info(fmt.Sprintf("write thresholdValue %s abount fileName %s",u.thresholdValue, fileName))
				 return ioutil.WriteFile(fileName, []byte(u.thresholdValue), 0775)
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
	return ioutil.WriteFile(path.Join(u.cgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0775)
}

func (u *subsystem) delete() error{
	if len(u.cgroupPath) < 1 {
		return fmt.Errorf("%s is null", u.cgroupPath);
	}
	return syscall.Rmdir(u.cgroupPath)
}

