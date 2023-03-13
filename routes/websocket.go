package routes

import (
	"gbADMIN/controllers"

	"github.com/gin-gonic/gin"
)

func WebsocketRoute(router *gin.Engine, controller controllers.WebSocket) {
	router.GET("/connect", controller.WebSocketHandler)
}
