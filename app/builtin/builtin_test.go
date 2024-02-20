package builtin

import (
	"testing"
	"reflect"
	"net"
	"time"
)

func TestSetBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(map[string]string{
		"Percy": "????",
	})

	t.Run("Test set Percy with value response Jackson", func(t *testing.T) {
		set := Set{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Percy", "Jackson"}
		copy(s.Expected, "+OK\r\n")

		set.Cmd(params)

		value, ok := s.Env[params[0]]
		if !ok {
			t.Error("variable don't set")
		}
		if value != params[1] {
			t.Errorf("Expected '%s', but has '%s'", params[1], value)
		}
		compareStrings(t, s.Expected, s.Out)
	})
}

func TestGetBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(map[string]string {
		"Percy": "Jackson",
		"Key": "Value",
	})

	t.Run("Test get key Percy and response Jackson", func(t *testing.T) {
		get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Percy"}
		copy(s.Expected, "$7\r\nJackson\r\n")

		get.Cmd(params)

		compareStrings(t, s.Expected, s.Out)
	})

	t.Run("Test get key Any and response nil", func(t *testing.T) {
		get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Any"}
		copy(s.Expected, "$-1\r\n")

		get.Cmd(params)

		compareStrings(t, s.Expected, s.Out)
	})
}

func TestpingBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(nil)

	t.Run("Test ping with string 'ping'", func(t *testing.T) {
		ping := Ping{Conn: s.Conn}
		params := []string{"ping"}
		copy(s.Expected, "+PONG\r\n")

		ping.Cmd(params)

		compareStrings(t, s.Expected, s.Out)
	})
}

func TestEchoBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(nil)

	t.Run("Test pass a \"hey\" string", func (t *testing.T) {
		echo := Echo{Conn: s.Conn}
		params := []string{"hey"}
		copy(s.Expected, "+hey\r\n")

		echo.Cmd(params)

		compareStrings(t, s.Expected, s.Out)
	})

	t.Run("Test pass a \"hey\" string", func (t *testing.T) {
		echo := Echo{Conn: s.Conn}
		params := []string{"hey", "ho"}
		copy(s.Expected, "+hey ho\r\n")

		echo.Cmd(params)

		compareStrings(t, s.Expected, s.Out)
	})
}

func compareStrings(t *testing.T, expected, received []byte) {
	t.Helper()
	if !reflect.DeepEqual(expected, received) {
		t.Errorf("expected value: '%s' len: %d cap: %d, but has value: '%s' len: %d cap: %d\n",
					expected, len(expected), cap(expected),
					received, len(received), cap(received))
	}
}

/*
	Test Conn struct Mock
*/
type TConn struct {
	In []byte
	Out []byte
}

func (c TConn) Close() error {return nil}
func (c TConn) SetDeadline(t time.Time) error {return nil}
func (c TConn) SetReadDeadline(t time.Time) error {return nil}
func (c TConn) SetWriteDeadline(t time.Time) error {return nil}

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


/*
	Mutex struct Mock
*/
type TMutex struct {}
func (m TMutex) Lock() {}
func (m TMutex) Unlock() {}
func (m TMutex) TryLock() bool {return false}

/*
	Setup out and in files to put result in vars
*/
type SetupFilesRDWR struct {
	In []byte
	Out []byte
	Expected []byte
	Conn TConn
	Env map[string]string
	Mutex IMutex
}

func (s *SetupFilesRDWR) config( values map[string]string) {
	s.In = make([]byte, 1024)
	s.Out = make([]byte, 1024)
	s.Expected = make([]byte, 1024)
	s.Conn = TConn{In: s.In, Out: s.Out}
	s.Env = values
	s.Mutex = TMutex{}
}