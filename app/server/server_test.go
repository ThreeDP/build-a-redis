package server

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/define"
	"github.com/codecrafters-io/redis-starter-go/app/builtin"
)

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

func setupRedisServer() RedisServer {
	return RedisServer{
		Env: make(map[string]builtin.EnvData),
		Time: TTime{},
		Args: []string{}, 		//os.Args[1:]
		Infos: make(map[string]map[string]string),
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
		s := setupRedisServer()
		s.Args = []string{}
		expectedAddr := TListener{
			Address: &TAddr{Netw: "tcp",
			Host: "0.0.0.0",
			Port: "6379"},
		}

		err := s.Listen(TNetListen{})
		
		validateSeverAddr(t, &s, &expectedAddr, err)
	})

	t.Run("Test Listen in custom port", func(t *testing.T) {
		s := setupRedisServer()
		s.Args = []string{"--port", "6380"}
		expectedAddr := TListener{
			Address: &TAddr{Netw: "tcp",
			Host: "0.0.0.0",
			Port: "6380"},
		}

		err := s.Listen(TNetListen{})

		validateSeverAddr(t, &s, &expectedAddr, err)
	})

	t.Run("Test Listen with flags --port --replicaof", func(t *testing.T) {
		s := setupRedisServer()
		s.Args = []string{"--port", "6380", "--replicaof", "0.0.0.0", "6379" }
		expectedAddr := TListener{
			Address: &TAddr{Netw: "tcp",
			Host: "0.0.0.0",
			Port: "6380"},
		}

		err := s.Listen(TNetListen{})
		validateSeverAddr(t, &s, &expectedAddr, err)
	})
}

func TestHandleArgs(t *testing.T) {
	t.Run("Test HandleArgs with no flags", func(t *testing.T) {
		s := setupRedisServer()
		s.Args = []string{""}

		s.HandleArgs()

		checkInfosMap(t, s.Infos["replication"], "role", "master")
		checkInfosMap(t, s.Infos["server"], "port", define.DEFAULPORT)
	})

	t.Run("Test HandleArgs with flag --port", func(t *testing.T) {
		s := setupRedisServer()
		s.Args = []string{"--port", "7589"}

		s.HandleArgs()
		
		checkInfosMap(t, s.Infos["replication"], "role", "master")
		checkInfosMap(t, s.Infos["server"], "port", "7589")
	})

	t.Run("Test HandleArgs with flag --port --replicaof ", func(t *testing.T) {
		s := setupRedisServer()
		s.Args = []string{"--port", "8000", "--replicaof", "localhost", "8000"}

		s.HandleArgs()
		
		checkInfosMap(t, s.Infos["replication"], "role", "slave")
		checkInfosMap(t, s.Infos["server"], "port", "8000")
	})
}

func TestSetCommands(t *testing.T) {
	s := setupRedisServer()

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

