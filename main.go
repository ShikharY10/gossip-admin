package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gbADMIN/config"
	"gbADMIN/controllers"
	"gbADMIN/epoll"
	"gbADMIN/handler"
	"gbADMIN/routes"
	"gbADMIN/schema"
	"gbADMIN/utils"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func removeServices(channel chan string, handle *handler.Handler) {
	for serviceName := range channel {
		handle.RemoveService(serviceName)
		handle.Cache.RemoveNode(serviceName)
	}
}

func readDataFromServices(channel chan []byte, handler *handler.Handler) {

	var packet schema.Packet
	for job := range channel {
		err := json.Unmarshal(job, &packet)
		if err != nil {
			continue
		}

		if packet.Type == "log" {
			jsonEncoded := utils.Decode(packet.Message)
			var log schema.Log
			err := json.Unmarshal(jsonEncoded, &log)
			if err != nil {
				continue
			}

			err = handler.Database.SetNewLog(log)
			if err != nil {
				continue
			}
		}

	}
}

func userIO(handler *handler.Handler, port string) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	green := color.New(color.FgHiMagenta, color.Bold).PrintfFunc()
	for {
		green("G_ADMIN:" + port + " > ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.Trim(line, "\n")
		commands := strings.Split(line, " ")

		if commands[0] == "clear" {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				continue
			}
			continue
		} else if commands[0] == "services" {
			services := handler.Services
			// fmt.Println("     Name     |    Type    |    Host    |    Port    ")
			for _, service := range services {
				fmt.Println(service.Name + " | " + service.Type + " | " + service.Host + " | " + service.Port)
			}
		}

	}

}

func main() {
	ENV := config.LoadENV()

	db, err := config.ConnectToDBs(ENV)
	if err != nil {
		log.Fatal(err)
	}

	Epoll, err := epoll.InitiatEpoll()
	if err != nil {
		log.Fatal(err)
	}

	handle := &handler.Handler{
		Database: &handler.DataBaseHandler{
			Mongo: *db.MongoDB,
		},
		Cache: &handler.CacheHandler{
			RedisClient: db.RedisDB,
		},
		Env:      ENV,
		Services: []schema.Service{},
	}

	controller := controllers.WebSocket{
		Epoll:   Epoll,
		Handler: handle,
	}

	go readDataFromServices(Epoll.DataPipeline, handle)
	go removeServices(Epoll.ClosePipeline, handle)
	go Epoll.StartEpollMonitoring()

	router := gin.New()
	routes.WebsocketRoute(router, controller)

	utils.ShowSucces("[STARTING GOSSIP ADMINISTRATOR] | PORT => ["+ENV.AdminPort+"]", true)
	go userIO(handle, ENV.AdminPort)
	router.Run(":" + ENV.AdminPort)

}
