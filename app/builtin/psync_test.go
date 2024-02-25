package builtin

import (
	"testing"
)

func TestPSyncBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	getTime := TTime{}
	s.config(nil, map[string]map[string]string{
		"replication": {
			"master_replid": "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
			"master_repl_offset": "0",
		}})

	t.Run("Test Request ['?', '-1'] command", func(t *testing.T) {
		sp := PSync{Conn: s.Conn, Infos: s.Infos, Now: getTime.Now()}
		params := []string{"?", "-1"}
		copy(s.Expected, "+FULLRESYNC 8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb 0\r\n")

		sp.Response(params)
		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test Response PSYNC ? -1 test", func(t *testing.T) {
		sp := PSync{Conn: s.Conn, Infos: s.Infos, Now: getTime.Now()}
		params := []string{"PSYNC", "?", "-1"}
		copy(s.Expected, "*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n")

		sp.Request(params)
		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})
}

func BenchmarkPSyncBuiltin(b *testing.B) {
	s := SetupFilesRDWR{}
	getTime := TTime{}
	s.config(nil, map[string]map[string]string{
		"replication": {
			"master_replid": "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
			"master_repl_offset": "0",
		},
	})
	rc := PSync{Conn: s.Conn, Infos: s.Infos, Now: getTime.Now()}
	
	b.Run("Benchmark handler Request psync ? -1", func(b *testing.B) {
		params := []string{"?", "-1"}
		for i := 0; i < b.N; i++ {
			rc.Response(params)
		}
	})

	b.Run("Benchmark handler Response psync ? -1", func(b *testing.B) {
		params := []string{"REPLCONF", "?", "-1"}
		for i := 0; i < b.N; i++ {
			rc.Request(params)
		}
	})
}
