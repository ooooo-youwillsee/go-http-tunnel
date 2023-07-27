package tunnel

import (
	"io"
	"net"
)

func CopyConn(localConn net.Conn, remoteConn net.Conn) chan error {
	errCh := make(chan error, 1)
	go copy(localConn, remoteConn, errCh)
	go copy(remoteConn, localConn, errCh)
	return errCh
}

func copy(src io.Reader, dst io.Writer, errCh chan error) {
	_, err := io.Copy(dst, src)
	if err != nil {
		errCh <- err
	}
}
