package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrNoConnections = errors.New("no connections")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("connections error: %w", err)
	}

	c.conn = conn

	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Send() error {
	_, err := io.Copy(c.conn, c.in)

	return err
}

func (c *client) Receive() error {
	if c.conn == nil {
		return ErrNoConnections
	}

	_, err := io.Copy(c.out, c.conn)

	return err
}
