package builtin

import (
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/define"
)

type TTime struct{}

func (t TTime) Now() time.Time {
	return time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)
}

func TestTimeStructNow(t *testing.T) {
	time := TTime{}
	expected := time.Now()
	checkVarDate(t, expected, time.Now())
}

func checkVarValue(t *testing.T, ok bool, received, expected string) {
	t.Helper()
	if !ok {
		t.Error("variable don't set")
	}
	if received != expected {
		t.Errorf("Expected '%s', but has '%s'", expected, received)
	}
}

func checkVarDate(t *testing.T, received, expected time.Time) {
	t.Helper()
	if received != expected {
		t.Errorf("Expected %v, but has %v\n", expected, received)
	}
}

func checkMustExpire(t *testing.T, received, expected bool) {
	t.Helper()
	if received != expected {
		t.Errorf("Expected %t, but has %t\n", expected, expected)
	}
}

func compareStrings(t *testing.T, expected, received []byte) {
	t.Helper()
	if !reflect.DeepEqual(expected, received) {
		t.Errorf("expected value: '%s' len: %d cap: %d, but has value: '%s' len: %d cap: %d\n",
			expected, len(expected), cap(expected),
			received, len(received), cap(received))
	}
}

func compareSubStringsInString(t *testing.T, expected []string, received []byte) {
	t.Helper()

	for _, str := range expected {
		if !strings.Contains(string(received), str) {
			t.Errorf("expected value: '%s' not found in '%s'\n", str, received)
		}
	}
}

/*
Setup out and in files to put result in vars
*/
type SetupFilesRDWR struct {
	In       []byte
	Out      []byte
	Expected []byte
	Conn     TConn
	Env      map[string]EnvData
	Mutex    IMutex
	TimeNow  time.Time
}

func (s *SetupFilesRDWR) config(data map[string]EnvData) {
	s.In = make([]byte, define.BUFFERSIZE)
	s.Out = make([]byte, define.BUFFERSIZE)
	s.Expected = make([]byte, define.BUFFERSIZE)
	s.Conn = TConn{In: s.In, Out: s.Out}
	s.Env = data
	s.Mutex = TMutex{}
	s.TimeNow = time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)
}

func (s *SetupFilesRDWR) reset() {
	s.Expected = make([]byte, define.BUFFERSIZE)
	s.In = make([]byte, define.BUFFERSIZE)
	s.Out = make([]byte, define.BUFFERSIZE)
	s.Conn = TConn{In: s.In, Out: s.Out}
}

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

/*
Mutex struct Mock
*/
type TMutex struct{}
func (m TMutex) Lock()         {}
func (m TMutex) Unlock()       {}
func (m TMutex) TryLock() bool { return false }
