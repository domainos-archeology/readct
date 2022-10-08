package cli

import "fmt"

func Index(paths []string) {
	ch := ReadTapes(paths)

	var fileName string

	for msg := range ch {
		switch m := msg.(type) {

		case NameMessage:
			fileName = m.Name

		case FileMessage:
			fmt.Printf("(file) %s  (%s %d)\n", fileName, typeName(UID{m.Header.TypeHigh, m.Header.TypeLow}), m.Header.Size)

		case DirMessage:
			fmt.Printf("(dir) %s\n", m.Name)

		case DirOldMessage:
			fmt.Printf("(dir) %s\n", m.Name)

		case LinkMessage:
			fmt.Printf("(link) %s -> %s\n", m.Name, m.Destination)

		default:
			// do nothing
		}
	}
}

func typeName(uid UID) string {
	for _, u := range uid_table {
		if u.uid.Equal(uid) {
			return u.desc
		}
	}
	return "unknown"
}

var uid_table = []struct {
	uid        UID
	desc       string
	headerSize int
}{
	/* Most common first as this is sequentially searched */
	/*  uid value       description     header size */
	{UID_UNSTRUCT, DESC_UNSTRUCT, 0},
	{UID_UASC, DESC_UASC, 32},
	{UID_REC, DESC_REC, 32},
	{UID_HDRU, DESC_HDRU, 32},
	{UID_COFF, DESC_COFF, 0},
	{UID_OBJ, DESC_OBJ, 0},
	{UID_NIL, DESC_NIL, 0},
	{UID_COMPRESS, DESC_COMPRESS, 0},
	{UID_BITMAP, DESC_BITMAP, 0},
	{UID_CASE_HM, DESC_CASE_HM, 0},
	{UID_CMPEXE, DESC_CMPEXE, 0},
	{UID_D3M_AREA, DESC_D3M_AREA, 0},
	{UID_D3M_SCH, DESC_D3M_SCH, 0},
	{UID_DIR, DESC_DIR, 0},
	{UID_DM_EDIT, DESC_DM_EDIT, 0},
	{UID_IPAD, DESC_IPAD, 0},
	{UID_MBX, DESC_MBX, 0},
	{UID_MT, DESC_MT, 0},
	{UID_NULL, DESC_NULL, 0},
	{UID_PAD, DESC_PAD, 0},
	{UID_PIPE, DESC_PIPE, 0},
	{UID_PTY_SLAVE, DESC_PTY_SLAVE, 0},
	{UID_PTY, DESC_PTY, 0},
	{UID_SIO, DESC_SIO, 0},
	{UID_TCP, DESC_TCP, 0},
	{UID_OSIO, DESC_OSIO, 0},
	{UID_DEV_TTY, DESC_DEV_TTY, 0},
}
