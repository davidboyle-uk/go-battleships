package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/dbx123/go-battleships/tcp"

	"github.com/rockwell-uk/go-logger/logger"
)

func HandleRequest(connId int, conn net.Conn, connMap *sync.Map) {
	defer func() {
		conn.Close()
		connMap.Delete(connId)
	}()

	reader := bufio.NewReader(conn)

	for {
		// Read inbound payload
		m, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		logger.Log(
			logger.LVL_INTERNAL,
			fmt.Sprintf("received %s", m),
		)

		// Parse
		p, err := tcp.ParseMessage(m)
		if err != nil {
			logger.Log(
				logger.LVL_ERROR,
				err.Error(),
			)
			return
		}

		responses, err := ProcessRequest(p)
		if err != nil {
			logger.Log(
				logger.LVL_ERROR,
				err.Error(),
			)
			return
		}

		// Write the response(s)
		for _, res := range responses {
			sendToClients(connMap, res)
		}
	}
}

func sendToClients(connMap *sync.Map, payload tcp.Proto) {
	connMap.Range(func(key, value interface{}) bool {
		if key == payload.Player || payload.Player == 0 {
			if conn, ok := value.(net.Conn); ok {
				if _, err := conn.Write([]byte(fmt.Sprintln(payload))); err != nil {
					logger.Log(
						logger.LVL_ERROR,
						fmt.Sprintf("error writing to client %s", err.Error()),
					)
				}
				logger.Log(
					logger.LVL_INTERNAL,
					fmt.Sprintf("sent %s", payload),
				)
			}
		}
		return true
	})
	if payload.Action == QUIT {
		os.Exit(0)
	}
}
