package builtin

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func (e *Info) mapToBulkString(m map[string]string, section string) string {
	str := fmt.Sprintf("## %s\n\n", section)
	for k, v := range m {
		str += fmt.Sprintf("%s:%s\n", k, v)
	}
	return str
}

func (e *Info) Response(params []string) {
	str := ""
	cParams := parser.FindNextRedisSerialization(params)
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
	str = str[:len(str)-1]
	size := len(str)
	e.Conn.Write([]byte(
		fmt.Sprintf("$%d\r\n%s\r\n", size, str),
	))
}

func (e *Info) Request(params []string) {
}

func (i *Info) SetConn(conn net.Conn) {
	i.Conn = conn
}

func (i *Info) GetConn() net.Conn {
	return i.Conn
}

func (i *Info) SetTimeNow(now time.Time) {
	i.Now = now
}

func (i *Info) GetTimeNow() time.Time {
	return i.Now
}