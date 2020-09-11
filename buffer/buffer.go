package buffer

import (
	"io"
)

const (
	defaultBufferSize      int = 4096
	defaultBufferLowestCap int = 512
)

type Buffer struct {
	w         int
	r         int
	buff      []byte
	lowestCap int
}

func NewBufferSize(size, lowestCap int) *Buffer {
	return &Buffer{w: 0, r: 0, buff: make([]byte, size), lowestCap: lowestCap}
}

func NewBuffer() *Buffer {
	return NewBufferSize(defaultBufferSize, defaultBufferLowestCap)
}

func (b *Buffer) Size() int {
	return b.w - b.r
}

func (b *Buffer) Reset() {
	b.w, b.r = 0, 0
}

func (b *Buffer) allocate(need int) {
	if free := cap(b.buff) - b.w; free < need {
		newBuff := b.buff
		if b.r+free < need {
			newBuff = make([]byte, cap(b.buff)*2+need)
		}
		copy(newBuff, b.buff[b.r:b.w])
		b.buff = newBuff
		b.w = b.w - b.r
		b.r = 0
	}
}

func (b *Buffer) ReadFrom(reader io.Reader) (int, error) {
	if b.r >= b.w {
		b.Reset()
	}

	b.allocate(b.lowestCap)
	n, err := reader.Read(b.buff[b.w:])
	if n > 0 {
		b.w += n
	}
	return n, err
}

func (b *Buffer) Peek(n int) []byte {
	if n > b.Size() {
		n = b.Size()
	}
	return b.buff[b.r : b.r+n]
}

func (b *Buffer) Discard(n int) {
	if n > b.Size() {
		n = b.Size()
	}
	b.r += n
}

func (b *Buffer) Read(n int) []byte {
	data := b.Peek(n)
	b.Discard(n)
	return data
}

func (b *Buffer) ReadAll() []byte {
	return b.Read(b.Size())
}

func (b *Buffer) Write(data []byte) (int, error) {
	if b.r >= b.w {
		b.Reset()
	}
	wlen := len(data)
	b.allocate(wlen)
	copy(b.buff[b.w:], data)
	b.w += wlen
	return b.Size(), nil
}

func (b *Buffer) WriteByte(data byte) {
	if b.r >= b.w {
		b.Reset()
	}

	b.allocate(1)
	b.buff[b.w] = data
	b.w += 1
}
