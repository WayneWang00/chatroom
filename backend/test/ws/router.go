package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Init() {
	fmt.Println("ws router")
	go Manager.Start()

	router := gin.Default()
	router.GET("/ws", WsPage)
	router.Run(":8088")
}
