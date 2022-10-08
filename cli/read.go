package cli

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

const BlockSize = 512

type NameMessage struct {
	Header NameHeader
	Name   string
}

type DirMessage struct {
	Header DirHeader
	Name   string
}

type DirOldMessage struct {
	Header DirOldHeader
	Name   string
}

type FileMessage struct {
	Header FileHeader
}

type FileOldMessage struct {
	Header FileOldHeader
}

type DataMessage struct {
	Data []byte
}

type LinkMessage struct {
	Header      LinkHeader
	Name        string
	Destination string
}

func ReadTapes(paths []string) chan any {
	ch := make(chan any)
	go func() {
		defer close(ch)
		for _, p := range paths {
			readTape(p, ch)
		}
	}()
	return ch
}

func findBlocks(file io.ReadSeeker) int64 {
	// this cannot possibly be the right way to check this..  the tapes either
	// have ctboot at the start (along with an rfc header), or they have a
	// (physical?) volume header.
	var blockOffset int64
	fileMagic := make([]byte, 3)
	err := binary.Read(file, binary.BigEndian, &fileMagic)
	if err != nil {
		panic("couldn't detect file magic")
	}

	if !bytes.Equal(fileMagic, []byte("VOL")) {
		// we need to skip past ctboot. I'm sure there's something in the RFC header
		// to tell us how big it is, but instead we just loop until we find a block
		// with 0xdeaffaed as the first 4 bytes in it
		for {
			blockOffset += BlockSize
			_, err := file.Seek(blockOffset, io.SeekStart)
			if err != nil {
				panic("couldn't seek to block")
			}
			blockMagic := make([]byte, 4)
			err = binary.Read(file, binary.BigEndian, &blockMagic)
			if err != nil {
				panic("couldn't read block magic")
			}
			if bytes.Equal(blockMagic, []byte{0xde, 0xaf, 0xfa, 0xed}) {
				blockOffset += BlockSize
				break
			}
		}

		// check to make sure that after ctboot, we have the volume header
		_, err := file.Seek(blockOffset, io.SeekStart)
		if err != nil {
			panic("couldn't seek to block")
		}
		err = binary.Read(file, binary.BigEndian, &fileMagic)
		if err != nil {
			panic("couldn't detect file magic")
		}

		if !bytes.Equal(fileMagic, []byte("VOL")) {
			panic("volume header not found after ctboot")
		}
	}

	for {
		blockOffset += BlockSize
		_, err := file.Seek(blockOffset, io.SeekStart)
		if err != nil {
			panic("couldn't seek to block")
		}
		blockMagic := make([]byte, 4)
		err = binary.Read(file, binary.BigEndian, &blockMagic)
		if err != nil {
			panic("couldn't read block magic")
		}
		if bytes.Equal(blockMagic, []byte{0xde, 0xaf, 0xfa, 0xed}) {
			return blockOffset + BlockSize
		}
	}
}

