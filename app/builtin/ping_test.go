package builtin

import (
	"reflect"
	"testing"
)

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

func TestPingSets(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(nil)

	t.Run("Test set Ping time now", func(t *testing.T) {
		ping := Ping{Conn: s.Conn}
		ping.SetTimeNow(s.TimeNow)

		if ping.Now != s.TimeNow {
			t.Errorf("Expected %v, but has %v\n", s.TimeNow, ping.Now)
		}
	})

	t.Run("Test set Ping conn", func(t *testing.T) {
		ping := Ping{}
		ping.SetConn(s.Conn)

		if !reflect.DeepEqual(ping.Conn, s.Conn) {
			t.Errorf("Expected %v, but has %v\n", s.Conn, ping.Conn)
		}
	})
}

func BenchmarkPingBuiltin(b *testing.B) {
	s := SetupFilesRDWR{}
	s.config(nil)
	params := []string{"ping"}
	ping := Ping{Conn: s.Conn}

	for i := 0; i < b.N; i++ {
		ping.Received(params)
	}
}
