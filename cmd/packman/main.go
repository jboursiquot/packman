package main

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/jboursiquot/packman"
)

const (
	port = ":8080"
)

var (
	// OK means we're all good
	OK = "OK"
	// FAIL means a business rule violated
	FAIL = "FAIL"
	// ERROR means bad message or otherwise unparsable
	ERROR = "ERROR"

	idxr = packman.NewIndexer(nil)
)

func main() {
	log.Println("Starting server...")

	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := ln.Accept()
		log.Printf("conn, err -> %v, %v", conn, err)
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		message = strings.TrimSpace(message)
		log.Printf("Message Received: '%v'", message)

		cmd, err := packman.CommandFromMessage(message)
		if err != nil {
			// log.Printf("Error: %v | Sending %v", err, []byte(ERROR+"\n"))
			_, wErr := conn.Write([]byte(ERROR + "\n"))
			if wErr != nil {
				log.Printf("Failed to send message back: %v. Skipping.", wErr)
			}
		}

		_, err = packman.ProcessCommand(cmd, &idxr)
		// log.Printf("ProcessCommand(%v, %v) -> (%v, %v)", cmd, &idxr, res, err)
		if err != nil {
			// log.Printf("Error: %v | Sending %v", err, []byte(FAIL+"\n"))
			_, wErr := conn.Write([]byte(FAIL + "\n"))
			if wErr != nil {
				log.Printf("Failed to send message back: %v. Skipping.", wErr)
			}
		}

		_, wErr := conn.Write([]byte(OK + "\n"))
		// log.Printf("conn.Write(%v) -> %d, %v", []byte(OK+"\n"), written, err)
		if wErr != nil {
			log.Printf("Failed to send message back: %v. Skipping.", wErr)
		}

	}
}
