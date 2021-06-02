package epoll

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func Test_client(t *testing.T) {
	conn, _ := net.Dial("tcp", "127.0.0.1:8888")
	for i := 0; i < 1; i++ {
		time.Sleep(time.Second)
		conn.Write([]byte("Hi"))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		fmt.Printf("string(buf)=%s\n", string(buf[:n]))
	}
}
