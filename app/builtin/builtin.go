package builtin

import (
	"net"
	"fmt"
	"time"
	"strconv"
	"strings"
)

type IMutex interface {
	Lock()
	TryLock() bool
	Unlock()
}

type Builtin interface {
	Cmd([]string)
}

type ITime interface {
	Now() time.Time
}

type Time struct {}
func (t Time) Now() time.Time {
	return time.Now()
}


const (
	BUFFERSIZE				= 1024
	RedisSimpleString 		= "+"
	RedisSimpleErrors 		= "-"
	RedisIntegers 			=":"
	RedisBulkStrings 		= "$"
	RedisArrays				= "*"
	RedisBooleans			= "#"
	RedisDoubles			= ","
	RedisBigNumbers			= "("
	RedisBulkErrors			= "!"
	RedisVerbatimStrings	= "="
	RedisMaps				= "%"
	RedisSets				= "~"
	RedisPushes				= ">"
)

var RedisSerialization = []string{
	RedisSimpleString,
	RedisSimpleErrors,
	RedisIntegers,
	RedisBulkStrings,
	RedisArrays,
	RedisBooleans,
	RedisDoubles,
	RedisBigNumbers,
	RedisBulkErrors,
	RedisVerbatimStrings,
	RedisMaps,
	RedisSets,
	RedisPushes,
}

func checkRedisSerialization(str string) bool {
	for _, rs := range RedisSerialization {
		if strings.HasPrefix(str, rs) {
			return true
		}
	}
	return false
}

func findNextRedisSerialization(params []string) []string {
	for i, param := range params {
		if checkRedisSerialization(param) {
			return params[:i]
		}
	}
	return params
}

type EnvData struct {
	Value string
	Expiry time.Time
	MustExpire bool
}

type Set struct {
	Conn net.Conn
	Env map[string]EnvData
	Mutex IMutex
	Now time.Time
}

func (s *Set) Cmd(params []string) {
	flag := false
	cParams := findNextRedisSerialization(params)
	if len(cParams) == 4 {
		if strings.ToUpper(cParams[2]) == "PX" {
			flag = true
		}
		expiryTime, err := strconv.Atoi(cParams[3])
		if err != nil {
			return
		}
		s.Now = s.Now.Add(time.Millisecond * time.Duration(expiryTime))
	}
	value := EnvData{Value: cParams[1], Expiry: s.Now, MustExpire: flag}
	s.Mutex.Lock()
	s.Env[cParams[0]] = value
	s.Mutex.Unlock()
	s.Conn.Write([]byte("+OK\r\n"))
}

type Get struct {
	Conn net.Conn
	Env map[string]EnvData
	Mutex IMutex
	Now time.Time
}

func (g *Get) Cmd(params []string) {
	cparams:= findNextRedisSerialization(params)
	g.Mutex.Lock()
	data, ok := g.Env[cparams[0]]
	g.Mutex.Unlock()
	if !ok || (data.Expiry.Before(g.Now) && data.MustExpire) {
		g.Conn.Write([]byte("$-1\r\n"))
		return
	}
	g.Conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(data.Value), data.Value)))
}

type Info struct {
	Conn net.Conn
}

func (e *Info) Cmd(params []string) () {
	if len(params) > 0 {
		switch params[0] {
			case "replication": 
				e.Conn.Write([]byte("*1\r\n$11\r\nrole:master\r\n"))
				return
		}
	}
	e.Conn.Write([]byte("*1\r\n$11\r\nrole:master\r\n"))
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
