package main

import (
	"fmt"
	"net"
	"strings"
)

type ClientDetails struct {
	Name string
}

type Client struct {
	Details *ClientDetails
	Conn    *net.TCPConn
}

func getClient(conn *net.TCPConn, details *ClientDetails) *Client {
	return &Client{
		Conn:    conn,
		Details: details,
	}
}

func (s *Server) handleMessage(client *Client, data string) {

	data = strings.Replace(data, "\n", "", -1)
	parts := strings.Split(data, " ")
	mode := parts[0]
	switch mode {
	case "/announce":
		s.announce(parts[1], client)
	case "/whisper":
		fmt.Println("Someone whispered")
	default:
		fmt.Println("didnt understand")
		return
	}
}
