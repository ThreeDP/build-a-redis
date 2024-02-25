package builtin

import (
	"testing"
)

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

func BenchmarkInfosBuiltin(b *testing.B) {
	s := SetupFilesRDWR{}
	params := []string{"RepLication"}
	inf := map[string]map[string]string{
		"replication": {
			"role":               "slave",
			"master_replid":      "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
			"master_repl_offset": "0",
		},
	}
	info := Info{Conn: s.Conn, Infos: inf}
	
	for i := 0; i < b.N; i++ {
		info.Received(params)
	}
}