func readTape(path string, ch chan<- any) {
	fmt.Println("read_ct: reading from", path)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	offset := findBlocks(file)
	// fmt.Println("[DEBUG] initial blocks offset =", offset)
	off, err := file.Seek(offset, io.SeekStart)
	if err != nil {
		panic(err)
	}
	if off != offset {
		panic("seek failed")
	}

	for {
		// read a block at a time from the tape file
		buf := make([]byte, BlockSize)
		err = binary.Read(file, binary.BigEndian, &buf)
		if err == io.EOF {
			// we're done with this tape?
			return
		}
		if err != nil {
			panic(fmt.Sprintf("failed to read block: %v", err))
		}

		// then create a reader from the block that we'll use for the rest of
		// the processing
		block := bytes.NewReader(buf)

		var header BlockHeader
		err = binary.Read(block, binary.BigEndian, &header)
		if err != nil {
			panic("failed to read block header")
		}

		// fmt.Println("[DEBUG] block seq:", header.SequenceNumber)
		// fmt.Println("[DEBUG] block size:", header.Size)

	processBlock:
		for {
			// each block is made up of 1-N sections, that start with a magic header
			var magic MagicHeader
			err = binary.Read(block, binary.BigEndian, &magic)
			if err == io.EOF {
				break
			}
			if err != nil {
				// fmt.Printf("[DEBUG] failed to read magic header: %v\n", err)
				continue processBlock
			}
			// fmt.Println("[DEBUG] magic type:", magic.Type1, magic.Type2, magic.Type())
			// fmt.Println("[DEBUG] magic size:", magic.Size)

			magicContent := make([]byte, magic.Size)
			err = binary.Read(block, binary.BigEndian, &magicContent)
			if err != nil {
				// fmt.Printf("[DEBUG] failed to read magic content: %v\n", err)
				continue processBlock
			}
			// continue to ensure content is aligned to a 2-byte boundary
			if magic.Size%2 == 1 {
				_, _ = block.ReadByte()
			}

			magicContentReader := bytes.NewReader(magicContent)

			switch magic.Type() {
			case MagicName, MagicNameOld:
				var name NameHeader
				err = binary.Read(magicContentReader, binary.BigEndian, &name)
				if err == io.EOF {
					break
				}
				nameBytes, err := io.ReadAll(magicContentReader)
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("[DEBUG] failed to read file name: %v\n", err)
					continue processBlock
				}
				if err == io.EOF {
					break
				}

				if magic.Type() == MagicNameOld {
					fmt.Println("OLD NAME")
					nameBytes = []byte("broken")
				}

				ch <- NameMessage{
					Header: name,
					Name:   string(nameBytes),
				}

			case MagicDir:
				var dir DirHeader
				err = binary.Read(magicContentReader, binary.BigEndian, &dir)
				if err == io.EOF {
					break
				}
				nameBytes, err := io.ReadAll(magicContentReader)
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("[DEBUG] failed to read dir name: %v\n", err)
					continue processBlock
				}

				ch <- DirMessage{
					Header: dir,
					Name:   string(nameBytes),
				}

			case MagicDirOld:
				var dir DirOldHeader
				err = binary.Read(magicContentReader, binary.BigEndian, &dir)
				if err == io.EOF {
					break
				}
				nameBytes, err := io.ReadAll(magicContentReader)
				if err == io.EOF {
					break
				}
				if err != nil {
					// fmt.Printf("[DEBUG] failed to read dir name: %v\n", err)
					continue processBlock
				}

				ch <- DirOldMessage{
					Header: dir,
					Name:   strings.ToLower(string(nameBytes)),
				}

			case MagicFile:
				var fileH FileHeader
				err = binary.Read(magicContentReader, binary.BigEndian, &fileH)
				if err == io.EOF {
					break
				}
				ch <- FileMessage{
					Header: fileH,
				}

			case MagicFileOld:
				var fileH FileOldHeader
				err = binary.Read(magicContentReader, binary.BigEndian, &fileH)
				if err == io.EOF {
					break
				}
				ch <- FileOldMessage{
					Header: fileH,
				}

			case MagicData:
				data, err := io.ReadAll(magicContentReader)
				if err != nil {
					panic(fmt.Sprintf("failed data: %v\n", err))
				}

				ch <- DataMessage{
					Data: data,
				}

			case MagicPopd:
				fmt.Println("[DEBUG] *** popd found")

			case MagicLink, MagicLinkOld:
				var link LinkHeader
				err = binary.Read(magicContentReader, binary.BigEndian, &link)
				if err == io.EOF {
					break
				}

				nameBytes := make([]byte, link.NameLength)
				err = binary.Read(magicContentReader, binary.BigEndian, &nameBytes)
				if err != nil {
					fmt.Printf("failed to read link name: %v (name length is %d)\n", err, link.NameLength)
					break
				}

				destBytes, err := io.ReadAll(magicContentReader)
				if err != nil {
					fmt.Printf("failed to read link destination: %v\n", err)
					break
				}

				ch <- LinkMessage{
					Header:      link,
					Name:        string(nameBytes),
					Destination: string(destBytes),
				}

			case MagicSub, MagicSubOld:
				// I'm sure we should do something here, but skip it for now...

			case MagicEmpty:
				// not sure what the point of this one is...

			case MagicACL, MagicACLOld:
				// TBD

			default:
				// ignore everything else for now
			}
		}
	}
}
