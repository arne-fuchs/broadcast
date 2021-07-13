package server

import (
	"github.com/iotaledger/goshimmer/packages/tangle"

	"github.com/iotaledger/hive.go/events"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/netutil/buffconn"

	"fmt"
	"io"
	"net"
	"strings"
)

type Connection struct {
	bufferedConn  *buffconn.BufferedConnection
	log           *logger.Logger
	active 	      bool
}


var connectionList = [256]Connection{}
var index = 0

// Listen starts a TCP listener and starts a Connection for each accepted connection
func Listen(bindAddress string, log *logger.Logger, shutdownSignal <-chan struct{}) error {

	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return fmt.Errorf("failed to start Broadcast daemon: %w", err)
	}

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			log.Infof("Started connection: %s", connection.RemoteAddr().String())
			go handleConnection(connection, log, shutdownSignal)
		}
	}()

	go func() {
		defer func(listener net.Listener) {
			err := listener.Close()
			if err != nil {
				log.Infof("Error closing listener: %s",err)
			}
		}(listener)

		<-shutdownSignal

		log.Infof("Closing Broadcast server...")
		for i := 0; i < index; i++ {
			connectionList[i].active = false
		}
		log.Infof("Closing Broadcast server... \tDone")
	}()

	return nil
}

func handleConnection(conn net.Conn, log *logger.Logger, shutdownSignal <-chan struct{}) {

	connectionList[index] = Connection{
		bufferedConn:  buffconn.NewBufferedConnection(conn, tangle.MaxMessageSize),
		log:           log,
		active: 	true,
	}
	index++
	bufferedConnDataReceived, bufferedConnClosed := connectionList[index-1].readLoop()

	select {
	case data := <-bufferedConnDataReceived:
		//No input required. For debugging it will be printed
		connectionList[index-1].log.Infof("Data received:%v",data)
		connectionList[index-1].active = false
	case <-shutdownSignal:
		connectionList[index-1].log.Infof("Shutdown signal received")
		for i := 0; i < index; i++ {
			connectionList[i].active = false
		}
		return
	case <-bufferedConnClosed:
		connectionList[index-1].log.Errorf("Connection lost")
		connectionList[index-1].active = false
		return
	}
}

func (connection *Connection) readLoop() (chan []byte, chan bool) {
	bufferedConnDataReceived := make(chan []byte)
	bufferedConnClosed := make(chan bool)

	go func() {
		{
			cl := events.NewClosure(func() { close(bufferedConnClosed) })
			connection.bufferedConn.Events.Close.Attach(cl)
			defer connection.bufferedConn.Events.Close.Detach(cl)
		}

		{
			cl := events.NewClosure(func(data []byte) {
				d := make([]byte, len(data))
				copy(d, data)
				bufferedConnDataReceived <- d
			})
			connection.bufferedConn.Events.ReceiveMessage.Attach(cl)
			defer connection.bufferedConn.Events.ReceiveMessage.Detach(cl)
		}

		if err := connection.bufferedConn.Read(); err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "Use of closed network connection") {
				connection.log.Warnw("Buffered connection read error", "err", err)
				connection.active = false
			}
		}
	}()

	return bufferedConnDataReceived, bufferedConnClosed
}

func Broadcast(data []byte) {
	for i := 0; i < index; i++ {
		if connectionList[i].active{
			_, err := connectionList[i].bufferedConn.Write(data)
			if err != nil {
				connectionList[i].log.Debugf("Error writing on connection: %s",err)
				connectionList[i].active = false
				return
			}
		}
	}
	//Tidy up array of unused connections
	if index == 250{
		newIndex := 0
		var newConnectionList = [256]Connection{}
		for i := 0; i < index; i++ {
			if connectionList[i].active{
				newConnectionList[newIndex] = connectionList[i]
				newIndex++
			}
		}
		index = newIndex
		connectionList = newConnectionList
	}
}
