package main

import (
	"flag"
	"fmt"

	"github.com/cheggaaa/pb/v3"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" || to == "" {
		fmt.Println("from and to arguments required")
		flag.Usage()
		return
	}

	bar := pb.Full.Start64(limit)
	ProgressBar = bar

	if err := Copy(from, to, offset, limit); err != nil {
		fmt.Println(err)
		return
	}
}
