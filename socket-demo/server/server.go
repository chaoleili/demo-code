package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	f := os.NewFile(3, "test.txt")
	defer f.Close()
	if f == nil {
		fmt.Println("f nil")
		return
	}
	l, err := net.FileListener(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConn(conn)
	}

}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
