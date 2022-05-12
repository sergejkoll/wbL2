package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ln, _ := net.Listen("tcp", ":8080")
	conn, _ := ln.Accept()
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)
	i := 0
	for {
		if i == 3 {
			conn.Close()
			return
		}
		msg := make([]byte, 1024)
		n, err := bufio.NewReader(conn).Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("Message Received: ", string(msg[:n]))
		conn.Write(msg[:n])
		i++
	}
}
