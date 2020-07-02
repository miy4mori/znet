package connection

import (
	"bytes"
	"syscall"
)

const size = 1 << 16

type Connection interface {

	// Read reads data from the connection.
	Read() (buffer *bytes.Buffer, err error)

	// Write writes data to the connection.
	Write(buffer *bytes.Buffer) (n int, err error)

	// Close closes the connection.
	Close() error

	// Get socket file descriptor of the connection.
	Socket() int
}

type connection struct {
	socket int // Socket File Descriptor
}

func NewConnection(socket int) Connection {
	return &connection{socket: socket}
}

func (c *connection) Write(buffer *bytes.Buffer) (n int, err error) {
	n, err = syscall.Write(c.socket, buffer.Bytes())
	if err != nil {
		return 0, err
	}
	return n, err
}

func (c *connection) Read() (buffer *bytes.Buffer, err error) {
	buffers := make([]byte, size)
	n, err := syscall.Read(c.socket, buffers)
	if err != nil {
		return nil, err
	}
	buffer = bytes.NewBuffer(buffers[:n])
	return buffer, err
}

func (c *connection) Close() (err error) {
	return syscall.Close(c.socket)
}

func (c *connection) Socket() int {
	return c.socket
}
