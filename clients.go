package main

import "net"

type Client struct {
	Conn *net.TCPConn
}

func getClient(conn *net.TCPConn) *Client {
	return &Client{
		Conn: conn,
	}
}
