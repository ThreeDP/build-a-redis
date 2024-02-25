package builtin

import (
	"strconv"
	"testing"
	"time"
)

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

func BenchmarkSetBuiltin(b *testing.B) {
	s := SetupFilesRDWR{}
	getTime := TTime{}
	s.config(map[string]EnvData{
		"Percy": {Value: "???", Expiry: s.TimeNow, MustExpire: false},
		"EX":    {Value: "?!?", Expiry: s.TimeNow.Add(time.Millisecond * 100), MustExpire: true},
	})
	set := Set{Conn: s.Conn, Env: s.Env, Mutex: s.Mutex, Now: getTime.Now()}

	for i := 0; i < b.N; i++ {
		set.Received([]string{strconv.Itoa(i), "test"})
	}
}