package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	quic "github.com/lucas-clemente/quic-go"
)

func main() {
	ln, err := quic.ListenAddr(":4242", generateTLSConfig(), nil)
	must(err)

	hub := newHub()
	go hub.run()

	for {
		sess, err := ln.Accept(context.Background())
		must(err)

		c := newConnection(
			sess,
			hub.commands,
			hub.registrations,
			hub.deregistrations,
		)
		go c.receive()
		go c.read()
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	tlsCert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	must(err)

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"slaq"},
	}
}

func must(err error) {
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
