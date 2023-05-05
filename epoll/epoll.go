package epoll

import (
	"fmt"
	"gbADMIN/schema"
	"reflect"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/sys/unix"
)

type EPOLL struct {
	FD            int
	Connections   map[int]*schema.Service
	Clients       map[string]*websocket.Conn
	Lock          *sync.RWMutex
	DataPipeline  chan []byte
	ClosePipeline chan schema.Service
}

func (e *EPOLL) websocketFD(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func InitiatEpoll() (*EPOLL, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	epoll := EPOLL{
		FD:            fd,
		Connections:   make(map[int]*schema.Service),
		Clients:       make(map[string]*websocket.Conn),
		Lock:          &sync.RWMutex{},
		DataPipeline:  make(chan []byte),
		ClosePipeline: make(chan schema.Service),
	}
	return &epoll, nil
}

func (e *EPOLL) Remove(conn websocket.Conn) error {
	fd := e.websocketFD(&conn)
	err := unix.EpollCtl(e.FD, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}

	service := e.Connections[fd]
	if service != nil {
		e.ClosePipeline <- *service
	}

	e.Lock.Lock()
	delete(e.Connections, fd)
	e.Lock.Unlock()

	return nil
}

func (e *EPOLL) Wait() ([]*websocket.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.FD, events, 100)
	if err != nil {
		return nil, err
	}
	e.Lock.RLock()
	defer e.Lock.RUnlock()
	var connections []*websocket.Conn
	for i := 0; i < n; i++ {
		service := e.Connections[int(events[i].Fd)]
		connections = append(connections, service.Conn)
	}
	return connections, nil
}

func (e *EPOLL) Add(service schema.Service) error {
	fd := e.websocketFD(service.Conn)
	err := unix.EpollCtl(e.FD,
		syscall.EPOLL_CTL_ADD,
		fd,
		&unix.EpollEvent{
			Events: unix.POLLIN | unix.POLLHUP,
			Fd:     int32(fd),
		},
	)
	if err != nil {
		return err
	}
	e.Lock.Lock()
	e.Connections[fd] = &service
	e.Lock.Unlock()
	fmt.Println("connections: ", e.Connections)
	return nil
}

func (e *EPOLL) StartEpollMonitoring() {
	for {
		connections, err := e.Wait()
		if err != nil {
			// ToDo: Log this error
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				continue
			}
			if _, msg, err := conn.ReadMessage(); err != nil {
				e.Remove(*conn)
			} else {
				e.DataPipeline <- msg
			}
		}
	}
}
