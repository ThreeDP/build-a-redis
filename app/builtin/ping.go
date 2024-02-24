package builtin

import (
	"net"
	"time"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func (p *Ping) Request(params []string) {
	p.Conn.Write([]byte(parser.ArrayEncode([]string{"PING"})))
}

func (p *Ping) Received(params []string) {
	p.Conn.Write([]byte("+PONG\r\n"))
}

func (p *Ping) SetConn(conn net.Conn) {
	p.Conn = conn
}

func (p *Ping) SetTimeNow(now time.Time) {
	p.Now = now
}