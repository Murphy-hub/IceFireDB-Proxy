package utils

import (
	"io"
)

var bufStepSize = 1024

type Reader struct {
	reader        io.Reader
	Buffer        []byte
	ReadPosition  int
	WritePosition int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{reader: reader, Buffer: make([]byte, bufStepSize)}
}

func (r *Reader) requestSpace(reqSize int) {
	ccap := cap(r.Buffer)
	if r.WritePosition+reqSize > ccap {
		newbuff := make([]byte, max(ccap*2, ccap+reqSize+bufStepSize))
		copy(newbuff, r.Buffer)
		r.Buffer = newbuff
	}
}

func (r *Reader) ReadSome(min int) error {
	r.requestSpace(min)
	nr, err := io.ReadAtLeast(r.reader, r.Buffer[r.WritePosition:], min)
	if err != nil {
		return err
	}
	r.WritePosition += nr
	return nil
}

func (r *Reader) RequireNBytes(num int) error {
	a := r.WritePosition - r.ReadPosition
	if a >= num {
		return nil
	}
	if err := r.ReadSome(num - a); err != nil {
		return err
	}
	return nil
}

func (r *Reader) GetNbytes(num int) (data []byte, err error) {
	err = r.RequireNBytes(num)

	if err != nil {
		return
	}

	data = r.Buffer[r.ReadPosition : r.ReadPosition+num]

	r.ReadPosition += num
	return
}

func (r *Reader) IsEnd() (ret bool) {
	ret = false

	if r.ReadPosition >= r.WritePosition {
		ret = true
		r.Reset()
	}

	return
}

func (r *Reader) Reset() {
	r.WritePosition = 0
	r.ReadPosition = 0
}
