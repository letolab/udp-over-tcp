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

func initLogs(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

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
	senderChan := make(chan []byte)
	bufPool := NewBufPool(10)
	go SenderLoop(senderChan, bufPool)
	StartUDPListener(senderChan, bufPool)
}

func StartUDPListener(senderChan chan []byte, bufPool *BufPool) {
	addr := net.UDPAddr{
		Port: 9229,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	oob := make([]byte, 32)

	for {

		buf := bufPool.Borrow()

		n, _, _, _, err := conn.ReadMsgUDP(buf, oob)
		if err != nil {
			panic(err)
		}

		Info.Println("Got UDP Message:", buf[:n])

		go func(b []byte) {
			Info.Println("Message waiting to send to TCP Sender")
			senderChan <- b
			Info.Println("Message sent to TCP Sender")
		}(buf[:n])
	}

}

func SenderLoop(senderChan chan []byte, bufPool *BufPool) {
	for {
		select {
		case buf := <-senderChan:
			conn, err := net.Dial("tcp", "localhost:6000")
			if err != nil {
				Warning.Panicln("Could not dial tcp", err.Error())
				break
			}
			defer bufPool.Return(buf)
			Info.Println("Got UDP Message over TCP Sender Channel")
			_, err = conn.Write(buf)
			if err != nil {
				Warning.Printf("Error while writing %s\n", err.Error())
			}
		}
	}
}
