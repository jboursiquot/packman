package main

import (
	"bufio"
	"log"
	"net"
)

const (
	port = ":8080"
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
			log.Println(err)
			continue
		}
		log.Printf("Message Received: %v\n", string(message))
	}
}
