package builtin

import (
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/define"
)

/*
Test Conn struct Mock
*/
type TConn struct {
	In  []byte
	Out []byte
}

func (c TConn) Close() error                       { return nil }
func (c TConn) SetDeadline(t time.Time) error      { return nil }
func (c TConn) SetReadDeadline(t time.Time) error  { return nil }
func (c TConn) SetWriteDeadline(t time.Time) error { return nil }

func (c TConn) Read(b []byte) (n int, err error) {
	copy(c.In, b)
	return len(b), nil
}

func (c TConn) Write(b []byte) (n int, err error) {
	copy(c.Out, b)
	return len(b), nil
}

func (c TConn) LocalAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}
}

func (c TConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}
}

func (c *TConn) NewConn() {
	c.In = make([]byte, define.BUFFERSIZE)
	c.Out = make([]byte, define.BUFFERSIZE)
}

/*
Mutex struct Mock
*/
type TMutex struct{}
func (m TMutex) Lock()         {}
func (m TMutex) Unlock()       {}
func (m TMutex) TryLock() bool { return false }