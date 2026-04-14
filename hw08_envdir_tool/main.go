package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Invalid args count.")
		return
	}
	folder, cmd := processArgs(args)

	env, err := ReadDir(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Exit(RunCmd(cmd, env))
}

func processArgs(args []string) (string, []string) {
	return args[0], args[1:]
}
