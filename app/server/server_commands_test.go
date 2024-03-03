package server

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/builtin"
)

func TestSetCommands(t *testing.T) {
	s := setupRedisServer(nil)

	t.Run("Test SetCommands with 'echo' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("echo", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command echo, but has not\n")
		}
		_, ok = b.(*builtin.Echo)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Echo{}, b)
		}
	})

	t.Run("Test SetCommands with 'ping' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("pinG", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command ping, but has not\n")
		}
		_, ok = b.(*builtin.Ping)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Ping{}, b)
		}
	})

	t.Run("Test SetCommands with 'info' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("iNfo", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command info, but has not\n")
		}
		_, ok = b.(*builtin.Info)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Info{}, b)
		}
	})

	t.Run("Test SetCommands with 'get' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("Get", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command get, but has not\n")
		}
		_, ok = b.(*builtin.Get)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Get{}, b)
		}
	})

	t.Run("Test SetCommands with 'set' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("SET", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command set, but has not\n")
		}
		_, ok = b.(*builtin.Set)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Set{}, b)
		}
	})

	t.Run("Test SetCommands with 'replconf' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("replconf", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command replconf, but has not\n")
		}
		_, ok = b.(*builtin.ReplConf)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.ReplConf{}, b)
		}
	})

	t.Run("Test SetCommands with 'psync' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("psynC", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command replconf, but has not\n")
		}
		_, ok = b.(*builtin.PSync)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.PSync{}, b)
		} 
	})

	t.Run("Test SetCommands with unknown command", func(t *testing.T) {
		s.SetCommands()
		_, ok := s.Commands("unknown", nil, s.Time.Now())
		if ok {
			t.Errorf("Expected unknown command, but has not\n")
		}
	})
}