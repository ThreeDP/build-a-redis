package builtin

import (
	"fmt"
	"io/ioutil"
	"net"
	"encoding/base64"
)

func SendRDBData(conn net.Conn) {
	rdb := "./app/rdb/data.rdb"
	b64, err := ioutil.ReadFile(rdb)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		conn.Write([]byte("-ERR error reading file\n"))
		return
	}

	decode, err := base64.StdEncoding.DecodeString(string(b64))
	if err != nil {
		fmt.Printf("Error decoding file: %s\n", err)
		conn.Write([]byte("-ERR error decoding file\n"))
		return
	}
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s", len(decode), decode)))
}