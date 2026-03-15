package main

import (
	"fmt"
	"os"
	"uDocker/src"
)

func main() {
	args := os.Args
	switch args[1] {
	case "run":
		var arguments []string
		for _, arg := range args[3:] {
			arguments = append(arguments, arg)
		}
		output, _ := src.RunCommand(args[2], arguments...)
		fmt.Println(output)
	default:
		fmt.Println("hola")
	}
}
