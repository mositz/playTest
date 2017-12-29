package main

import "strconv"

type config struct {}

const  (
	memory = "memory"
	cpu = "cpu"
	cpuset = "cpuset"
)

func (self *config) getDefault(name string,val string) interface{}{
	if(val == "") {
		switch name {
			case "memory":
				return strconv.Itoa(1024*1024)
			case "cpu":
				return "0-1"
			case "cpuset":
				return "1024"
		default:
			return val
		}
	}
	return val
}