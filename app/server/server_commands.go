package server

import (
	"net"
	"time"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/builtin"
)

func (s *RedisServer) SetCommands() {
	commands := map[string]builtin.Builtin{
		"echo":     &builtin.Echo{Conn: nil, Now: time.Time{}},
		"ping":     &builtin.Ping{Conn: nil, Now: time.Time{}},
		"info":     &builtin.Info{Infos: s.Infos, Conn: nil, Now: time.Time{}},
		"get":      &builtin.Get{Mutex: s.Mutex, Env: s.Env, Conn: nil, Now: time.Time{}},
		"set":      &builtin.Set{Mutex: s.Mutex, Env: s.Env, Conn: nil, Now: time.Time{}},
		"replconf": &builtin.ReplConf{Conn: nil, Env: s.Env, Now: time.Time{}},
		"psync":	&builtin.PSync{Conn: nil, Now: time.Time{}, Infos: s.Infos},
	}
	s.Commands = func(key string, conn net.Conn, now time.Time) (builtin.Builtin, bool) {
		elem, ok := commands[strings.ToLower(key)]
		if ok {
			elem.SetConn(conn)
			elem.SetTimeNow(now)
		}
		return elem, ok
	}
}