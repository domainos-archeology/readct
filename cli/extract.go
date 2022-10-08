package cli

import (
	"fmt"
	"os"
	"path"
)

var fileName string
var file *os.File
var remainingFileSize int

func Extract(paths []string) {
	ch := ReadTapes(paths)
	for msg := range ch {
		switch m := msg.(type) {

		case DirMessage:
			fmt.Printf("(dir) %s\n", m.Name)
			if m.Name[0] == '/' || m.Name[0] == '.' {
				panic("bad dir name")
			}
			err := os.MkdirAll(m.Name, 0755)
			if err != nil {
				panic(fmt.Sprintf("failed to make directory structure %v", err))
			}

		case DirOldMessage:
			fmt.Printf("(dir) %s\n", m.Name)
			if m.Name[0] == '/' || m.Name[0] == '.' {
				panic("bad dir name")
			}
			err := os.MkdirAll(m.Name, 0755)
			if err != nil {
				panic(fmt.Sprintf("failed to make directory structure %v", err))
			}

		case NameMessage:
			fileName = m.Name

		case FileMessage:
			if fileName == "" {
				panic("no active file...")
			}
			fmt.Printf("(file) %s  (%s %d)\n", fileName, typeName(UID{m.Header.TypeHigh, m.Header.TypeLow}), m.Header.Size)
			remainingFileSize = int(m.Header.Size)
			var err error
			// why do we have to create the dir here?  shouldn't the DirMessages already have done so?
			err = os.MkdirAll(path.Dir(fileName), 0755)
			if err != nil {
				panic(fmt.Sprintf("failed to create containing directory %v", err))
			}
			file, err = os.Create(fileName)
			if err != nil {
				panic(fmt.Sprintf("failed to create file %v", err))
			}

		case FileOldMessage:
			if fileName == "" {
				panic("no active file...")
			}
			fmt.Printf("(file) %s  (unknown %d)\n", fileName, m.Header.Size)
			remainingFileSize = int(m.Header.Size)
			var err error
			// why do we have to create the dir here?  shouldn't the DirMessages already have done so?
			err = os.MkdirAll(path.Dir(fileName), 0755)
			if err != nil {
				panic(fmt.Sprintf("failed to create containing directory %v", err))
			}

			file, err = os.Create(fileName)
			if err != nil {
				panic(fmt.Sprintf("failed to create file %v", err))
			}

		case DataMessage:
			if file == nil {
				// fmt.Println("no active file, what is this data for?")
				break
			}
			if remainingFileSize == 0 {
				panic("already completed writing file")
			}
			toWrite := len(m.Data)
			if remainingFileSize < toWrite {
				toWrite = remainingFileSize
			}
			// write the data to the file
			n, err := file.Write(m.Data[:toWrite])
			if err != nil {
				panic(fmt.Sprintf("failed to write to file %v", err))
			}
			if n != toWrite {
				panic(fmt.Sprintf("failed to write all data to file %v", err))
			}

			remainingFileSize -= toWrite

			// decrement the written size from remainingFileSize
			if remainingFileSize == 0 {
				file.Close()
				file = nil
				fileName = ""
			}

		case LinkMessage:
			fmt.Printf("(link) %s -> %s\n", m.Name, m.Destination)
			// why do we have to create the dir here?  shouldn't the DirMessages already have done so?
			err := os.MkdirAll(path.Dir(m.Name), 0755)
			if err != nil {
				panic(fmt.Sprintf("failed to create containing directory %v", err))
			}
			os.Symlink(m.Destination, m.Name)

		default:
			// do nothing
		}
	}
}
