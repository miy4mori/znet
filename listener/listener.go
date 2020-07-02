// +build linux

package listener

import (
	"syscall"
	"znet/connection"
)

const backlog = 1 << 8

type Listener interface {
	// Listen
	Listen(port int) (err error)
	// Accept
	Accept() (c connection.Connection, err error)
	// Close
	Close() error
	// Socket file descriptor
	Socket() int
}

type listener struct {
	socket  int
	backlog int
}

func NewListener() Listener {
	return &listener{
		backlog: backlog,
	}
}

func (l *listener) Listen(port int) (err error) {
	// create socket
	l.socket, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return err
	}
	// create inet4 address
	inet4 := &syscall.SockaddrInet4{
		Port: port,
	}
	// bind address
	err = syscall.Bind(l.socket, inet4)
	if err != nil {
		return err
	}
	// listen socket
	err = syscall.Listen(l.socket, l.backlog)
	if err != nil {
		return err
	}

	return nil
}

func (l *listener) Accept() (c connection.Connection, err error) {
	nfd, _, err := syscall.Accept(l.socket)
	if err != nil {
		return nil, err
	}
	c = connection.NewConnection(nfd)

	return c, nil
}

func (l *listener) Socket() int {
	return l.socket
}

func (l *listener) Close() error {
	return syscall.Close(l.socket)
}
