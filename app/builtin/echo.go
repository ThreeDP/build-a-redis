package builtin

import (
	"net"
	"time"
)

func (e *Echo) Response(params []string) {
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

func (e *Echo) Request(params []string) {
	
}

func (e *Echo) SetConn(conn net.Conn) {
	e.Conn = conn
}

func (e *Echo) SetTimeNow(now time.Time) {
	e.Now = now
}
