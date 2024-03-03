package server

import (
	"net"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func (s *RedisServer) HandleResponse(conn net.Conn, buf string) {
	rpp := parser.RedisProtocolParser{Idx: 0}
	it, _ := rpp.ParserProtocol(buf)
	cmd := it.([]string)
	fmt.Printf("%s\n", cmd[0])
}