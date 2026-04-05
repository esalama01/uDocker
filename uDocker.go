package main

import (
	"fmt"
	"os"
	"uDocker/src"
)

func main() {
	args := os.Args
	switch args[1] {
	case "run": //to run you need to grant the binary root access.
		src.Run()
	case "child":
		src.Child()
	default:
		fmt.Println("Invalid Command")
	}
}
