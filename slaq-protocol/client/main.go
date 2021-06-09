package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	quic "github.com/lucas-clemente/quic-go"
)

func main() {
	cp, err := x509.SystemCertPool()
	must(err)

	data, err := ioutil.ReadFile("../../minica/minica.pem")
	must(err)

	cp.AppendCertsFromPEM(data)

	tlsConf := &tls.Config{
		RootCAs:            cp,
		InsecureSkipVerify: false,
		NextProtos:         []string{"slaq"},
	}

	session, err := quic.DialAddr("quic.server:4242", tlsConf, nil)
	must(err)

	stream, err := session.OpenStreamSync(context.Background())
	must(err)

	message := "REG ilija"
	fmt.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	must(err)

	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	must(err)

	fmt.Printf("Client: Got '%s'\n", buf)
}

func must(err error) {
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
