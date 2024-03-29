package builtin

import (
	"testing"
)

func TestReplConfBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	getTime := TTime{}
	s.config(nil, nil)

	t.Run("Test handler Request reply conf listening-port 6380", func(t *testing.T) {
		rc := ReplConf{Conn: s.Conn, Env: s.Env, Now: getTime.Now()}
		params := []string{"listening-port", "6380"}
		copy(s.Expected, "+OK\r\n")

		rc.Response(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test handler Request reply conf capa npsync2", func(t *testing.T) {
		rc := ReplConf{Conn: s.Conn, Env: s.Env, Now: getTime.Now()}
		params := []string{"capa", "npsync2"}
		copy(s.Expected, "+OK\r\n")

		rc.Response(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test handler Reponse reply conf listening-port 6380", func(t *testing.T) {
		rc := ReplConf{Conn: s.Conn, Env: s.Env, Now: getTime.Now()}
		params := []string{"REPLCONF", "listening-port", "6380"}
		copy(s.Expected, "*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n6380\r\n")

		rc.Request(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})
}

func BenchmarkReplConfBuiltin(b *testing.B) {
	s := SetupFilesRDWR{}
	getTime := TTime{}
	s.config(nil, nil)
	rc := ReplConf{Conn: s.Conn, Env: s.Env, Now: getTime.Now()}
	
	b.Run("Benchmark handler Request reply conf listening-port 6380", func(b *testing.B) {
		params := []string{"listening-port", "6380"}
		for i := 0; i < b.N; i++ {
			rc.Response(params)
		}
	})

	b.Run("Benchmark handler Response reply conf listening-port 6380", func(b *testing.B) {
		params := []string{"REPLCONF", "listening-port", "6380"}
		for i := 0; i < b.N; i++ {
			rc.Request(params)
		}
	})
}