package controllers

import (
	"gbADMIN/epoll"
	"gbADMIN/handler"
	"gbADMIN/schema"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
	Epoll   *epoll.EPOLL
	Handler *handler.Handler
}

func (ctrl *WebSocket) WebSocketHandler(c *gin.Context) {
	serviceType := c.Query("type")
	name := c.Query("name")
	port := c.Query("port")

	conn, err := ctrl._webSocketHandler(c.Writer, c.Request)
	if err != nil {
		log.Println(err.Error())
	}

	service := schema.Service{
		Type: serviceType,
		Name: name,
		Host: c.Request.RemoteAddr,
		Port: port,
	}

	ctrl.Handler.Services = append(ctrl.Handler.Services, service)

	ctrl.Epoll.Lock.Lock()
	ctrl.Epoll.Clients[name] = conn
	ctrl.Epoll.Lock.Unlock()
}

func (ctrl *WebSocket) _webSocketHandler(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	if err := ctrl.Epoll.Add(conn); err != nil {
		log.Println("Failed to add connection: ", err)
		conn.Close()
		return nil, err
	}
	return conn, nil
}
