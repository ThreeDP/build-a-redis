package server

import (
	"testing"
	"reflect"

	"github.com/codecrafters-io/redis-starter-go/app/builtin"
	"github.com/codecrafters-io/redis-starter-go/app/define"
)

func TestHandleRequest(t *testing.T) {
	tm := TTime{}
	s := setupRedisServer(map[string]builtin.EnvData{
		"Percy": {Value: "Jackson", Expiry: tm.Now(), MustExpire: false},
	})
	conn := TConn{}
	s.SetCommands()

	t.Run("Test *2\\r\\n$4\\r\\necho\\r\\n$2\\r\\nhi\\r\\n command", func(t *testing.T) {
		conn.NewConn()
		expected := make([]byte, define.BUFFERSIZE)
		buf := "*2\r\n$4\r\necho\r\n$2\r\nhi\r\n"
		copy(expected, "+hi\r\n")

		s.HandleRequest(&conn, buf)

		if !reflect.DeepEqual(conn.Out, expected) {
			t.Errorf("Expected '%v', but has '%v'\n", string(expected), string(conn.Out))
		}
	})

	t.Run("Test *2\\r\\n$3\\r\\nget\\r\\n$5\\r\\nPercy\\r\\n command", func(t *testing.T) {
		conn.NewConn()
		expected := make([]byte, define.BUFFERSIZE)
		buf := "*2\r\n$3\r\nget\r\n$5\r\nPercy\r\n"
		copy(expected, "$7\r\nJackson\r\n")

		s.HandleRequest(&conn, buf)
		if !reflect.DeepEqual(conn.Out, expected) {
			t.Errorf("Expected '%v', but has '%v'\n", string(expected), string(conn.Out))
		}
	})

	t.Run("Test *2\\r\\n$3\\r\\nget\\r\\n$7\\r\\nunknown\\r\\n command", func(t *testing.T) {	
		conn.NewConn()
		expected := make([]byte, define.BUFFERSIZE)
		buf := "*2\r\n$3\r\nget\r\n$7\r\nunknown\r\n"
		copy(expected, "$-1\r\n")

		s.HandleRequest(&conn, buf)
		if !reflect.DeepEqual(conn.Out, expected) {
			t.Errorf("Expected '%v', but has '%v'\n", string(expected), string(conn.Out))
		}
	})

	t.Run("Test *3\\r\\n$3\\r\\nset\\r\\n$6\\r\\junior\\r\\n$4\r\\nluis\r\n command", func(t *testing.T) {
		conn.NewConn()
		expected := make([]byte, define.BUFFERSIZE)
		buf := "*3\r\n$3\r\nset\r\n$6\r\njunior\r\n$4\r\nluis\r\n"
		copy(expected, "+OK\r\n")

		s.HandleRequest(&conn, buf)
		if !reflect.DeepEqual(conn.Out, expected) {
			t.Errorf("Expected '%v', but has '%v'\n", string(expected), string(conn.Out))
		}
	})

	t.Run("Test *2\\r\\n$4\\r\\nunknown\\r\\n$3\\r\\n123\\r\\n command", func(t *testing.T) {
		conn.NewConn()
		expected := make([]byte, define.BUFFERSIZE)
		buf := "*2\r\n$7\r\nunknown\r\n$3\r\n123\r\n"
		copy(expected, "-ERR unknown command 'unknown'\r\n")

		s.HandleRequest(&conn, buf)
		if !reflect.DeepEqual(conn.Out, expected) {
			t.Errorf("Expected '%v', but has '%v'\n", string(expected), string(conn.Out))
		}
	})
}
