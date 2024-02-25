package builtin

import (
	"testing"
	"time"
)

func TestGetBuiltin(t *testing.T) {
	s := SetupFilesRDWR{}
	s.config(map[string]EnvData{
		"Percy": {Value: "Jackson", Expiry: s.TimeNow, MustExpire: false},
		"Key":   {Value: "Value", Expiry: s.TimeNow.Add(-10 * time.Millisecond), MustExpire: true},
	}, nil)

	t.Run("Test get key Percy and response Jackson", func(t *testing.T) {
		get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Percy"}
		copy(s.Expected, "$7\r\nJackson\r\n")

		get.Response(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test get key Any and response nil", func(t *testing.T) {
		get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Any"}
		copy(s.Expected, "$-1\r\n")

		get.Response(params)

		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})

	t.Run("Test get expired Key Any and response nil", func(t *testing.T) {
		get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
		params := []string{"Key"}
		copy(s.Expected, "$-1\r\n")

		get.Response(params)
		compareStrings(t, s.Expected, s.Out)
		s.reset()
	})
}

func BenchmarkGetBuiltin(b *testing.B) {
	s := SetupFilesRDWR{}
	s.config(map[string]EnvData{
		"Percy": {Value: "Jackson", Expiry: s.TimeNow, MustExpire: false},
		"Key":   {Value: "Value", Expiry: s.TimeNow.Add(-10 * time.Millisecond), MustExpire: true},
	}, nil)
	get := Get{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex}
	params := []string{"Percy"}

	for i := 0; i < b.N; i++ {
		get.Response(params)
	}
}
