package buffer

import (
	"bytes"
	"testing"
)

func TestBuffer(t *testing.T) {
	b := NewBufferSize(32, 8)
	wd := []byte("0123456789ABCDEF")
	b.Write(wd)
	if b.Size() != len(wd) || b.w != len(wd) || b.r != 0 {
		t.Errorf("b error:%v", b)
	}

	b.Write(wd)
	if b.Size() != 32 || b.w != 32 || b.r != 0 {
		t.Errorf("b error:%v", b)
	}

	t.Log(b)

	b.Write(wd)

	t.Log(b)

	rdata := b.Peek(16)
	if bytes.Compare(rdata, wd) != 0 {
		t.Errorf("b.Peek compare failed:%s", string(rdata))
	}
	t.Log(b)

	b.Discard(16)
	if b.r != 16 {
		t.Errorf("b.Discard error:%v", b)
	}
	t.Log(b)

	rdata = b.Read(16)
	if bytes.Compare(rdata, wd) != 0 {
		t.Errorf("b.Read compare failed:%s", string(rdata))
	}
	t.Log(b)

	if b.Size() != 16 {
		t.Errorf("b.Size error: buff[%s] len[%d], r:%d, w:%d", string(b.buff), b.Size(), b.r, b.w)
	}
	t.Log(b)
}
