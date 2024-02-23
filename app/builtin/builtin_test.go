package builtin

import (
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

type TTime struct{}

func (t TTime) Now() time.Time {
	return time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)
}

func TestInfoBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(nil)

	t.Run("Test Info only command", func(t *testing.T) {
		i := map[string]map[string]string{
			"replication": {
				"role": "master",
			},
		}
		info := Info{Conn: s.Conn, Infos: i}
		params := []string{}
		copy(s.Expected, "$27\r\n## Replication\n\nrole:master\r\n")

		info.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test Info command with replication arg", func(t *testing.T) {
		i := map[string]map[string]string{
			"replication": {
				"role": "slave",
			},
		}
		info := Info{Conn: s.Conn, Infos: i}
		params := []string{"RepLication"}
		copy(s.Expected, "$26\r\n## Replication\n\nrole:slave\r\n")

		info.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test Info command with replication arg more keys", func(t *testing.T) {
		i := map[string]map[string]string{
			"replication": {
				"role":               "slave",
				"master_replid":      "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
				"master_repl_offset": "0",
			},
		}
		info := Info{Conn: s.Conn, Infos: i}
		params := []string{"RepLication"}
		expected := []string{
			"$102\r\n",
			"## Replication\n\n",
			"master_repl_offset:0",
			"role:slave",
			"master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		}

		info.Received(params)

		compareSubStringsInString(t, expected, s.Out)
		s.reset()
	})
}

func TestSetBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	getTime := TTime{}
	s.config(map[string]EnvData{
		"Percy": {Value: "???", Expiry: s.TimeNow, MustExpire: false},
		"EX":    {Value: "?!?", Expiry: s.TimeNow.Add(time.Millisecond * 100), MustExpire: true},
	})

	t.Run("Test set Percy with value response Jackson", func(t *testing.T) {
		set := Set{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex, Now: getTime.Now()}
		params := []string{"Percy", "Jackson"}
		copy(s.Expected, "+OK\r\n")

		set.Received(params)

		data, ok := s.Env[params[0]]
		checkVarValue(t, ok, data.Value, params[1])
		checkVarDate(t, data.Expiry, s.TimeNow)
		checkMustExpire(t, data.MustExpire, false)
		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test set Percy PX 100 with value response Jackson", func(t *testing.T) {
		set := Set{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex, Now: getTime.Now()}
		params := []string{"Minute", "10Sec", "Px", "100"}
		copy(s.Expected, "+OK\r\n")

		set.Received(params)

		data, ok := s.Env[params[0]]
		checkVarValue(t, ok, data.Value, params[1])
		checkVarDate(t, data.Expiry, s.TimeNow.Add(time.Millisecond*100))
		checkMustExpire(t, data.MustExpire, true)
		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test set {minute, 10sec, PX, 100, $7}", func(t *testing.T) {
		set := Set{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex, Now: getTime.Now()}
		params := []string{"Minute", "10Sec", "PX", "100", "$7"}
		copy(s.Expected, "+OK\r\n")

		set.Received(params)

		data, ok := s.Env[params[0]]
		checkVarValue(t, ok, data.Value, params[1])
		checkVarDate(t, data.Expiry, s.TimeNow.Add(time.Millisecond*100))
		checkMustExpire(t, data.MustExpire, true)
		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})
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

func TestGetBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(map[string]EnvData{
		"Percy": {Value: "Jackson", Expiry: s.TimeNow, MustExpire: false},
		"Key":   {Value: "Value", Expiry: s.TimeNow.Add(-10 * time.Millisecond), MustExpire: true},
	})

	t.Run("Test get key Percy and response Jackson", func(t *testing.T) {
		get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Percy"}
		copy(s.Expected, "$7\r\nJackson\r\n")

		get.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test get key Any and response nil", func(t *testing.T) {
		get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Any"}
		copy(s.Expected, "$-1\r\n")

		get.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test get expired Key Any and response nil", func(t *testing.T) {
		get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Key"}
		copy(s.Expected, "$-1\r\n")

		get.Received(params)
		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})
}

func TestPingBuiltinReceived(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(nil)

	t.Run("Test ping with string 'ping'", func(t *testing.T) {
		ping := Ping{Conn: s.Conn}
		params := []string{"ping"}
		copy(s.Expected, "+PONG\r\n")

		ping.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})
}

func TestPingBuiltinRequest(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(nil)

	t.Run("Test ping with string 'ping'", func(t *testing.T) {
		ping := Ping{Conn: s.Conn}
		params := []string{"ping"}
		copy(s.Expected, "*1\r\n$4\r\nPING\r\n")

		ping.Request(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

}

func BenchmarkPingBuiltin(b *testing.B) {
	s := SetupFilesRDWR{}
	s.config(nil)
	params := []string{"ping"}
	copy(s.Expected, "+PONG\r\n")

	for i := 0; i < b.N; i++ {
		ping := Ping{Conn: s.Conn}

		ping.Received(params)
		s.reset()
	}
}

func TestEchoBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(nil)

	t.Run("Test pass a \"hey\" string", func(t *testing.T) {
		echo := Echo{Conn: s.Conn}
		params := []string{"hey"}
		copy(s.Expected, "+hey\r\n")

		echo.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test pass a \"hey\" string", func(t *testing.T) {
		echo := Echo{Conn: s.Conn}
		params := []string{"hey", "ho"}
		copy(s.Expected, "+hey ho\r\n")

		echo.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
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

func compareSubStringsInString(t *testing.T, expected []string, received []byte) {
	t.Helper()

	for _, str := range expected {
		if !strings.Contains(string(received), str) {
			t.Errorf("expected value: '%s' not found in '%s'\n", str, received)
		}
	}
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
	s.In = make([]byte, BUFFERSIZE)
	s.Out = make([]byte, BUFFERSIZE)
	s.Expected = make([]byte, BUFFERSIZE)
	s.Conn = TConn{In: s.In, Out: s.Out}
	s.Env = data
	s.Mutex = TMutex{}
	s.TimeNow = time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)
}

func (s *SetupFilesRDWR) reset() {
	s.Expected = make([]byte, BUFFERSIZE)
	s.In = make([]byte, BUFFERSIZE)
	s.Out = make([]byte, BUFFERSIZE)
	s.Conn = TConn{In: s.In, Out: s.Out}
}
