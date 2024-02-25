package builtin

import (
	"time"
	"net"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func (p *PSync) Request(params []string) {
	p.Conn.Write([]byte(parser.ArrayEncode(params)))
}

func (p *PSync) Response(params []string) {
	p.Conn.Write([]byte("+FULLRESYNC <REPL_ID> 0\r\n"))
}

func (p *PSync) SetConn(conn net.Conn) {
	p.Conn = conn
}

func (p *PSync) SetTimeNow(now time.Time) {
	p.Now = now
}
//