package utils

import (
	"io"
	"sync/atomic"
)

type progressWriter struct {
	w       io.Writer
	counter int64
}

func NewProgressWriter(w io.Writer) *progressWriter {
	return &progressWriter{
		w: w,
	}
}

var _ io.Writer = (*progressWriter)(nil)

func (ww *progressWriter) Write(p []byte) (n int, err error) {
	n, err = ww.w.Write(p)

	if err == nil {
		atomic.AddInt64(&ww.counter, int64(n))
	}
	return
}

func (ww *progressWriter) SetWriter(w io.Writer) io.Writer {
	ww.w = w
	return ww
}

func (ww *progressWriter) N() int64 {
	return atomic.LoadInt64(&ww.counter)
}
