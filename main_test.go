package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

var sourceServer *RedisTestServer
var destServer *RedisTestServer

func ClearRedis() {
	sourceServer.conn.Do("flushdb")
	destServer.conn.Do("flushdb")
}

func StartTestServers() {
	fmt.Println("Starting redis...")
	sourceServer = NewRedisTestServer("6377")
	destServer = NewRedisTestServer("6277")
}

func StopTestServers() {
	fmt.Println("Stopping redis...")
	sourceServer.Stop()
	destServer.Stop()
}

func NewRedisTestServer(port string) *RedisTestServer {
	srv := &RedisTestServer{
		port: port,
		url:  fmt.Sprintf("127.0.0.1:%s", port),
	}

	srv.Start()

	return srv
}

type RedisTestServer struct {
	cmd  *exec.Cmd
	port string
	url  string
	conn redis.Conn
}

func (s *RedisTestServer) Start() {
	args := fmt.Sprintf("--port %s", s.port)
	s.cmd = exec.Command("redis-server", args)

	err := s.cmd.Start()
	time.Sleep(2 * time.Second)

	conn, err := redis.Dial("tcp", s.url)
	s.conn = conn

	if err != nil {
		panic("Could not start redis")
	}
}

func (s *RedisTestServer) Stop() {
	s.cmd.Process.Signal(syscall.SIGTERM)
	s.cmd.Process.Wait()
}

func TestMain(m *testing.M) {
	StartTestServers()

	result := m.Run()

	StopTestServers()
	os.Exit(result)
}
