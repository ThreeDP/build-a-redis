package server

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/builtin"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

const (
	DEFAULPORT = "6379"
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

type RedisServer struct {
	Env      map[string]builtin.EnvData
	Infos    map[string]map[string]string
	Mutex    sync.Mutex
	Time     builtin.ITime
	Args     []string
	Idx      int
	Listener net.Listener
}

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
	s.insertInfos("server", "port", DEFAULPORT)
	s.insertInfos("replication", "role", "master")
	s.insertInfos("replication", "master_replid", "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb")
	s.insertInfos("replication", "master_repl_offset", "0")
	listener := fmt.Sprintf(
		"name=%s,bind=%s,bind=%s,port=%s",
		"tcp",
		"*",
		"0.0.0.0",
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

func (s *RedisServer) SlaveConnMaster() error {
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
		defer conn.Close()
		ping := &builtin.Ping{Conn: conn}
		fmt.Print("Ping\n")
		ping.Request([]string{"PING"})
	} 
	return nil
}

func (s *RedisServer) Run() {
	defer s.Listener.Close()
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go s.handleClient(conn)
	}
}

func (s *RedisServer) getCommands(conn net.Conn) func(key string) (builtin.Builtin, bool) {
	commands := map[string]builtin.Builtin{
		"echo": &builtin.Echo{Conn: conn},
		"ping": &builtin.Ping{Conn: conn},
		"info": &builtin.Info{Conn: conn, Infos: s.Infos},
		"get": &builtin.Get{Conn: conn, Env: s.Env,
			Mutex: &s.Mutex, Now: s.Time.Now()},
		"set": &builtin.Set{Conn: conn, Env: s.Env,
			Mutex: &s.Mutex, Now: s.Time.Now()},
	}
	return func(key string) (builtin.Builtin, bool) {
		elem, ok := commands[key]
		return elem, ok
	}
}

func (s *RedisServer) handleCommand(buf string, conn net.Conn) {
	var b builtin.Builtin
	commands := s.getCommands(conn)
	rpp := parser.RedisProtocolParser{Idx: 0}
	it, _ := rpp.ParserProtocol(buf)
	cmd := it.([]string)
	b, ok := commands(strings.ToLower(cmd[0]))
	if !ok {
		err := fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd[0])
		conn.Write([]byte(err))
		return
	}
	b.Received(cmd[1:])
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
