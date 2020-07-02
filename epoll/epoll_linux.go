// +build linux

package epoll

import (
	"sync"
	"syscall"
	"znet/connection"
)

const (
	size = 1 << 8
	msec = 1 << 7
)

type Epoll struct {
	epoll   int                           // Epoll File Descriptor
	sockets map[int]connection.Connection // Tcp Socket File Descriptor
	locker  *sync.RWMutex
}

// New Epoll
func NewEpoll() (*Epoll, error) {
	var epoll int
	var err error
	epoll, err = syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}
	return &Epoll{
		epoll:   epoll,
		sockets: make(map[int]connection.Connection, size),
		locker:  &sync.RWMutex{},
	}, nil
}

// Add Socket File Descriptor
func (e *Epoll) Add(c connection.Connection) error {
	ee := &syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		Fd:     int32(c.Socket()),
	}
	if err := syscall.EpollCtl(e.epoll, syscall.EPOLL_CTL_ADD, c.Socket(), ee); err != nil {
		return err
	}
	e.locker.Lock()
	defer e.locker.Unlock()
	e.sockets[c.Socket()] = c
	return nil
}

// Wait EpollEvent
func (e *Epoll) Wait() ([]connection.Connection, error) {
	events := make([]syscall.EpollEvent, size)
retry:
	n, err := syscall.EpollWait(e.epoll, events, msec)
	if err != nil {
		if err == syscall.EINTR {
			goto retry
		}
		return nil, err
	}
	if n == 0 {
		goto retry
	}
	e.locker.RLock()
	defer e.locker.RUnlock()
	var connections []connection.Connection
	for i := 0; i < n; i++ {
		// TODO: Handle Events
		if isErrorEvent(&events[i]) {
			socket := int(events[i].Fd)
			c := e.sockets[socket]
			_ = c.Close()
			continue
		}
		if isReadEvent(&events[i]) {
			socket := int(events[i].Fd)
			c := e.sockets[socket]
			connections = append(connections, c)
			continue
		}
	}
	return connections, nil
}

// Remove Socket File Descriptor
func (e *Epoll) Remove(c connection.Connection) error {
	if err := syscall.EpollCtl(e.epoll, syscall.EPOLL_CTL_DEL, c.Socket(), nil); err != nil {
		return err
	}
	e.locker.Lock()
	defer e.locker.Unlock()
	delete(e.sockets, c.Socket())
	return nil
}

// Close Epoll File Descriptor
func (e *Epoll) Close() error {
	return syscall.Close(e.epoll)
}

func isReadEvent(e *syscall.EpollEvent) bool {
	return 0 != e.Events&(syscall.EPOLLIN|syscall.EPOLLPRI|syscall.EPOLLRDHUP)
}

func isErrorEvent(e *syscall.EpollEvent) bool {
	return ((e.Events & syscall.EPOLLHUP) != 0) && ((e.Events & syscall.EPOLLIN) == 0)
}
