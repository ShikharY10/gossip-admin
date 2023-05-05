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

	service := schema.Service{
		Type: serviceType,
		Name: name,
		Host: c.Request.RemoteAddr,
		Port: port,
	}

	_service, err := ctrl._webSocketHandler(c.Writer, c.Request, service)
	if err != nil {
		log.Println(err.Error())
	}

	ctrl.Handler.Services = append(ctrl.Handler.Services, _service)

	// ctrl.Epoll.Lock.Lock()
	// ctrl.Epoll.Clients[name] = conn
	// ctrl.Epoll.Lock.Unlock()

	ctrl.Handler.Cache.RegisterNode(service.Type, service.Name)
}

func (ctrl *WebSocket) _webSocketHandler(w http.ResponseWriter, r *http.Request, service schema.Service) (*schema.Service, error) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	service.Conn = conn
	if err := ctrl.Epoll.Add(service); err != nil {
		log.Println("Failed to add connection: ", err)
		conn.Close()
		return nil, err
	}
	return &service, nil
}
