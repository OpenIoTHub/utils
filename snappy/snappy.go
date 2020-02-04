package snappy

import (
	"github.com/golang/snappy"
	"io"
)

type compStream struct {
	conn io.ReadWriteCloser
	w    *snappy.Writer
	r    *snappy.Reader
}

func (c *compStream) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

func (c *compStream) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	err = c.w.Flush()
	return n, err
}

func (c *compStream) Close() error {
	return c.conn.Close()
}

func NewCompStream(OldConn io.ReadWriteCloser) *compStream {
	c := new(compStream)
	c.conn = OldConn
	c.w = snappy.NewBufferedWriter(OldConn)
	c.r = snappy.NewReader(OldConn)
	return c
}
