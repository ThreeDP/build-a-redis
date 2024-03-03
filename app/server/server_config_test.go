package server

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/define"
)

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