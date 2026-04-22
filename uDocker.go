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
		m := src.Manifest_sha("hello-world", "latest")
		src.Pull_layers("hello-world", m)
		fmt.Print(m)
		fmt.Print("\n")
	case "miaaaw":
		if err := src.ExtractTarGz("layer.tar.gz", "./output"); err != nil {
			panic(err)
		}
	//*/
	default:
		fmt.Println("Invalid Command")
	}
}

//// now i need to convert a list of layers into an actual filesys
