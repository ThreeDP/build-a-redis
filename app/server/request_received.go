package server

import (
	"net"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/builtin"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func (s *RedisServer) HandleRequest(conn net.Conn, buf string) {
	var b builtin.Builtin
	rpp := parser.RedisProtocolParser{Idx: 0}
	it, _ := rpp.ParserProtocol(buf)
	cmd := it.([]string)
	b, ok := s.Commands(cmd[0], conn, s.Time.Now())
	if !ok { // Error Response
		err := fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd[0])
		conn.Write([]byte(err))
		return
	}
	b.Response(cmd[1:])
}