package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
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

	fmt.Printf("Server is listening on localhost%s\n", s.Cfg.Addr)
	l, err := net.Listen("tcp", s.Cfg.Addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Got a client connection at", conn.RemoteAddr().String())
		go s.handleConn(conn.(*net.TCPConn))
	}
}

func (s *Server) handleConn(conn *net.TCPConn) {
	var det ClientDetails

	conn.Write([]byte("Please give Client details in one line Json\n"))

	det_buff := make([]byte, 512)
	n, err := conn.Read(det_buff)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	det_buff = det_buff[:n-1]
	err = json.Unmarshal(det_buff, &det)
	if err != nil {
		fmt.Println(err.Error())
		conn.Write([]byte("Invalid Details\n"))
		return
	}

	client := getClient(conn, det)

	s.lock.Lock()
	s.Clients = append(s.Clients, client)
	s.lock.Unlock()

	conn.Write([]byte("You can announce or whisper messages now\n"))

	for {
		buff := make([]byte, 512)
		n, err := client.Conn.Read(buff)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		data := string(buff[:n])
		s.handleMessage(client, data)
	}
}

func (s *Server) announce(_message []string, announcer *Client) {

	message := strings.Join(_message, " ")
	s.lock.RLock()
	for _, client := range s.Clients {
		if client.Details.Name != announcer.Details.Name {
			go func(c *Client) {
				c.Conn.Write([]byte(fmt.Sprintf("%s Announced: %s\n", announcer.Details.Name, message)))
			}(client)
		}

	}
	s.lock.RUnlock()
}

func (s *Server) whisper(_message []string, sender *Client, reciever_name string) {

	message := strings.Join(_message, " ")

	var reciver *Client
	for _, client := range s.Clients {
		if client.Details.Name == reciever_name {
			reciver = client
			break
		}
	}
	if reciver != nil {
		reciver.Conn.Write([]byte(fmt.Sprintf("%s whispered to you: %s\n", sender.Details.Name, message)))
	} else {
		sender.Conn.Write([]byte("Sorry, No client with given username"))
	}
}
