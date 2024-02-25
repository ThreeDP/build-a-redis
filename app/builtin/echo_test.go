package builtin

import (
	"testing"
)

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

	t.Run("Test pass a [\"hey\", \"ho\"]  string", func(t *testing.T) {
		echo := Echo{Conn: s.Conn}
		params := []string{"hey", "ho"}
		copy(s.Expected, "+hey ho\r\n")

		echo.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test pass a [] string", func(t *testing.T) {

		echo := Echo{Conn: s.Conn}
		params := []string{}
		copy(s.Expected, "+\r\n")

		echo.Received(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})
}

func BenchmarkEchoBuiltin(b *testing.B) {
	s := SetupFilesRDWR{}
	s.config(nil)
	params := []string{"hey", "ho", "lets", "go"}
	echo := Echo{Conn: s.Conn}

	for i := 0; i < b.N; i++ {
		echo.Received(params)
	}
}