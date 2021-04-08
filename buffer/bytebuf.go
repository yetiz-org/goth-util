package buffer

type ByteBuf struct {
	buf                   []byte
	readIndex, writeIndex int
}
