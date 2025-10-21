package modbus

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	conn    net.Conn
	timeout time.Duration
}

func NewClient(ip string, port int, timeout time.Duration) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		timeout: timeout,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
	if c.conn != nil {
		c.conn.SetDeadline(time.Now().Add(timeout))
	}
}
