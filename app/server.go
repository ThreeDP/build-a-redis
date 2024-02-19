package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type RedisServer struct {
	cn net.Conn
	rc RedisCommand
}

type RedisCommand struct {

}

func (rc *RedisCommand) BuiltinEcho(cmd []string, cn net.Conn) {
	for _, c := range cmd {
		cn.Write([]byte(c))
	}
}

func (rc *RedisCommand) BuiltinPing(cmd []string, cn net.Conn) {
	cn.Write([]byte("+PONG\r\n"))
}

func (s *RedisServer) handleCommand(buf string) {
	rpp := RedisProtocolParser{idx:0}
	it, _ := rpp.ParserProtocol(buf)
	cmd := it.([]string)
	cmd[0] = strings.ToLower(cmd[0])
	switch cmd[0] {
		case "echo":
			s.rc.BuiltinEcho(cmd[1:], s.cn)
		case "ping":
			s.rc.BuiltinPing(cmd[1:], s.cn)
	}
}

func (s *RedisServer) handleClient() {
	defer s.cn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := s.cn.Read(buf)
		if err != nil {
			return
		}
		if n == 0 {
			break
		}
		s.handleCommand(string(buf))
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		s := RedisServer{}
		s.cn, err = l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go s.handleClient()
	}
}
