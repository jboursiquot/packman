package main

import (
	"bufio"
	"io"
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
	OK = []byte("OK\n")
	// FAIL means a business rule violated
	FAIL = []byte("FAIL\n")
	// ERROR means bad message or otherwise unparsable
	ERROR = []byte("ERROR\n")

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
			if err != io.EOF {
				log.Println(err)
			}
			return
		}
		message = strings.TrimSpace(message)
		// log.Printf("%v | Message: %v", conn, message)

		cmd, err := packman.CommandFromMessage(message)
		if err != nil {
			log.Printf("REQ=%v, RES=%v", message, string(ERROR))
			_, wErr := conn.Write(ERROR)
			if wErr != nil {
				log.Printf("REQ=%v, %v", message, wErr)
			}
			continue
		}

		_, err = packman.ProcessCommand(cmd, &idxr)
		if err != nil {
			log.Printf("REQ=%v  RES=%v", message, string(FAIL))
			_, wErr := conn.Write(FAIL)
			if wErr != nil {
				log.Printf("REQ=%v, %v", message, wErr)
			}
		} else {
			log.Printf("REQ=%v, RES=%v", message, string(OK))
			_, wErr := conn.Write(OK)
			if wErr != nil {
				log.Printf("REQ=%v, %v", message, wErr)
			}
		}

	}
}
