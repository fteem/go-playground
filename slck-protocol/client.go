package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
)

var (
	DELIMITER = []byte(`\r\n`)
)

type client struct {
	conn       net.Conn
	outbound   chan<- command
	register   chan<- *client
	deregister chan<- *client
	username   string
}

func newClient(conn net.Conn, o chan<- command, r chan<- *client, d chan<- *client) *client {
	return &client{
		conn:       conn,
		outbound:   o,
		register:   r,
		deregister: d,
	}
}

func (c *client) read() error {
	for {
		msg, err := bufio.NewReader(c.conn).ReadBytes('\n')
		if err == io.EOF {
			// Connection closed, deregister client
			c.deregister <- c
			return nil
		}

		if err != nil {
			return err
		}

		c.handle(msg)
	}
	return nil
}

func (c *client) handle(message []byte) {
	cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))
	args := bytes.TrimSpace(bytes.TrimPrefix(message, cmd))

	switch string(cmd) {
	case "REG":
		if err := c.reg(args); err != nil {
			c.err(err)
		}
	case "JOIN":
		if err := c.join(args); err != nil {
			c.err(err)
		}
	case "LEAVE":
		if err := c.leave(args); err != nil {
			c.err(err)
		}
	case "MSG":
		if err := c.msg(args); err != nil {
			c.err(err)
		}
	case "CHNS":
		c.chns()
	case "USRS":
		c.usrs()
	default:
		c.err(fmt.Errorf("Unknown command %s", cmd))
	}
}

func (c *client) reg(args []byte) error {
	u := bytes.TrimSpace(args)
	if u[0] != '@' {
		return fmt.Errorf("Username must begin with @")
	}
	if len(u) == 0 {
		return fmt.Errorf("Username cannot be blank")
	}

	c.username = string(u)
	c.register <- c

	return nil
}

func (c *client) join(args []byte) error {
	channelID := bytes.TrimSpace(args)
	if channelID[0] != '#' {
		return fmt.Errorf("ERR Channel ID must begin with #")
	}

	c.outbound <- command{
		recipient: string(channelID),
		sender:    c.username,
		id:        JOIN,
	}
	return nil
}

func (c *client) leave(args []byte) error {
	channelID := bytes.TrimSpace(args)
	if channelID[0] == '#' {
		return fmt.Errorf("ERR channelID must start with '#'")
	}

	c.outbound <- command{
		recipient: string(channelID),
		sender:    c.username,
		id:        LEAVE,
	}
	return nil
}

func (c *client) msg(args []byte) error {
	args = bytes.TrimSpace(args)
	if args[0] != '#' && args[0] != '@' {
		return fmt.Errorf("recipient must be a channel ('#name') or user ('@user')")
	}

	recipient := bytes.Split(args, []byte(" "))[0]
	if len(recipient) == 0 {
		return fmt.Errorf("recipient must have a name")
	}

	args = bytes.TrimSpace(bytes.TrimPrefix(args, recipient))
	l := bytes.Split(args, DELIMITER)[0]
	length, err := strconv.Atoi(string(l))
	if err != nil {
		return fmt.Errorf("body length must be present")

	}
	if length == 0 {
		return fmt.Errorf("body length must be at least 1")
	}

	padding := len(l) + len(DELIMITER) // Size of the body length + the delimiter
	body := args[padding : padding+length]

	c.outbound <- command{
		recipient: string(recipient),
		sender:    c.username,
		body:      body,
		id:        MSG,
	}

	return nil
}

func (c *client) chns() {
	c.outbound <- command{
		sender: c.username,
		id:     CHNS,
	}
}

func (c *client) usrs() {
	c.outbound <- command{
		sender: c.username,
		id:     USRS,
	}
}

func (c *client) err(e error) {
	c.conn.Write([]byte("ERR " + e.Error() + "\n"))
}
