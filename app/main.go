package main

import (
	"fmt"
	"os"
	"github.com/codecrafters-io/redis-starter-go/app/builtin"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)



func main() {
	fmt.Println("Logs from your program will appear here!")
	s := server.RedisServer{
		Env: make(map[string]builtin.EnvData),
		Time: builtin.Time{},
		Args: os.Args,
		Infos: make(map[string]map[string]string),
	}
	err := s.Listen(server.NetListen{})
	if err != nil {
		fmt.Println("Failed to bind to port " + s.Infos["server"]["port"])
		os.Exit(1)
	}
	s.Run()
}
