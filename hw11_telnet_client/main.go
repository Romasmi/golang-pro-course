package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Romasmi/golang-pro-course/hw11_telnet_client/internal/validator" //nolint:depguard
)

type cliArgs struct {
	timeout *time.Duration
	host    string
	port    int
}

func main() {
	args, err := getCliArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args *cliArgs) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	address := net.JoinHostPort(args.host, strconv.Itoa(args.port))
	client := NewTelnetClient(address, *args.timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Close()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	type result struct {
		source string
		err    error
	}
	resCh := make(chan result, 2)

	go func() {
		resCh <- result{source: "send", err: client.Send()}
	}()

	go func() {
		resCh <- result{source: "receive", err: client.Receive()}
	}()

	select {
	case <-ctx.Done():
		return nil
	case res := <-resCh:
		if res.err != nil && !errors.Is(res.err, io.EOF) && !errors.Is(res.err, net.ErrClosed) {
			return res.err
		}
		if res.source == "send" {
			fmt.Fprintln(os.Stderr, "...EOF")
		} else {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		}
	}
	return nil
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
