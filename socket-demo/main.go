package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"
)

func CreateUnixSocket(path string) (net.Listener, error) {
	// BSDs have a 104 limit
	if err := os.MkdirAll(filepath.Dir(path), 0660); err != nil {
		return nil, err
	}
	if err := unix.Unlink(path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return net.Listen("unix", path)
}

func main() {
	socket := filepath.Join("/tmp", "demo-code.sock")
	l, err := CreateUnixSocket(socket)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("./server/server")
	f, err := l.(*net.UnixListener).File()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	cmd.ExtraFiles = append(cmd.ExtraFiles, f)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.Dial("unix", l.Addr().String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	time.Sleep(100 * time.Second)
	return
}
