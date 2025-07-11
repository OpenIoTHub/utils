package cryptoSnappy

import (
	"github.com/OpenIoTHub/utils/v2/crypto/conn"
	"github.com/OpenIoTHub/utils/v2/snappy"
	"io"
	"net"
)

func Convert(oldConn net.Conn, key []byte) (closer io.ReadWriteCloser, err error) {
	enConn, err := conn.WithEncryption(oldConn, key)
	if err != nil {
		return nil, err
	}
	comConn := snappy.NewCompStream(enConn)
	return comConn, nil
}
