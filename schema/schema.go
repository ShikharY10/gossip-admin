package schema

import "github.com/gorilla/websocket"

type Packet struct {
	NodeName string `json:"name"`
	Type     string `json:"type"`
	Message  string `json:"message"`
}

type Log struct {
	TimeStamp   string `bson:"timeStamp" json:"timeStamp"`
	ServiceType string `bson:"serviceType" json:"serviceType"`
	Type        string `bson:"type" json:"type"`
	FileName    string `bson:"fileName" json:"fileName"`
	LineNumber  int    `bson:"lineNumber" json:"lineNumber"`
	Message     string `bson:"errorMessage" json:"errorMessage"`
}

type Service struct {
	Type string
	Name string
	Host string
	Port string
	Conn *websocket.Conn
}
