package main

import (
	"fmt"
	//"io"
	"os"
	"uDocker/src"
	//"runtime"
	//"net/http"
)

func main() {
	args := os.Args
	switch args[1] {
	case "run": //to run you you to grant root permissions.
		src.Run()

		//fmt.Printf("%v\n",runtime.NumCPU())
	case "child":
		src.Child()
	///*
	case "pull":
		src.Pull(args[2], args[3])
		//*/
	default:
		fmt.Println("Invalid Command")
	}
}

//// now i need to convert a list of layers into an actual filesys
