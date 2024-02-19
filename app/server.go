package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var mutex sync.Mutex

type RedisServer struct {
	cn net.Conn
	rc RedisCommand
}

type RedisCommand struct {
	Map *map[string]string
}

func (rc *RedisCommand) EchoIn(cmd []string, cn net.Conn) {
	var str string
	size := len(cmd)
	for i := 0; i < size; i++ {
		if i == size -1 {
			str += cmd[i]
			continue 
		}
		str += cmd[i] + " "
	}
	cn.Write([]byte("+" + str + "\r\n"))
}

func (rc *RedisCommand) GetIn(cmd []string, cn net.Conn) {
	mutex.Lock()
	value, ok := (*rc.Map)[cmd[0]]
	mutex.Unlock()
	if !ok {
		cn.Write([]byte("$-1\r\n"))
		return
	}
	cn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
}

func (rc *RedisCommand) SetIn(cmd []string, cn net.Conn) {
	mutex.Lock()
	(*rc.Map)[cmd[0]] = cmd[1]
	mutex.Unlock()
	cn.Write([]byte("+OK\r\n"))
}

func (rc *RedisCommand) PingIn(cmd []string, cn net.Conn) {
	cn.Write([]byte("+PONG\r\n"))
}

func (s *RedisServer) handleCommand(buf string) {
	rpp := RedisProtocolParser{idx:0}
	it, _ := rpp.ParserProtocol(buf)
	cmd := it.([]string)
	cmd[0] = strings.ToLower(cmd[0])
	switch cmd[0] {
		case "echo":
			s.rc.EchoIn(cmd[1:], s.cn)
		case "ping":
			s.rc.PingIn(cmd[1:], s.cn)
		case "get":
			s.rc.GetIn(cmd[1:], s.cn)
		case "set":
			s.rc.SetIn(cmd[1:], s.cn)
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
	env := make(map[string]string) 
	for {
		s := RedisServer{rc: RedisCommand{Map: &env}}
		s.cn, err = l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go s.handleClient()
	}
}
