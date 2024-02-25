package builtin

import (
	"net"
	"time"
)

type IMutex interface {
	Lock()
	TryLock() bool
	Unlock()
}

type Builtin interface {
	Response([]string)
	Request([]string)
	GetConn() net.Conn
	SetConn(net.Conn)
	GetTimeNow() time.Time
	SetTimeNow(time.Time)
}

type ITime interface {
	Now() time.Time
}

type Time struct{}

func (t Time) Now() time.Time {
	return time.Now()
}

type EnvData struct {
	Value      string
	Expiry     time.Time
	MustExpire bool
}

type Set struct {
	Env   map[string]EnvData
	Mutex IMutex
	Conn  net.Conn
	Now   time.Time
}

type Info struct {
	Infos map[string]map[string]string
	Conn  net.Conn
	Now   time.Time
}

type Echo struct {
	Conn net.Conn
	Now  time.Time
}

type Ping struct {
	Conn net.Conn
	Now  time.Time
}

type PSync struct {
	Conn net.Conn
	Now  time.Time
	Infos map[string]map[string]string
}
