package builtin

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
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

type Time struct{}

func (t Time) Now() time.Time {
	return time.Now()
}

const (
	BUFFERSIZE           = 1024
	RedisSimpleString    = "+"
	RedisSimpleErrors    = "-"
	RedisIntegers        = ":"
	RedisBulkStrings     = "$"
	RedisArrays          = "*"
	RedisBooleans        = "#"
	RedisDoubles         = ","
	RedisBigNumbers      = "("
	RedisBulkErrors      = "!"
	RedisVerbatimStrings = "="
	RedisMaps            = "%"
	RedisSets            = "~"
	RedisPushes          = ">"
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
	Value      string
	Expiry     time.Time
	MustExpire bool
}

type Set struct {
	Conn  net.Conn
	Env   map[string]EnvData
	Mutex IMutex
	Now   time.Time
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
	Conn  net.Conn
	Env   map[string]EnvData
	Mutex IMutex
	Now   time.Time
}

func (g *Get) Cmd(params []string) {
	cparams := findNextRedisSerialization(params)
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
	Conn  net.Conn
	Infos map[string]map[string]string
}

func (e *Info) mapToBulkString(m map[string]string, section string) string {
	// str := fmt.Sprintf("# %s\n\n", section)
	str := ""
	for k, v := range m {
		str += fmt.Sprintf("%s:%s\n", k, v)
	}
	return str
}

func (e *Info) Cmd(params []string) {
	str := ""
	cParams := findNextRedisSerialization(params)
	if len(cParams) > 0 {
		for _, v := range cParams {
			key := strings.ToLower(v)
			section, ok := e.Infos[key]
			if !ok {
				continue
			}
			str += e.mapToBulkString(section, strings.Title(key))
		}
	} else {
		for key, v := range e.Infos {
			str += e.mapToBulkString(v, strings.Title(key))
		}
	}
	str = str[:len(str) - 1]
	size := len(str)
	e.Conn.Write([]byte(
		fmt.Sprintf("$%d\r\n%s\r\n", size, str),
	))
}


type Echo struct {
	Conn net.Conn
}

func (e *Echo) Cmd(params []string) {
	var str string
	size := len(params)
	for i := 0; i < size; i++ {
		if i == size-1 {
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
