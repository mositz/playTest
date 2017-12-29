package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"syscall"
	"os/exec"
)
func main(){
	app := cli.NewApp()
	mountRoot := &mountRoot{
		oldPath:"/root/busybox",
	}
	config := &config{}
	app.Commands = []cli.Command{
		{
			Name:"run",
			Flags:[]cli.Flag{
				cli.BoolFlag{
					Name:"ti",
				},
				cli.StringFlag{
					Name:"m",
				},
				cli.StringFlag{
					Name: "cpushare",
				},
				cli.StringFlag{
					Name: "cpuset",
				},
			},
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
				if context.Bool("ti") {
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
				}
				cmd.Dir = mountRoot.oldPath

				if err:=cmd.Start();err != nil {
					log.Error( err.Error())
					return
				}
				subs := []subsystem{
					{
						mountInfo:"/proc/self/mountinfo",
						name:"memory",
						underPath:"testCgroup",
						cgroupRelativeFile:"memory.limit_in_bytes",
						thresholdValue: config.getDefault(memory,context.String("m")).(string),
					},
					{
						mountInfo:"/proc/self/mountinfo",
						name:"cpu",
						underPath:"testCgroup",
						cgroupRelativeFile:"cpu.shares",
						thresholdValue:config.getDefault(cpu,  context.String("cpushare")).(string),
					},
					{
						mountInfo:"/proc/self/mountinfo",
						name:"cpuset",
						underPath:"testCgroup",
						cgroupRelativeFile:"cpuset.cpus",
						thresholdValue:config.getDefault(cpuset,  context.String("cpuset")).(string),
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

				if err := mountRoot.mountNewRoot();err != nil{
					log.Error(err.Error())
				}
				log.Info("mount finish" + commond)

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







