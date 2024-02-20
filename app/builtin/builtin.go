package builtin

import (
	"net"
	"fmt"
)

type IMutex interface {
	Lock()
	TryLock() bool
	Unlock()
}

type Builtin interface {
	Cmd([]string)
}

type Set struct {
	Conn net.Conn
	Env map[string]string
	Mutex IMutex
}

func (s *Set) Cmd(params []string) {
	s.Mutex.Lock()
	s.Env[params[0]] = params[1]
	s.Mutex.Unlock()
	s.Conn.Write([]byte("+OK\r\n"))
}

type Get struct {
	Conn net.Conn
	Env map[string]string
	Mutex IMutex
}

func (g *Get) Cmd(params []string) {
	g.Mutex.Lock()
	value, ok := g.Env[params[0]]
	g.Mutex.Unlock()
	if !ok {
		g.Conn.Write([]byte("$-1\r\n"))
		return
	}
	g.Conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
}

type Echo struct {
	Conn net.Conn
}

func (e *Echo) Cmd(params []string) () {
	var str string
	size := len(params)
	for i := 0; i < size; i++ {
		if i == size -1 {
			str += params[i]
			continue 
		}
		str += params[i] + " "
	}
	e.Conn.Write([]byte("+" + str + "\r\n"))
}

type Ping struct {
	Conn net.Conn
}

func (p *Ping) Cmd(params []string) {
	p.Conn.Write([]byte("+PONG\r\n"))
}
