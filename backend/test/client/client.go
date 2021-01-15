package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

var addr = flag.String("addr", "192.168.56.101:8088", "http service address")

func main() {
	fmt.Println("test client")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read err:", err)
			return
		}

		fmt.Printf("read message:%+v\n", string(message))
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(2 * time.Second)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}
