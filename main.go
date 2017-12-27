package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"syscall"
	"os/exec"
	"fmt"
)
func main(){
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:"run",
			Action: func(context *cli.Context) {
				initCmd, err := os.Readlink("/proc/self/exe")
				if err != nil {
					log.Errorf("get init process error %v", err)
					return
				}
				commond := context.Args().Get(0)
				log.Infof("run command %s", commond)
				commonds := []string{"init",commond}
				cmd := exec.Command(initCmd, commonds...)
			 	cmd.SysProcAttr = 	&syscall.SysProcAttr{
					Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
						syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
				}
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err:=cmd.Start();err != nil {
					log.Error("%s",err)
					return
				}

				subs := []subsystem{
					{
						mountInfo:"/proc/self/mountinfo",
						name:"memory",
						underPath:"testCgroup",
						cgroupRelativeFile:"memory.limit_in_bytes",
						thresholdValue:fmt.Sprintf("%d",1024*1024*100),
					},
					{
						mountInfo:"/proc/self/mountinfo",
						name:"cpu",
						underPath:"testCgroup",
						cgroupRelativeFile:"cpu.shares",
						thresholdValue:fmt.Sprintf("%d",512),
					},
					{
						mountInfo:"/proc/self/mountinfo",
						name:"cpuset",
						underPath:"testCgroup",
						cgroupRelativeFile:"cpuset.cpus",
						thresholdValue:fmt.Sprintf("%s","0-1"),
					},
				}

				for _,u:= range subs {
					uu := &u
					uu.set()
					uu.apply(cmd.Process.Pid)
				}
				cmd.Wait()
				os.Exit(-1)
			},
		},
		{
			Name:"init",
			Action: func(context *cli.Context) {
				commond := context.Args().Get(0)
				log.Infof("init command %s", commond)

				defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
				syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
				argv := []string{commond}
				syscall.Exec(commond, argv,os.Environ())
			},
		},
	}

	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	app.Run(os.Args)
}







