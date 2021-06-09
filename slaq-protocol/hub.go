package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type hub struct {
	channels        map[string]*channel
	connections     map[string]*connection
	commands        chan command
	deregistrations chan *connection
	registrations   chan *connection
}

func newHub() *hub {
	return &hub{
		registrations:   make(chan *connection),
		deregistrations: make(chan *connection),
		connections:     make(map[string]*connection),
		channels:        make(map[string]*channel),
		commands:        make(chan command),
	}
}

func (h *hub) run() {
	for {
		select {
		case conn := <-h.registrations:
			h.register(conn)
		case conn := <-h.deregistrations:
			h.deregister(conn)
		case cmd := <-h.commands:
			switch cmd.id {
			case JOIN:
				h.joinChannel(cmd)
			case LEAVE:
				h.leaveChannel(cmd)
			case MSG:
				h.message(cmd)
			case USRS:
				h.listUsers(cmd)
			case CHNS:
				h.listChannels(cmd)
			default:
				// Freak out?
			}
		}
	}
}

func (h *hub) register(c *connection) {
	if _, exists := h.connections[c.username]; exists {
		c.username = fmt.Sprintf("%s-%d", c.username, rand.Intn(1000))
		h.connections[c.username] = c
	} else {
		h.connections[c.username] = c
	}
}

func (h *hub) deregister(c *connection) {
	if _, exists := h.connections[c.username]; exists {
		delete(h.connections, c.username)

		for _, channel := range h.channels {
			delete(channel.connections, c)
		}
	}
}

func (h *hub) joinChannel(cmd command) {
	if conn, ok := h.connections[cmd.sender]; ok {
		if channel, ok := h.channels[cmd.recipient]; ok {
			// Channel exists, join
			channel.connections[conn] = true
		} else {
			// Channel doesn't exists, create and join
			ch := newChannel(cmd.recipient)
			ch.connections[conn] = true
			h.channels[cmd.recipient] = ch
		}
		if stream, ok := conn.streams[cmd.streamID]; ok {
			stream.Write([]byte("OK\n"))
		}
	}
}

func (h *hub) leaveChannel(cmd command) {
	if conn, ok := h.connections[cmd.sender]; ok {
		if channel, ok := h.channels[cmd.recipient]; ok {
			delete(channel.connections, conn)

			if stream, ok := conn.streams[cmd.streamID]; ok {
				stream.Write([]byte("OK\n"))
			}
		}
	}
}

func (h *hub) message(cmd command) {
	if sender, ok := h.connections[cmd.sender]; ok {
		switch cmd.recipient[0] {
		case '#':
			if channel, ok := h.channels[cmd.recipient]; ok {
				if _, ok := channel.connections[sender]; ok {
					channel.broadcast(sender.username, cmd.body)
				}
			} else {
				if stream, ok := sender.streams[cmd.streamID]; ok {
					stream.Write([]byte("ERR no such channel"))
				}
			}
		case '@':
			if user, ok := h.connections[cmd.recipient]; ok {
				if stream, ok := user.streams[cmd.streamID]; ok {
					msg := append([]byte(sender.username+": "), cmd.body...)
					msg = append(msg, '\n')

					stream.Write(msg)
				}
			} else {
				if stream, ok := sender.streams[cmd.streamID]; ok {
					stream.Write([]byte("ERR no such user"))
				}
			}
		default:
			if stream, ok := sender.streams[cmd.streamID]; ok {
				stream.Write([]byte("ERR MSG command"))
			}
		}
	}
}

func (h *hub) listUsers(cmd command) {
	if sender, ok := h.connections[cmd.sender]; ok {
		var names []string

		for _, c := range h.connections {
			names = append(names, "@"+c.username+" ")
		}

		resp := strings.Join(names, ", ")

		if stream, ok := sender.streams[cmd.streamID]; ok {
			stream.Write([]byte(resp + "\n"))
		}
	}
}

func (h *hub) listChannels(cmd command) {
	if sender, ok := h.connections[cmd.sender]; ok {
		if stream, ok := sender.streams[cmd.streamID]; ok {
			if len(h.channels) == 0 {
				stream.Write([]byte("ERR no channels found\n"))
			} else {
				var names []string
				for c := range h.channels {
					names = append(names, "#"+c+" ")
				}

				resp := strings.Join(names, ", ")
				stream.Write([]byte(resp + "\n"))
			}
		}
	}
}
