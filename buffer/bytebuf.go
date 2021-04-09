package buffer

import (
	"fmt"
	"io"
)

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
	ReadableBytes() int
	Cap() int
	Grow(i int)
	WriteBytes(bs []byte)
	WriteByteBuf(buf ByteBuf)
	WriteInt(i int)
	WriteUInt(i uint)
	WriteInt64(i int64)
	WriteUInt64(i uint64)
	ReadBytes(len int) ([]byte, )
	ReadByteBuf(len int) (ByteBuf, )
	ReadInt() (int, )
	ReadUInt() (uint, )
	ReadInt64() (int64, )
	ReadUInt64() (uint64, )
}

var NilObject = fmt.Errorf("nil object")
var OutOfRange = fmt.Errorf("out of range")
var InsufficientSize = fmt.Errorf("insufficient size")

//
//func NewByteBuf() ByteBuf {
//	return &*DefaultByteBuf{}
//}

func NewByteBuf(bs []byte) ByteBuf {
	buf := &DefaultByteBuf{}
	buf.WriteBytes(bs)
	return buf
}

type DefaultByteBuf struct {
	buf                                                        []byte
	readerIndex, writerIndex, prevReaderIndex, prevWriterIndex int
}

func (b *DefaultByteBuf) WriteByte(c byte) error {
	b.prepare(1)
	b.buf[b.writerIndex] = c
	b.writerIndex++
	return nil
}

func (b *DefaultByteBuf) Write(p []byte) (n int, err error) {
	pl := len(p)
	if pl == 0 {
		return 0, nil
	}

	b.prepare(pl)
	copy(b.buf[b.writerIndex:], p)
	b.writerIndex += pl
	return pl, nil
}

func (b *DefaultByteBuf) ReadByte() (byte, error) {
	if b.readerIndex == b.writerIndex {
		return 0, OutOfRange
	}

	b.readerIndex++
	return b.buf[b.readerIndex-1], nil
}

func (b *DefaultByteBuf) Read(p []byte) (n int, err error) {
	cpLen := b.ReadableBytes()
	if cpLen == 0 {
		return 0, io.EOF
	}

	if cpLen > len(p) {
		cpLen = len(p)
	}

	copy(p, b.buf[b.readerIndex:b.readerIndex+cpLen])
	b.readerIndex += cpLen
	return cpLen, nil
}

func (b *DefaultByteBuf) WriteString(s string) (n int, err error) {
	return b.Write([]byte(s))
}

func (b *DefaultByteBuf) MarkReaderIndex() ByteBuf {
	b.prevReaderIndex = b.readerIndex
	return b
}

func (b *DefaultByteBuf) ResetReaderIndex() ByteBuf {
	b.readerIndex = b.prevReaderIndex
	b.prevReaderIndex = 0
	return b
}

func (b *DefaultByteBuf) MarkWriterIndex() ByteBuf {
	b.prevWriterIndex = b.writerIndex
	return b
}

func (b *DefaultByteBuf) ResetWriterIndex() ByteBuf {
	b.writerIndex = b.prevWriterIndex
	b.prevWriterIndex = 0
	return b
}

func (b *DefaultByteBuf) Reset() ByteBuf {
	b.buf = b.buf[:0]
	b.readerIndex = 0
	b.writerIndex = 0
	b.prevReaderIndex = 0
	b.prevWriterIndex = 0
	return b
}

func (b *DefaultByteBuf) Bytes() []byte {
	return b.buf[b.readerIndex:b.writerIndex]
}

func (b *DefaultByteBuf) ReadableBytes() int {
	return b.writerIndex - b.readerIndex
}

func (b *DefaultByteBuf) Cap() int {
	return len(b.buf)
}

func (b *DefaultByteBuf) Grow(i int) {
	tb := make([]byte, b.Cap()+i)
	if b.prevReaderIndex == 0 {
		copy(tb, b.buf[b.readerIndex:])
	} else {
		copy(tb, b.buf[b.prevReaderIndex:])
	}
}

func (b *DefaultByteBuf) WriteBytes(bs []byte) {
	pl := len(bs)
	b.prepare(pl)
	copy(b.buf[b.writerIndex:], bs)
	b.writerIndex += pl
}

func (b *DefaultByteBuf) WriteByteBuf(buf ByteBuf) {
	if buf == nil {
		panic(NilObject)
	}

	b.WriteBytes(buf.Bytes())
}

func (b *DefaultByteBuf) WriteInt(i int) {
	panic("implement me")
}

func (b *DefaultByteBuf) WriteUInt(i uint) {
	panic("implement me")
}

func (b *DefaultByteBuf) WriteInt64(i int64) {
	panic("implement me")
}

func (b *DefaultByteBuf) WriteUInt64(i uint64) {
	panic("implement me")
}

func (b *DefaultByteBuf) ReadBytes(len int) []byte {
	if b.ReadableBytes() < len {
		panic(InsufficientSize)
	}

	b.readerIndex += len
	return b.buf[b.readerIndex-len : b.readerIndex]
}

func (b *DefaultByteBuf) ReadByteBuf(len int) ByteBuf {
	buf := &DefaultByteBuf{}
	buf.WriteBytes(b.ReadBytes(len))
	return buf
}

func (b *DefaultByteBuf) ReadInt() int {
	panic("implement me")
}

func (b *DefaultByteBuf) ReadUInt() uint {
	panic("implement me")
}

func (b *DefaultByteBuf) ReadInt64() int64 {
	panic("implement me")
}

func (b *DefaultByteBuf) ReadUInt64() uint64 {
	panic("implement me")
}

func (b *DefaultByteBuf) prepare(i int) {
	if b.writerIndex+i < b.Cap() {
		b.Grow(b.Cap() * 2)
	}
}
