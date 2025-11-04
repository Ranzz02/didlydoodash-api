package whiteboardws

import (
	"fmt"

	"github.com/Stenoliv/didlydoodash_api/internal/db/models"
	"github.com/Stenoliv/didlydoodash_api/internal/ws"
	"github.com/gin-gonic/gin"
)

type Wbhandler struct {
	Hub *Hub
}

func (wbh *Wbhandler) HandleConnections(w *gin.Context) {
	ws, err := ws.WebsocketUpgrader.Upgrade(w.Writer, w.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	Client := &Client{Conn: ws, UserID: *models.CurrentUser, RoomID: w.Param("wbID"), Message: make(chan *WhiteboardMessage)}
	wbh.Hub.Register <- Client
	fmt.Println("New client connected")
	go Client.readMessage(wbh)
	go Client.writeMessage()

}
func NewHandler() *Wbhandler {
	hub := NewHub()
	return &Wbhandler{Hub: hub}
}
