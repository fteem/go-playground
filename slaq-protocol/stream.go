package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"

	quic "github.com/lucas-clemente/quic-go"
)

func handle(stream quic.Stream, outbound chan command, username string) {
	for {
		msg := make([]byte, 64)
		_, err := io.ReadFull(stream, msg)
		if err != nil {
			log.Printf("Error: %v\n", err)
			break
		}
		fmt.Println(msg)

		cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(msg, []byte(" "))[0]))
		args := bytes.TrimSpace(bytes.TrimPrefix(msg, cmd))

		switch string(cmd) {
		case "REG":
			if err := register(outbound, stream.StreamID(), args); err != nil {
				stream.Write([]byte("ERR " + err.Error() + "\n"))
			}
		case "JOIN":
			if err := join(outbound, args, username); err != nil {
				stream.Write([]byte("ERR " + err.Error() + "\n"))
			}
		case "LEAVE":
			if err := leave(outbound, args, username); err != nil {
				stream.Write([]byte("ERR " + err.Error() + "\n"))
			}
		case "MSG":
			if err := message(outbound, stream.StreamID(), args, username); err != nil {
				stream.Write([]byte("ERR " + err.Error() + "\n"))
			}
		case "CHNS":
			chns(stream.StreamID(), outbound, username)
		case "USRS":
			usrs(stream.StreamID(), outbound, username)
		default:
			stream.Write([]byte("ERR Unknown command \n"))
		}
	}
}

func register(outbound chan command, sid quic.StreamID, args []byte) error {
	u := bytes.TrimSpace(args)
	if u[0] != '@' {
		return fmt.Errorf("Username must begin with @")
	}
	if len(u) == 0 {
		return fmt.Errorf("Username cannot be blank")
	}

	outbound <- command{
		sender:   string(u),
		id:       REG,
		streamID: sid,
	}

	return nil
}

func join(outbound chan command, args []byte, username string) error {
	channelID := bytes.TrimSpace(args)
	if channelID[0] != '#' {
		return fmt.Errorf("ERR Channel ID must begin with #")
	}

	outbound <- command{
		recipient: string(channelID),
		sender:    username,
		id:        JOIN,
	}
	return nil
}

func leave(outbound chan command, args []byte, username string) error {
	channelID := bytes.TrimSpace(args)
	if channelID[0] == '#' {
		return fmt.Errorf("ERR channelID must start with '#'")
	}

	outbound <- command{
		recipient: string(channelID),
		sender:    username,
		id:        LEAVE,
	}
	return nil
}

func message(outbound chan command, sid quic.StreamID, args []byte, username string) error {
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

	outbound <- command{
		recipient: string(recipient),
		sender:    username,
		body:      body,
		id:        MSG,
		streamID:  sid,
	}

	return nil
}

func chns(sid quic.StreamID, outbound chan command, username string) {
	outbound <- command{
		sender:   username,
		id:       CHNS,
		streamID: sid,
	}
}

func usrs(sid quic.StreamID, outbound chan command, username string) {
	outbound <- command{
		sender:   username,
		id:       USRS,
		streamID: sid,
	}
}
