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
	case "miaw":
		m := src.Manifest("hello-world", "latest")
		fmt.Print(m)
	//*/
	default:
		fmt.Println("Invalid Command")
	}
}
