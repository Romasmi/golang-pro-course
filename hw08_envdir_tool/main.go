package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Invalid args count.")
	}
	folder, cmd, cmdArgs := processArgs(args)
	fmt.Printf("Folder: %s, Command: %s, Args: %v\n", folder, cmd, cmdArgs)
}

func processArgs(args []string) (string, string, []string) {
	return args[0], args[1], args[2:]
}
