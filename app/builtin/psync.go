package builtin

import (
	"fmt"
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func (p *PSync) Request(params []string) {
	p.Conn.Write([]byte(parser.ArrayEncode(params)))
}

func (p *PSync) Response(params []string) {
	str := fmt.Sprintf("+FULLRESYNC %s 0\r\n", p.Infos["replication"]["master_replid"])
	p.Conn.Write([]byte(str))
}

func (p *PSync) SetConn(conn net.Conn) {
	p.Conn = conn
}

func (p *PSync) SetTimeNow(now time.Time) {
	p.Now = now
}
//