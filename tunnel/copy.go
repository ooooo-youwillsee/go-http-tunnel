package tunnel

import (
	"io"
	"net"
	"sync"
)

func copyDataOnConn(conn1 net.Conn, conn2 net.Conn) {
	var wg sync.WaitGroup
	// read data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyConn(conn1, conn2)
	}()

	// write data
	wg.Add(1)
	go func() {
		defer wg.Done()
		copyConn(conn2, conn1)
	}()
	wg.Wait()
}

func copyConn(conn1 net.Conn, conn2 net.Conn) {
	_, err := io.Copy(conn2, conn1)
	if err != nil {
		return
	}
}

//func copyConn(conn1 net.Conn, conn2 net.Conn) {
//	buf := make([]byte, 1024*1024)
//	var (
//		n   int
//		err error
//		nn  int
//	)
//	for {
//		n, err = conn1.Read(buf)
//		if err != nil {
//			return
//		}
//		nn = 0
//		for nn < n {
//			nnn, err := conn2.Write(buf[nn:n])
//			if err != nil {
//				return
//			}
//			nn += nnn
//		}
//	}
//}
