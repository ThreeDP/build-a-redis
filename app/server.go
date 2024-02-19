package main

import (
	"fmt"
	"net"
	"os"
)

func handleClient(cn net.Conn) {
	defer cn.Close()
	buf := make([]byte, 1024)

	n, err := cn.Read(buf)
	if err != nil {
		return
	}

	fmt.Println("received data", buf[:n])
	cn.Write([]byte("+PONG\r\n"))
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		cn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		handleClient(cn)
	}
}
