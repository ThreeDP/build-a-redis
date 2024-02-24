package builtin

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func (s *Set) Received(params []string) {
	flag := false
	cParams := parser.FindNextRedisSerialization(params)
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

func (s *Set) SetConn(conn net.Conn) {
	s.Conn = conn
}

func (s *Set) SetTimeNow(now time.Time) {
	s.Now = now
}