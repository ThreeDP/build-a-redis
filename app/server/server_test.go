package server

import (
	"net"
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

func NewDefaultTRedisServer() RedisServer {
	timeNow := time.Date(2009, 1, 1, 12, 0, 0, 0, time.UTC)
	return RedisServer{
		Env: map[string]builtin.EnvData {
			"Percy": {Value: "Jackson", Expiry: timeNow, MustExpire: false},
			"Key":   {Value: "Value", Expiry: timeNow.Add(-10 * time.Millisecond), MustExpire: true},
		},
		Infos: map[string]map[string]string{
			"replication": {
				"role":					"master",
				"master_replid":		"8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
				"master_repl_offset":	"0",
				"master_host":			"localhost",
				"master_port":			define.DEFAULPORT,
				"master_link_status":	"down",
			},
			"server": {
				"port":					define.DEFAULPORT,
				"listener0":			"name=tcp,bind=*,bind=localhost,port=" + define.DEFAULPORT,
			},
		},
		Mutex: TMutex{},
		Time: TTime{},
		Args: []string{},
		Idx: 0,
		Listener: nil,
		Action: nil,
		Commands: nil,
	}
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





// func TestHandleResponse(t *testing.T) {
// 	tm := TTime{}
// 	s := setupRedisServer(map[string]builtin.EnvData{
// 		"Percy": {Value: "Jackson", Expiry: tm.Now(), MustExpire: false},
// 	})
// 	conn := TConn{}
// 	s.SetCommands()
// }