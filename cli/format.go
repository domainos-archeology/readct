package read_ct

type Timestamp struct {
	Time  int32 // Apollo time (seconds since 1/1/1980?)
	Extra int32 // extra precision, plus sometimes a node ID
}

type BlockHeader struct {
	SequenceNumber int32
	Timestamp      Timestamp
	Size           int16 // size of data within the block
}
