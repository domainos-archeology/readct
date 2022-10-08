package cli

type Timestamp struct {
	Time  uint32 // Apollo time (seconds since 1/1/1980?)
	Extra uint32 // extra precision, plus sometimes a node ID
}

type BlockHeader struct {
	SequenceNumber int32
	Timestamp      Timestamp
	Size           uint16 // size of data within the block
}

type MagicHeader struct {
	Type1 int16  /* data type (see below) */
	Size  uint16 /* size of data in next section */
	Type2 int16  /* data type (see below) */
}

type ACL byte
type Magic int32
type Inode int64

func (h MagicHeader) Type() Magic {
	return makeMagic(int32(h.Type1), int32(h.Type2))
}

func makeMagic(t1, t2 int32) Magic {
	return Magic(t1<<16 | t2)
}

var (
	MagicSub     = makeMagic(9, 2)
	MagicSubOld  = makeMagic(9, 1)
	MagicEmpty   = makeMagic(8, 1)
	MagicOpt     = makeMagic(6, 2)
	MagicOptOld  = makeMagic(6, 1)
	MagicACL     = makeMagic(10, 2)
	MagicACLOld  = makeMagic(10, 1)
	MagicName    = makeMagic(2, 2)
	MagicNameOld = makeMagic(2, 1)
	MagicFile    = makeMagic(0, 2)
	MagicFileOld = makeMagic(0, 1)
	MagicData    = makeMagic(1, 1) /* file contents */
	MagicDir     = makeMagic(3, 3)
	MagicDirOld  = makeMagic(3, 2)
	MagicPopd    = makeMagic(4, 2)
	MagicPopdOld = makeMagic(4, 1)
	MagicLink    = makeMagic(5, 2)
	MagicLinkOld = makeMagic(5, 1)
)

func (m Magic) String() string {
	switch m {
	case MagicSub:
		return "SUB"
	case MagicSubOld:
		return "SUB_OLD"
	case MagicEmpty:
		return "EMPTY_MAGIC"
	case MagicOpt:
		return "OPT_MAGIC"
	case MagicOptOld:
		return "OPT_OLD_MAGIC"
	case MagicACL:
		return "ACL_MAGIC"
	case MagicACLOld:
		return "ACL_OLD_MAGIC"
	case MagicName:
		return "NAME_MAGIC"
	case MagicNameOld:
		return "NAME_OLD_MAGIC"
	case MagicFile:
		return "FILE_MAGIC"
	case MagicFileOld:
		return "FILE_OLD_MAGIC"
	case MagicData:
		return "DATA_MAGIC"
	case MagicDir:
		return "DIR_MAGIC"
	case MagicDirOld:
		return "DIR_OLD_MAGIC"
	case MagicPopd:
		return "POPD_MAGIC"
	case MagicPopdOld:
		return "POPD_OLD_MAGIC"
	case MagicLink:
		return "LINK_MAGIC"
	case MagicLinkOld:
		return "LINK_OLD_MAGIC"
	default:
		return "UNKNOWN"
	}
}

type SubHeader struct {
	Ignore1   int16 // unused
	Timestamp Timestamp
	Ignore2   int16 // unused
}

type NameHeader struct {
	Inode  Inode
	Ignore int32 // unused
	// name bytes follow directly after
}

type CommonHeader struct {
	Ignore1  int32
	Inode    Inode /* inode of this file/directory */ /* probably actually Apollo UID of the object */
	TypeHigh int32 /* type_uid.high */
	TypeLow  int32 /* type_uid.low */
	Size     int32 /* actual size */
	Ignore4  int32
	Mtime    Timestamp /* modification time */
	Itime1   Timestamp
	Itime2   Timestamp
	Itime3   Timestamp
	Dinode   Inode /* inode of parent directory */
	Ignore5  int32
	Itime4   Timestamp
	Itime5   Timestamp
	Ignore6  int32
	Ignore7  int32
	Uacl     ACL /* acl for owner */
	Gacl     ACL /* acl for group */
	Zacl     ACL /* acl for organization */
	Oacl     ACL /* acl for world */
	Ignore8  int32
	Uid      int32 /* user id */
	Gid      int32 /* group id */
	Oid      int32 /* organization id */
	Nlink    int16 /* number of hard links */
	Pad      int16
}

