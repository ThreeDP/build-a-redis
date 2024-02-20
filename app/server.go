package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"github.com/codecrafters-io/redis-starter-go/app/builtin"
)

type RedisServer struct {
	Env map[string]string
	Mutex sync.Mutex
}

func (s *RedisServer) handleCommand(buf string, conn net.Conn) {
	rpp := RedisProtocolParser{idx:0}
	it, _ := rpp.ParserProtocol(buf)
	cmd := it.([]string)
	cmd[0] = strings.ToLower(cmd[0])
	var b builtin.Builtin
	switch cmd[0] {
		case "echo":
			b = &builtin.Echo{Conn: conn}
		case "ping":
			b = &builtin.Ping{Conn: conn}
		case "get":
			b = &builtin.Get{Conn: conn, Env: s.Env, Mutex: &s.Mutex}
		case "set":
			b = &builtin.Set{Conn: conn, Env: s.Env, Mutex: &s.Mutex}
		default:
			return
	}
	b.Cmd(cmd[1:])
}

func (s *RedisServer) handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		if n == 0 {
			break
		}
		s.handleCommand(string(buf), conn)
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()
	s := RedisServer{Env: make(map[string]string)}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go s.handleClient(conn)
	}
}
