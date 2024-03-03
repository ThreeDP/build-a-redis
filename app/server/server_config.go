package server

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/define"
)

// functions responsible for insert the server informations
func (s *RedisServer) insertInfosItem(section, key, value string) {
	if _, ok := s.Infos[section][key]; !ok {
		s.Infos[section][key] = value
	}
}

func (s *RedisServer) insertInfos(section, key, value string) {
	if _, ok := s.Infos[section]; ok {
		s.insertInfosItem(section, key, value)
	} else {
		s.Infos[section] = make(map[string]string)
		s.insertInfosItem(section, key, value)
	}
}

func (s *RedisServer) defaultInfos() {
	s.insertInfos("server", "port", define.DEFAULPORT)
	s.insertInfos("replication", "role", "master")
	s.insertInfos("replication", "master_replid", "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb")
	s.insertInfos("replication", "master_repl_offset", "0")
	listener := fmt.Sprintf(
		"name=%s,bind=%s,bind=%s,port=%s",
		"tcp",
		"*",
		"localhost",
		s.Infos["server"]["port"],
	)
	s.insertInfos("server", "listener0", listener)
}

func (s *RedisServer) HandleArgs() {
	size := len(s.Args)
	defer s.defaultInfos()
	for s.Idx = 0; s.Idx < size; s.Idx++ {
		switch s.Args[s.Idx] {
		case "--port":
			s.Idx++
			s.insertInfos("server", "port", s.Args[s.Idx])
		case "--replicaof":
			s.insertInfos("replication", "role", "slave")
			s.insertInfos("replication", "master_host", s.Args[s.Idx+1])
			s.insertInfos("replication", "master_port", s.Args[s.Idx+2])
			s.insertInfos("replication", "master_link_status", "down")
			s.Idx += 2
		}
	}
}
