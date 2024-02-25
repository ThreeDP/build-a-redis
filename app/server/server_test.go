package server

import (
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/builtin"
	"github.com/codecrafters-io/redis-starter-go/app/define"
)

/*
Test Conn struct Mock
*/
type TConn struct {
	In  []byte
	Out []byte
}

func (c TConn) Close() error                       { return nil }
func (c TConn) SetDeadline(t time.Time) error      { return nil }
func (c TConn) SetReadDeadline(t time.Time) error  { return nil }
func (c TConn) SetWriteDeadline(t time.Time) error { return nil }

func (c TConn) Read(b []byte) (n int, err error) {
	copy(c.In, b)
	return len(b), nil
}

func (c TConn) Write(b []byte) (n int, err error) {
	copy(c.Out, b)
	return len(b), nil
}

func (c TConn) LocalAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}
}

func (c TConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}
}

func (c *TConn) NewConn() {
	c.In = make([]byte, define.BUFFERSIZE)
	c.Out = make([]byte, define.BUFFERSIZE)
}

/*
Mutex struct Mock
*/
type TMutex struct{}
func (m TMutex) Lock()         {}
func (m TMutex) Unlock()       {}
func (m TMutex) TryLock() bool { return false }

type TTime struct {}
func (t TTime) Now() time.Time {
	return time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)
}

type TNetListen struct {}
func (l TNetListen) Listen(network, address string) (net.Listener, error){
	addr := strings.Split(address, ":")
	return &TListener{Address: &TAddr{Netw: network, Host: addr[0], Port: addr[1]}}, nil
}

type TListener struct {
	Address net.Addr
}
func (l *TListener) Close() error {return nil}
func (l *TListener) Addr() net.Addr {return l.Address}
func (l *TListener) Accept() (net.Conn, error) {return nil, nil}

type TAddr struct {
	Netw string
	Host string
	Port string
}
func (a *TAddr) Network() string {return a.Netw}
func (a *TAddr) String() string {return a.Host + ":" + a.Port}

func checkInfosMap(t *testing.T,
	received map[string]string,
	key, expected string) {

	t.Helper()
	if received[key] != expected {
		t.Errorf("expected role %s, but has %s\n", expected, received["role"])
	}
}

func setupRedisServer(env map[string]builtin.EnvData) RedisServer {
	return RedisServer{
		Env: env,
		Time: TTime{},
		Args: []string{}, 		//os.Args[1:]
		Infos: make(map[string]map[string]string),
		Mutex: TMutex{},
	}
}

func validateSeverAddr(t *testing.T, s *RedisServer, expectedAddr *TListener, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Error: listen on addr: %v", s.Listener.Addr())
	}
	if expectedAddr.Addr().Network() != s.Listener.Addr().Network() {
		t.Errorf("Expected network %s, but has %s\n", expectedAddr.Addr().Network(), s.Listener.Addr().Network())
	}
	if expectedAddr.Addr().String() != s.Listener.Addr().String() {
		t.Errorf("Expected host and port %s, but has %s\n", expectedAddr.Addr().String(), s.Listener.Addr().String())
	}
}

func TestListenServer(t *testing.T) {
	
	t.Run("Test Listen in default", func(t *testing.T) {
		s := setupRedisServer(nil)
		s.Args = []string{}
		expectedAddr := TListener{
			Address: &TAddr{Netw: "tcp",
			Host: "localhost",
			Port: "6379"},
		}

		err := s.Listen(TNetListen{})
		
		validateSeverAddr(t, &s, &expectedAddr, err)
	})

	t.Run("Test Listen in custom port", func(t *testing.T) {
		s := setupRedisServer(nil)
		s.Args = []string{"--port", "6380"}
		expectedAddr := TListener{
			Address: &TAddr{Netw: "tcp",
			Host: "localhost",
			Port: "6380"},
		}

		err := s.Listen(TNetListen{})

		validateSeverAddr(t, &s, &expectedAddr, err)
	})

	t.Run("Test Listen with flags --port --replicaof", func(t *testing.T) {
		s := setupRedisServer(nil)
		s.Args = []string{"--port", "6380", "--replicaof", "localhost", "6379" }
		expectedAddr := TListener{
			Address: &TAddr{Netw: "tcp",
			Host: "localhost",
			Port: "6380"},
		}

		err := s.Listen(TNetListen{})
		validateSeverAddr(t, &s, &expectedAddr, err)
	})
}

func TestHandleArgs(t *testing.T) {
	t.Run("Test HandleArgs with no flags", func(t *testing.T) {
		s := setupRedisServer(nil)
		s.Args = []string{""}

		s.HandleArgs()

		checkInfosMap(t, s.Infos["replication"], "role", "master")
		checkInfosMap(t, s.Infos["server"], "port", define.DEFAULPORT)
	})

	t.Run("Test HandleArgs with flag --port", func(t *testing.T) {
		s := setupRedisServer(nil)
		s.Args = []string{"--port", "7589"}

		s.HandleArgs()
		
		checkInfosMap(t, s.Infos["replication"], "role", "master")
		checkInfosMap(t, s.Infos["server"], "port", "7589")
	})

	t.Run("Test HandleArgs with flag --port --replicaof ", func(t *testing.T) {
		s := setupRedisServer(nil)
		s.Args = []string{"--port", "8000", "--replicaof", "localhost", "8000"}

		s.HandleArgs()
		
		checkInfosMap(t, s.Infos["replication"], "role", "slave")
		checkInfosMap(t, s.Infos["server"], "port", "8000")
	})
}

func TestSetCommands(t *testing.T) {
	s := setupRedisServer(nil)

	t.Run("Test SetCommands with 'echo' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("echo", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command echo, but has not\n")
		}
		_, ok = b.(*builtin.Echo)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Echo{}, b)
		}
	})

	t.Run("Test SetCommands with 'ping' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("pinG", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command ping, but has not\n")
		}
		_, ok = b.(*builtin.Ping)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Ping{}, b)
		}
	})

	t.Run("Test SetCommands with 'info' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("iNfo", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command info, but has not\n")
		}
		_, ok = b.(*builtin.Info)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Info{}, b)
		}
	})

	t.Run("Test SetCommands with 'get' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("Get", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command get, but has not\n")
		}
		_, ok = b.(*builtin.Get)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Get{}, b)
		}
	})

	t.Run("Test SetCommands with 'set' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("SET", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command set, but has not\n")
		}
		_, ok = b.(*builtin.Set)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.Set{}, b)
		}
	})

	t.Run("Test SetCommands with 'replconf' command", func(t *testing.T) {
		s.SetCommands()
		b, ok := s.Commands("replconf", nil, s.Time.Now())
		if !ok {
			t.Errorf("Expected command replconf, but has not\n")
		}
		_, ok = b.(*builtin.ReplConf)
		if !ok {
			t.Errorf("Expected conn %T, but has %T\n", builtin.ReplConf{}, b)
		}
	})

	t.Run("Test SetCommands with unknown command", func(t *testing.T) {
		s.SetCommands()
		_, ok := s.Commands("unknown", nil, s.Time.Now())
		if ok {
			t.Errorf("Expected unknown command, but has not\n")
		}
	})
}

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


// func TestHandleResponse(t *testing.T) {
// 	tm := TTime{}
// 	s := setupRedisServer(map[string]builtin.EnvData{
// 		"Percy": {Value: "Jackson", Expiry: tm.Now(), MustExpire: false},
// 	})
// 	conn := TConn{}
// 	s.SetCommands()
// }