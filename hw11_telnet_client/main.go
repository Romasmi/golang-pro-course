package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Romasmi/golang-pro-course/hw11_telnet_client/internal/validator"
)

type cliArgs struct {
	timeout *time.Duration
	host    string
	port    int
}

func main() {
	args, err := getCliArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewTelnetClient(fmt.Sprintf("%s:%d", args.host, args.port), *args.timeout, os.Stdin, os.Stdout)
	defer client.Close()
	err = client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(args)
}

func getCliArgs() (*cliArgs, error) {
	timeout := flag.Duration("timeout", 10*time.Second, "enter timeout like --timeout=10s, --timeout=1m etc")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid args count. Usage: go-telnet [--timeout=1m] <host> <port>")
	}

	host := args[0]
	if err := validator.ValidateHost(host); err != nil {
		return nil, err
	}
	// get port
	port, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	if err := validator.ValidatePort(port); err != nil {
		return nil, err
	}
	return &cliArgs{
		timeout: timeout,
		host:    host,
		port:    port,
	}, nil
}
