package conn

import (
	"io"
	"sync"
)

func WithEncryption(rwc io.ReadWriteCloser, key []byte) (io.ReadWriteCloser, error) {
	w, err := NewWriter(rwc, key)
	if err != nil {
		return nil, err
	}
	return WrapReadWriteCloser(NewReader(rwc, key), w, func() error {
		return rwc.Close()
	}), nil
}

type ReadWriteCloser struct {
	r       io.Reader
	w       io.Writer
	closeFn func() error

	closed bool
	mu     sync.Mutex
}

// closeFn will be called only once
func WrapReadWriteCloser(r io.Reader, w io.Writer, closeFn func() error) io.ReadWriteCloser {
	return &ReadWriteCloser{
		r:       r,
		w:       w,
		closeFn: closeFn,
		closed:  false,
	}
}

func (rwc *ReadWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.r.Read(p)
}

func (rwc *ReadWriteCloser) Write(p []byte) (n int, err error) {
	return rwc.w.Write(p)
}

func (rwc *ReadWriteCloser) Close() (errRet error) {
	rwc.mu.Lock()
	if rwc.closed {
		rwc.mu.Unlock()
		return
	}
	rwc.closed = true
	rwc.mu.Unlock()

	var err error
	if rc, ok := rwc.r.(io.Closer); ok {
		err = rc.Close()
		if err != nil {
			errRet = err
		}
	}

	if wc, ok := rwc.w.(io.Closer); ok {
		err = wc.Close()
		if err != nil {
			errRet = err
		}
	}

	if rwc.closeFn != nil {
		err = rwc.closeFn()
		if err != nil {
			errRet = err
		}
	}
	return
}
