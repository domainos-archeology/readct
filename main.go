package main

import (
	"os"

	"github.com/domainos-archeology/readct/cli"
)

func main() {
	cli.ReadTape(os.Args[1])
}
