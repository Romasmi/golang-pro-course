package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	const output = "Hello, OTUS!"
	fmt.Println(reverse.String(output))
}
