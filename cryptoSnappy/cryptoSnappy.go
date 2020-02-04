package cryptoSnappy

import (
	"git.iotserv.com/iotserv/utils/crypto/conn"
	"git.iotserv.com/iotserv/utils/snappy"
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
