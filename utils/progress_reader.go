package utils

import (
	"io"
	"sync/atomic"
)

type ProgressReader interface {
	io.Reader
	N() int64
	SetReader(reader io.Reader) io.Reader
}

type progressReaderImpl struct {
	r       io.Reader
	counter int64
}

func NewProgressReader(r io.Reader) ProgressReader {
	return &progressReaderImpl{
		r: r,
	}
}

var _ io.Reader = (*progressReaderImpl)(nil)

func (rr *progressReaderImpl) Read(p []byte) (n int, err error) {
	n, err = rr.r.Read(p)

	if err == nil {
		atomic.AddInt64(&rr.counter, int64(n))
	}
	return
}

func (rr *progressReaderImpl) N() int64 {
	return atomic.LoadInt64(&rr.counter)
}

func (rr *progressReaderImpl) SetReader(reader io.Reader) io.Reader {
	rr.r = reader
	return rr
}
