package server

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/builtin"
	"github.com/codecrafters-io/redis-starter-go/app/define"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
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
	Mutex    builtin.IMutex
	Time     builtin.ITime
	Args     []string
	Idx      int
	Listener net.Listener
	Action   func(net.Conn, string)
	Commands func(key string, conn net.Conn, now time.Time) (builtin.Builtin, bool)
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

func (s *RedisServer) SetCommands() {
	commands := map[string]builtin.Builtin{
		"echo":     &builtin.Echo{Conn: nil, Now: time.Time{}},
		"ping":     &builtin.Ping{Conn: nil, Now: time.Time{}},
		"info":     &builtin.Info{Infos: s.Infos, Conn: nil, Now: time.Time{}},
		"get":      &builtin.Get{Mutex: s.Mutex, Env: s.Env, Conn: nil, Now: time.Time{}},
		"set":      &builtin.Set{Mutex: s.Mutex, Env: s.Env, Conn: nil, Now: time.Time{}},
		"replconf": &builtin.ReplConf{Conn: nil, Env: s.Env, Now: time.Time{}},
	}
	s.Commands = func(key string, conn net.Conn, now time.Time) (builtin.Builtin, bool) {
		elem, ok := commands[strings.ToLower(key)]
		if ok {
			elem.SetConn(conn)
			elem.SetTimeNow(now)
		}
		return elem, ok
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
		defer conn.Close()
		ping := &builtin.Ping{Conn: conn}
		ping.Request([]string{"PING"})
		s.Handler2(conn, s.HandleResponse)
		conn, err = net.Dial("tcp",
			fmt.Sprintf("%s:%s",
				s.Infos["replication"]["master_host"],
				s.Infos["replication"]["master_port"],
			),
		)
		if err != nil {
			return err
		}
		rc := &builtin.ReplConf{Conn: conn}
		rc.Request([]string{"REPLCONF", "listening-port", s.Infos["server"]["port"]})
		s.Handler2(conn, s.HandleResponse)
		conn, err = net.Dial("tcp",
		fmt.Sprintf("%s:%s",
		s.Infos["replication"]["master_host"],
				s.Infos["replication"]["master_port"],
			),
		)
		if err != nil {
			return err
		}
		rc = &builtin.ReplConf{Conn: conn}
		rc.Request([]string{"REPLCONF", "capa", "npsync2"})
		s.Handler2(conn, s.HandleResponse)
	}
	return nil
}

func (s *RedisServer) HandleRequest(conn net.Conn, buf string) {
	var b builtin.Builtin
	rpp := parser.RedisProtocolParser{Idx: 0}
	it, _ := rpp.ParserProtocol(buf)
	cmd := it.([]string)
	b, ok := s.Commands(cmd[0], conn, s.Time.Now())
	if !ok { // Error Response
		err := fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd[0])
		conn.Write([]byte(err))
		return
	}
	b.Response(cmd[1:])
}

func (s *RedisServer) HandleResponse(conn net.Conn, buf string) {
	fmt.Printf("%s\n", buf)
	rpp := parser.RedisProtocolParser{Idx: 0}
	it, _ := rpp.ParserProtocol(buf)
	cmd := it.([]string)
	fmt.Printf("%s\n", cmd[0])
}

func (s *RedisServer) Handler(conn net.Conn, handler func(net.Conn, string)) {
	defer conn.Close()
	buf := make([]byte, 1024)
	s.Action = handler

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		if n == 0 {
			return
		}
		s.Action(conn, string(buf))
	}
}

func (s *RedisServer) Handler2(conn net.Conn, handler func(net.Conn, string)) {
	defer conn.Close()
	buf := make([]byte, 1024)
	s.Action = handler

	// for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		if n == 0 {
			return
		}
		s.Action(conn, string(buf))
	// }
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
