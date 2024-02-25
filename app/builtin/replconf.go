package builtin

import (
	"net"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

type ReplConf struct {
	Conn net.Conn
	Env  map[string]EnvData
	Now  time.Time
}

func (rc *ReplConf) Request(params []string) {
	rc.Conn.Write([]byte(parser.ArrayEncode(params)))
}

func (rc *ReplConf) Response(params []string) {
	rc.Conn.Write([]byte("+OK\r\n"))
}

func (rc *ReplConf) SetConn(conn net.Conn) {
	rc.Conn = conn
}

func (rc * ReplConf) GetConn() net.Conn {
	return rc.Conn
}

func (rc *ReplConf) SetTimeNow(now time.Time) {
	rc.Now = now
}

func (rc *ReplConf) GetTimeNow() time.Time {
	return rc.Now
}