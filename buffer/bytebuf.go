package buffer

import "io"

type ByteBuf interface {
	io.ByteWriter
	io.Writer
	io.ByteReader
	io.Reader
	io.StringWriter
	MarkReaderIndex() ByteBuf
	ResetReaderIndex() ByteBuf
	MarkWriterIndex() ByteBuf
	ResetWriterIndex() ByteBuf
	Reset() ByteBuf
	Bytes() []byte
	Len() int
	Cap() int
	Grow(i int)
	WriteBytes(bs []byte) error
	WriteByteBuf(buf ByteBuf) error
	WriteInt(i int) error
	WriteUInt(i uint) error
	WriteInt64(i int64) error
	WriteUInt64(i uint64) error
	ReadBytes(len int) ([]byte, error)
	ReadByteBuf(len int) (ByteBuf, error)
	ReadInt() (int, error)
	ReadUInt() (uint, error)
	ReadInt64() (int64, error)
	ReadUInt64() (uint64, error)
}

type DefaultByteBuf struct {
	buf                                                            []byte
	readerIndex, writerIndex, markedReaderIndex, markedWriterIndex int
}
//
//func NewByteBuf() ByteBuf {
//	return &DefaultByteBuf{}
//}
