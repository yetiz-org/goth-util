package buf

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ByteBuf interface {
	io.Writer
	io.Reader
	ReaderIndex() int
	WriterIndex() int
	MarkReaderIndex() ByteBuf
	ResetReaderIndex() ByteBuf
	MarkWriterIndex() ByteBuf
	ResetWriterIndex() ByteBuf
	Reset() ByteBuf
	Bytes() []byte
	ReadableBytes() int
	Cap() int
	Grow(v int)
	WriteByte(c byte)
	WriteBytes(bs []byte)
	WriteString(s string)
	WriteByteBuf(buf ByteBuf)
	WriteReader(reader io.Reader)
	WriteInt16(v int16)
	WriteUInt16(v uint16)
	WriteInt32(v int32)
	WriteUInt32(v uint32)
	WriteInt64(v int64)
	WriteUInt64(v uint64)
	ReadByte() byte
	ReadBytes(len int) []byte
	ReadByteBuf(len int) ByteBuf
	ReadInt16() int16
	ReadUInt16() uint16
	ReadInt32() int32
	ReadUInt32() uint32
	ReadInt64() int64
	ReadUInt64() uint64
}

var ErrNilObject = fmt.Errorf("nil object")
var ErrInsufficientSize = fmt.Errorf("insufficient size")

func NewByteBuf(bs []byte) ByteBuf {
	buf := &DefaultByteBuf{}
	buf.WriteBytes(bs)
	return buf
}

func EmptyByteBuf() ByteBuf {
	return &DefaultByteBuf{}
}

type DefaultByteBuf struct {
	buf                                                        []byte
	readerIndex, writerIndex, prevReaderIndex, prevWriterIndex int
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

func (b *DefaultByteBuf) ReaderIndex() int {
	return b.readerIndex
}

func (b *DefaultByteBuf) WriterIndex() int {
	return b.writerIndex
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

func (b *DefaultByteBuf) Grow(v int) {
	tb := make([]byte, b.Cap()+v)
	if b.prevReaderIndex == 0 {
		offset := b.readerIndex - b.prevReaderIndex
		copy(tb, b.buf[b.readerIndex:])
		b.readerIndex -= offset
		b.writerIndex -= offset
		if b.prevWriterIndex > 0 {
			b.prevWriterIndex -= offset
		}
	} else {
		copy(tb, b.buf[b.prevReaderIndex:])
		offset := b.prevReaderIndex
		b.prevReaderIndex = 0
		b.readerIndex -= offset
		b.writerIndex -= offset
		if b.prevWriterIndex > 0 {
			b.prevWriterIndex -= offset
		}
	}

	b.buf = tb
}

func (b *DefaultByteBuf) WriteByte(c byte) {
	b.prepare(1)
	b.buf[b.writerIndex] = c
	b.writerIndex++
}

func (b *DefaultByteBuf) WriteBytes(bs []byte) {
	pl := len(bs)
	b.prepare(pl)
	copy(b.buf[b.writerIndex:], bs)
	b.writerIndex += pl
}

func (b *DefaultByteBuf) WriteByteBuf(buf ByteBuf) {
	if buf == nil {
		panic(ErrNilObject)
	}

	b.WriteBytes(buf.Bytes())
}

func (b *DefaultByteBuf) WriteReader(reader io.Reader) {
	if reader == nil {
		panic(ErrNilObject)
	}

	if bs, err := io.ReadAll(reader); err != nil {
		panic(err)
	} else {
		b.WriteBytes(bs)
	}
}

func (b *DefaultByteBuf) WriteString(s string) {
	b.WriteBytes([]byte(s))
}

func (b *DefaultByteBuf) WriteInt16(v int16) {
	b.WriteUInt16(uint16(v))
}

func (b *DefaultByteBuf) WriteUInt16(v uint16) {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, v)
	b.prepare(len(bs))
	b.WriteBytes(bs)
}

func (b *DefaultByteBuf) WriteInt32(v int32) {
	b.WriteUInt32(uint32(v))
}

func (b *DefaultByteBuf) WriteUInt32(v uint32) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, v)
	b.prepare(len(bs))
	b.WriteBytes(bs)
}

func (b *DefaultByteBuf) WriteInt64(v int64) {
	b.WriteUInt64(uint64(v))
}

func (b *DefaultByteBuf) WriteUInt64(v uint64) {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, v)
	b.prepare(len(bs))
	b.WriteBytes(bs)
}

func (b *DefaultByteBuf) ReadByte() byte {
	if b.readerIndex == b.writerIndex {
		panic(ErrInsufficientSize)
	}

	b.readerIndex++
	return b.buf[b.readerIndex-1]
}

func (b *DefaultByteBuf) ReadBytes(len int) []byte {
	if b.ReadableBytes() < len {
		panic(ErrInsufficientSize)
	}

	b.readerIndex += len
	return b.buf[b.readerIndex-len : b.readerIndex]
}

func (b *DefaultByteBuf) ReadByteBuf(len int) ByteBuf {
	buf := &DefaultByteBuf{}
	buf.WriteBytes(b.ReadBytes(len))
	return buf
}

func (b *DefaultByteBuf) ReadInt16() int16 {
	return int16(b.ReadUInt16())
}

func (b *DefaultByteBuf) ReadUInt16() uint16 {
	return binary.BigEndian.Uint16(b.ReadBytes(2))
}

func (b *DefaultByteBuf) ReadInt32() int32 {
	return int32(b.ReadUInt32())
}

func (b *DefaultByteBuf) ReadUInt32() uint32 {
	return binary.BigEndian.Uint32(b.ReadBytes(4))
}

func (b *DefaultByteBuf) ReadInt64() int64 {
	return int64(b.ReadUInt64())
}

func (b *DefaultByteBuf) ReadUInt64() uint64 {
	return binary.BigEndian.Uint64(b.ReadBytes(8))
}

func (b *DefaultByteBuf) prepare(i int) {
	if b.Cap() == 0 {
		b.Grow(32)
	}

	for ; b.writerIndex+i > b.Cap(); {
		b.Grow(b.Cap())
	}
}
