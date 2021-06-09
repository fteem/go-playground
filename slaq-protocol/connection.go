package main

import (
	"context"

	quic "github.com/lucas-clemente/quic-go"
)

var (
	DELIMITER = []byte(`\r\n`)
)

type connection struct {
	session    quic.Session
	outbound   chan<- command
	register   chan<- *connection
	deregister chan<- *connection
	inbound    chan command
	streams    map[quic.StreamID]quic.Stream
	username   string
}

func newConnection(s quic.Session, o chan<- command, r chan<- *connection, d chan<- *connection) *connection {
	return &connection{
		session:    s,
		outbound:   o,
		register:   r,
		deregister: d,
		streams:    make(map[quic.StreamID]quic.Stream),
	}
}

func (conn *connection) read() error {
	for {
		stream, err := conn.session.AcceptStream(context.Background())
		if err != nil {
			conn.deregister <- conn
			return err
		}
		conn.streams[stream.StreamID()] = stream

		go handle(stream, conn.inbound, conn.username)
	}
	return nil
}

func (conn *connection) receive() {
	for cmd := range conn.inbound {
		switch cmd.id {
		case REG:
			conn.username = cmd.sender
			conn.register <- conn
		default:
			conn.outbound <- cmd
		}
	}
}
