package main

import (
	"io"
	"log"
	"net"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLogs(infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
func main() {
	initLogs(os.Stdout, os.Stdout, os.Stderr)
	TCPListner()
}

func TCPListner() {
	l, err := net.Listen("tcp", "localhost:6000")
	if err != nil {
		panic(err)
	}

	defer l.Close()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			Warning.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		Info.Println("Got new connection")
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 8192)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			Warning.Println("Got error while reading from conn:", err.Error())
			return
		}
		Info.Printf("Read %v from connection", buf[:n])
	}
}
