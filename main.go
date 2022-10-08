package main

import (
	"os"

	"github.com/domainos-archeology/readct/cli"
)

func main() {
	if os.Args[1] == "x" {
		cli.Extract(os.Args[2:])
	} else {
		cli.Index(os.Args[1:])
	}
}
