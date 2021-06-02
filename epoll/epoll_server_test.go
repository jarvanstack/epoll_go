package epoll

import (
	"net"
	"syscall"
	"testing"
)
const (
	maxListenNumber = 10
	maxEpollNumber  = 1000
	port            = 8888
	ip              = "127.0.0.1"
)
//socket->bind->listen->accept->send/recv->closesocket
func Test_epoll_server(t *testing.T) {
	var server_addr  *syscall.SockaddrInet4
	//1.socket
	//ForkLock Doc:2) Socket. Does not block. Use the ForkLock.
	syscall.ForkLock.Lock()
	listen_fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	syscall.ForkLock.Unlock()
	// Allow reuse of recently-used addresses.
	if err = syscall.SetsockoptInt(listen_fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		syscall.Close(listen_fd)
	}
	//2.bind
	server_addr = &syscall.SockaddrInet4{Port: port}
	copy(server_addr.Addr[:], net.ParseIP(ip))
	if err = syscall.Bind(listen_fd, server_addr); err != nil {
	}
	//3.Listen for incoming connections.
	if err = syscall.Listen(listen_fd, maxListenNumber); err != nil {
	}
	//---------------epoll-------------

	//1.create
	epfd, _ := syscall.EpollCreate(maxEpollNumber)
	//2.control
	event := &syscall.EpollEvent{Fd: int32(listen_fd), Events: syscall.EPOLLIN }
	_ = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, listen_fd, event)
	//3.wait
	myEvents := make([]syscall.EpollEvent, maxEpollNumber)
	//var sockaddr syscall.Sockaddr
	var client_fd int
	for true {
		activeEN, _ := syscall.EpollWait(epfd, myEvents, maxEpollNumber)
		for i := 0; i < activeEN; i++ {
			if myEvents[i].Fd == int32(listen_fd) {
				//client_fd create.
				client_fd, _, err = syscall.Accept(listen_fd)
				eventOP := syscall.EPOLLIN
				event.Events = uint32(eventOP)
				event.Fd = int32(client_fd)
				syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, client_fd, event)
			}else if (myEvents[i].Events & syscall.EPOLLIN) >0 {
				//epoll_in is read.
				client_fd = int(myEvents[i].Fd)
				buf := make([]byte,1024)
				readN, _ := syscall.Read(client_fd, buf)
				if readN<0 {
				}else if readN == 0 {
					syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, client_fd, event)
				}else {
					//write back now.
					//or add epoll event to write back.(like before.)
					syscall.Write(client_fd, []byte("Hi,I am Server1"))
				}
			}else if (myEvents[i].Events & syscall.EPOLLOUT) >0 {
				//epoll_out is write.
				client_fd = int(myEvents[i].Fd)
				syscall.Write(client_fd, []byte("Hi,I am Server2"))
			}
		}
	}
	syscall.Close(epfd)
	syscall.Close(listen_fd)
}
