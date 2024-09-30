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
	Details ClientDetails
	Conn    *net.TCPConn
}

func getClient(conn *net.TCPConn, details ClientDetails) *Client {
	return &Client{
		Conn:    conn,
		Details: details,
	}
}

func (s *Server) handleMessage(client *Client, data string) {

	data = strings.Replace(data, "\r\n", "", -1)
	parts := strings.Split(data, " ")
	mode := parts[0]
	switch mode {
	case "/announce":
		if len(parts) < 2 {
			client.Conn.Write([]byte("Please add the announcement\n"))
			return
		}
		s.announce(parts[1:], client)
	case "/whisper":
		if len(parts) < 3 {
			client.Conn.Write([]byte("Need username and message to whisper\n"))
			return
		}
		s.whisper(parts[2:], client, parts[1])
	default:
		fmt.Println("didnt understand")
		return
	}
}
