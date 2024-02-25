package main

import (
	"fmt"
	"os"
	"sync"
	"github.com/codecrafters-io/redis-starter-go/app/builtin"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

func main() {
	s := server.RedisServer{
		Env: make(map[string]builtin.EnvData),
		Time: builtin.Time{},
		Args: os.Args,
		Infos: make(map[string]map[string]string),
		Mutex: &sync.Mutex{},
	}
	
	err := s.Listen(server.NetListen{})
	if err != nil {
		fmt.Println("Failed to bind to port " + s.Infos["server"]["port"])
		os.Exit(1)
	}
	
	s.SetCommands()
	err = s.SlaveConnMaster()
	if err != nil {
		fmt.Println("Failed to bind to port " + s.Infos["replication"]["master_port"])
		os.Exit(1)
	}
	s.Run()
}
