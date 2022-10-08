package read_ct

import (
	"fmt"
	"os"
)

func ReadTape() {
	fmt.Println("read_ct: reading from", os.Args[1])

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	off, err := file.Seek(0xc00, os.SEEK_SET)
	if err != nil {
		panic(err)
	}
	if off != 0xc00 {
		panic("seek failed")
	}
}
