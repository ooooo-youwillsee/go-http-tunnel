package tunnel

import (
	"io"
)

func Copy(localConn io.ReadWriteCloser, remoteConn io.ReadWriteCloser) chan error {
	errCh := make(chan error, 1)
	go copy(localConn, remoteConn, errCh)
	go copy(remoteConn, localConn, errCh)
	return errCh
}

func copy(src io.Reader, dst io.Writer, errCh chan error) {
	_, err := io.Copy(dst, src)
	errCh <- err
}
