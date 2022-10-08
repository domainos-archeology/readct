package main

import (
	"os"

	"github.com/domainos-archeology/read_ct_go/read_ct"
)

func main() {
	read_ct.ReadTape(os.Args[1])
}
