package server

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/builtin"
)

// Interface responsible for decoupling the
// net.Listen function.
type INetIListen interface {
	Listen(network, address string) (net.Listener, error)
}

// struct responsible for implementing the
// INetIListen interface and call the net.Listen function.
type NetListen struct{}

func (l NetListen) Listen(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func NewRedisServer() RedisServer {
	return RedisServer {
		Env: make(map[string]builtin.EnvData),
		Infos: make(map[string]map[string]string),
		Mutex: &sync.Mutex{},
		Time: builtin.Time{},
		Args: os.Args,
		Idx: 0,
		Listener: nil,
		Action: nil,
		Commands: nil,
	}
}

type RedisServer struct {
	Env      map[string]builtin.EnvData
	Infos    map[string]map[string]string
	Mutex    builtin.IMutex
	Time     builtin.ITime
	Args     []string
	Idx      int
	Listener net.Listener
	Action   func(net.Conn, string)
	Commands func(key string, conn net.Conn, now time.Time) (builtin.Builtin, bool)
}

func (s *RedisServer) Listen(nl INetIListen) error {
	s.HandleArgs()
	serverArgs := strings.Split(s.Infos["server"]["listener0"], ",")
	network := strings.Split(serverArgs[0], "=")[1]
	address := strings.Split(serverArgs[2], "=")[1]
	port := strings.Split(serverArgs[3], "=")[1]

	l, err := nl.Listen(network, address+":"+port)
	if err != nil {
		return err
	}
	s.Listener = l
	fmt.Printf("Server listening on %s://%s:%s\n", network, address, port)
	return nil
}

func (s *RedisServer) Handler(conn net.Conn, handler func(net.Conn, string)) {
	defer conn.Close()
	buf := make([]byte, 1024)
	s.Action = handler

	for {
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			return
		}
		s.Action(conn, string(buf))
	}
}

func (s *RedisServer) Run() {
	defer s.Listener.Close()
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go s.Handler(conn, s.HandleRequest)
	}
}
