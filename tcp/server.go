package tcp

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/rockwell-uk/go-logger/logger"
)

type sServer struct {
	listener net.Listener
}

type operation struct {
	op  string
	res chan int
	n   int
}

var (
	connMap = &sync.Map{}
	s       sServer
	ctrl    = make(chan struct{}, 1)
	ops     = make(chan operation)
)

func init() {
	var n int
	go func() {
		for o := range ops {
			n += o.n
			if o.op == "CHECK" {
				o.res <- n
			}
		}
	}()
}

func add(i int) {
	ops <- operation{n: i}
}

func done() {
	ops <- operation{n: -1}
}

func wait() {
	for {
		res := make(chan int)
		ops <- operation{"CHECK", res, 0}
		if <-res == 0 {
			return
		}
	}
}

func Start(host, port string, requestHandler func(connId int, c net.Conn, connMap *sync.Map)) {
	// Start the tcp server
	logger.Log(
		logger.LVL_APP,
		fmt.Sprintf("starting on port %s", port),
	)

	// Store connected clients
	var connId int

	l, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%s", host, port),
	)
	if err != nil {
		logger.Log(
			logger.LVL_ERROR,
			err.Error(),
		)
		return
	}

	logger.Log(
		logger.LVL_APP,
		"started server",
	)
	s = sServer{
		listener: l,
	}

	go func() {
		// Client loop
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-ctrl:
					return
				default:
					logger.Log(
						logger.LVL_ERROR,
						err.Error(),
					)
					break
				}
			}

			connId++
			connMap.Store(connId, conn)
			add(1)
			go func() {
				defer done()
				requestHandler(connId, conn, connMap)
			}()
		}
	}()
}

func Stop() {
	logger.Log(
		logger.LVL_APP,
		"shutting down server",
	)

	// // close connections
	connMap.Range(func(key, value interface{}) bool {
		if conn, ok := value.(net.Conn); ok {
			payload := "quit"
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
			conn.Close()
		}
		return true
	})

	// Signal close
	ctrl <- struct{}{}

	// Close listener
	s.listener.Close()

	// Wait for all connections to finish
	wait()

	// Exit message (Client loop has closed)
	logger.Log(
		logger.LVL_APP,
		"server exited properly",
	)
}

func GracefulExit() {
	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	// Log
	logger.Log(
		logger.LVL_APP,
		"received sigterm",
	)

	// Stop TCP service
	Stop()

	// Finally stop our logger
	logger.Stop()
}
