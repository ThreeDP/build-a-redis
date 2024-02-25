package builtin

import (
	"fmt"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

type Get struct {
	Env   map[string]EnvData
	Mutex IMutex
	Conn  net.Conn
	Now   time.Time
}

func (g *Get) Response(params []string) {
	cparams := parser.FindNextRedisSerialization(params)
	g.Mutex.Lock()
	data, ok := g.Env[cparams[0]]
	g.Mutex.Unlock()
	if !ok || (data.Expiry.Before(g.Now) && data.MustExpire) {
		g.Conn.Write([]byte("$-1\r\n"))
		return
	}
	g.Conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(data.Value), data.Value)))
}

func (g *Get) Request(params []string) {
}

func (g *Get) SetConn(conn net.Conn) {
	g.Conn = conn
}

func (g *Get) GetConn() net.Conn {
	return g.Conn
}

func (g *Get) SetTimeNow(now time.Time) {
	g.Now = now
}

func (g *Get) GetTimeNow() time.Time {
	return g.Now
}