package ws

import (
	"chatRoom/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[*Client]bool),
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsPage(c *gin.Context) {
	fmt.Println("websocket")
	conn, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	worker, err := utils.NewWorker(1)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := &Client{
		ID:     strconv.FormatInt(worker.GetId(), 10),
		Socket: conn,
		Send:   make(chan []byte),
	}

	Manager.Register <- client

	go client.Read()
	go client.Writer()
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.Clients[conn] = true
			jsonMsg, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			manager.Send(jsonMsg, conn)
		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				jsonMsg, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				manager.Send(jsonMsg, conn)
			}
		case message := <-manager.Broadcast:
			for conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if ignore != conn {
			manager.Broadcast <- message
		}
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			return
		}
		jsonMsg, _ := json.Marshal(&Message{Sender: c.ID, Content: string(message)})
		Manager.Broadcast <- jsonMsg
	}
}

func (c *Client) Writer() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
