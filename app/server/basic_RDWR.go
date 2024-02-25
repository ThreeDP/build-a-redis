package server

import (
	"net"
)

func (s *RedisServer) Read(conn net.Conn, buf []byte) (int, error) {
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return n, err
	}
	return n, nil
}