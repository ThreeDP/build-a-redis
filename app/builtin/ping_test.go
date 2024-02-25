package builtin

import (
	"testing"
)

func TestPingBuiltinRequest(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(nil, nil)

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
	s.config(nil, nil)
	params := []string{"ping"}
	ping := Ping{Conn: s.Conn}

	for i := 0; i < b.N; i++ {
		ping.Response(params)
	}
}