type CommonOldHeader struct {
	Ignore1 int32
	Inode   Inode /* inode of this file/directory */
	Ignore2 int32
	Ignore3 int32
	Iinode  Inode
	Size    int32 /* actual size */
	Ignore4 int32
	Atime   int32 /* access time, no extra prec */
	Mtime   int32 /* modification time, no extra prec */
	Dinode  Inode /* inode of parent directory */
	Ignore5 int16
	Pad     int16
}

type DirHeader struct {
	CommonHeader
	Ignore [32]int32
	// directory name bytes follow directly after
}

type DirOldHeader struct {
	CommonOldHeader
	Ignore [2]int32
	Inode1 Inode
	Inode2 Inode
	// directory name bytes follow directly after
}

type FileHeader struct {
	CommonHeader
	Ignore [6]int32
}

type FileOldHeader struct {
	CommonOldHeader
	Ignore [2]int32
}

type LinkHeader struct {
	Ignore     int16
	NameLength int32
	// link name bytes follow directly after (length == NameLength above)
	// link destination follow directly after
}

type UID struct {
	High int32
	Low  int32
}

func (u UID) Equal(o UID) bool {
	return u.High == o.High && u.Low == o.Low
}

var (
	UID_BITMAP     = UID{0x317, 0x0}
	UID_CASE_HM    = UID{0x316, 0x0}
	UID_CMPEXE     = UID{0x325, 0x0}
	UID_COFF       = UID{0x322, 0x0}
	UID_D3M_AREA   = UID{0x30e, 0x0}
	UID_D3M_SCH    = UID{0x30f, 0x0}
	UID_DDF        = UID{0x30b, 0x0}
	UID_DEV_TTY    = UID{0x324, 0x0}
	UID_DIR        = UID{0x312, 0x0}
	UID_DM_EDIT    = UID{0x31a, 0x0}
	UID_HDRU       = UID{0x301, 0x0}
	UID_IPAD       = UID{0x309, 0x0}
	UID_IPCSOCK    = UID{0x31f, 0x0}
	UID_LHEAP      = UID{0x319, 0x0}
	UID_MBX        = UID{0x30c, 0x0}
	UID_MT         = UID{0x314, 0x0}
	UID_NULL       = UID{0x30d, 0x0}
	UID_OBJ        = UID{0x302, 0x0}
	UID_OBJLIB     = UID{0x318, 0x0}
	UID_OSIO       = UID{0x326, 0x0}
	UID_OS_PG_FILE = UID{0x323, 0x0}
	UID_PAD        = UID{0x305, 0x0}
	UID_PIPE       = UID{0x310, 0x0}
	UID_PTY_SLAVE  = UID{0x31d, 0x0}
	UID_PTY        = UID{0x31c, 0x0}
	UID_REC        = UID{0x300, 0x0}
	UID_SIO        = UID{0x30a, 0x0}
	UID_SLINK      = UID{0x31e, 0x0}
	UID_SYSBOOT    = UID{0x315, 0x0}
	UID_TCP        = UID{0x31b, 0x0}
	UID_UASC       = UID{0x311, 0x0}
	UID_UNIX_DIR   = UID{0x313, 0x0}
	UID_UNSTRUCT   = UID{0x321, 0x0}
	UID_VTE        = UID{0x320, 0x0}
	UID_UNDEF      = UID{0x304, 0x0}
	UID_NIL        = UID{0x0, 0x0}
	UID_COMPRESS   = UID{0x40c9eb63, 0x40018ec0}
)

const (
	DESC_BITMAP    = "bitmap"
	DESC_CASE_HM   = "case_hm"
	DESC_CMPEXE    = "cmpexe"
	DESC_COFF      = "coff"
	DESC_D3M_AREA  = "d3m_area"
	DESC_D3M_SCH   = "d3m_sch"
	DESC_DEV_TTY   = "dev_tty"
	DESC_DIR       = "dir"
	DESC_DM_EDIT   = "dm_edit"
	DESC_HDRU      = "hdru"
	DESC_IPAD      = "ipad"
	DESC_MBX       = "mbx"
	DESC_MT        = "mt"
	DESC_NULL      = "null"
	DESC_OBJ       = "obj"
	DESC_OSIO      = "osio"
	DESC_PAD       = "pad"
	DESC_PIPE      = "pipe"
	DESC_PTY_SLAVE = "pty_slave"
	DESC_PTY       = "pty"
	DESC_REC       = "rec"
	DESC_SIO       = "sio"
	DESC_TCP       = "tcp"
	DESC_UASC      = "uasc"
	DESC_UNSTRUCT  = "unstruct"
	DESC_NIL       = "nil"
	DESC_COMPRESS  = "compress"
)
