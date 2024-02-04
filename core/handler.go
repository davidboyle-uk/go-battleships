package core

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"go-battleships/logger"
	"go-battleships/tcp"
)

func HandleRequest(ct context.Context, connId int, conn net.Conn, connMap *sync.Map) {
	defer func() {
		conn.Close()
		connMap.Delete(connId)
	}()

	reader := bufio.NewReader(conn)

	for {
		select {
		case <-ct.Done():
			sendToClients(connMap, tcp.Proto{
				Action: QUIT,
			})
			os.Exit(1)
		default:
			// Read inbound payload
			message, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					logger.Log(
						logger.LVL_ERROR,
						fmt.Sprintf("error reading input %s", err.Error()),
					)
				}
				return
			}

			// String message
			m := string(message)
			logger.Log(
				logger.LVL_DEBUG,
				fmt.Sprintf("recieved packet: %s", m),
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
					logger.LVL_DEBUG,
					fmt.Sprintf("sent %s", payload),
				)
			}
		}
		return true
	})
}
