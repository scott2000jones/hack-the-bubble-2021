package main

import (
    "fmt"
    "net"
    "bufio"
	"time"
)

func main() {
    p :=  make([]byte, 2048)
    conn, err := net.Dial("udp", "127.0.0.1:22068")
    if err != nil {
        fmt.Printf("Some error %v", err)
        return
    }
	for {
		fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
		time.Sleep(time.Second * 1)
	}
    
    conn.Close()
}