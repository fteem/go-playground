package main

import "context"

type channel struct {
	name        string
	connections map[*connection]bool
}

func newChannel(name string) *channel {
	return &channel{
		name:        name,
		connections: make(map[*connection]bool),
	}
}

func (c *channel) broadcast(s string, m []byte) {
	msg := append([]byte(s), ": "...)
	msg = append(msg, m...)
	msg = append(msg, '\n')

	for cl := range c.connections {
		// Open stream and write to it
		go func() {
			stream, err := cl.session.OpenStreamSync(context.Background())
			if err == nil {
				stream.Write(msg)
			}
		}()
	}
}
