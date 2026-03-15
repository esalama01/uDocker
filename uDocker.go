package main

import(
	"fmt"
	"os"
)

func main(){
	args := os.Args
	switch args[1]{
	case "run":
		//iterate through args[3:] and passs them toa  slice of strings, and then pass the slice as func param.
		//it must be RunCommand(args[2], the slice)
	}
default:
	fmt.Println("hola")
}