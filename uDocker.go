package main

import (
	"fmt"
	"os"
	"uDocker/src"
	//"runtime"
)

func main() {
	args := os.Args
	switch args[1] {
	case "run": //to run you you to grant root permissions.
		src.Run()
		
		//fmt.Printf("%v\n",runtime.NumCPU())
	case "child":
		src.Child()
	case "miaw":
		src.Configure_cgroups()
	default:	
		fmt.Println("Invalid Command")
	}
}
