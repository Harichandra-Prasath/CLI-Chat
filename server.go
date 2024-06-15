package main

import (
	"fmt"
	"net"
	"sync"
)

type Config struct {
	Addr string
}

type Server struct {
	Cfg     *Config
	Clients []*Client
	lock    sync.RWMutex
}

func GetServer(cfg *Config) *Server {
	return &Server{
		Cfg:     cfg,
		Clients: make([]*Client, 0),
	}
}

func (s *Server) Serve() {
	l, err := net.Listen("tcp", s.Cfg.Addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		go s.handleConn(conn.(*net.TCPConn))
	}
}

func (s *Server) handleConn(conn *net.TCPConn) {
	client := getClient(conn)
	s.lock.Lock()
	s.Clients = append(s.Clients, client)
	s.lock.Unlock()
	for {
		buff := make([]byte, 512)
		_, err := client.Conn.Read(buff)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print(string(buff))
	}
}
