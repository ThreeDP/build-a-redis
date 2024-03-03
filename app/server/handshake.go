package server

import (
	"net"
	"fmt"
)

func (s *RedisServer) HandShake(conn net.Conn, handler func(net.Conn, string)) {
	s.Action = handler
	buf := make([]byte, 1024)
	params := [][]string{
		{"PING"},
		{"REPLCONF", "listening-port", s.Infos["server"]["port"]},
		{"REPLCONF", "capa", "npsync2"},
		{"PSYNC", "?", "-1"},
	}

	for _, param := range params {
		b, ok := s.Commands(param[0], conn, s.Time.Now())
		if !ok {return}
		b.Request(param)
		s.Read(conn, buf)
		s.Action(conn, string(buf))
	}
}

func (s *RedisServer) SlaveConnMaster() error {
	fmt.Printf("Connecting to %s\n", s.Infos["replication"]["role"])
	if s.Infos["replication"]["role"] == "slave" {
		conn, err := net.Dial("tcp",
		fmt.Sprintf("%s:%s",
			s.Infos["replication"]["master_host"],
			s.Infos["replication"]["master_port"],
		),
		)
		if err != nil {
			return err
		}
		s.HandShake(conn, s.HandleResponse)
		defer conn.Close()
	}
	return nil
}
