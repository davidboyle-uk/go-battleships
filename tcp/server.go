package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"

	"go-battleships/logger"
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
	s    sServer
	ctrl = make(chan struct{}, 1)
	ops  = make(chan operation)
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

func Start(host, port string, requestHandler func(ct context.Context, connId int, c net.Conn, connMap *sync.Map)) {
	// Start the tcp server
	logger.Log(
		logger.LVL_INFO,
		fmt.Sprintf("starting on port %s", port),
	)

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
		logger.LVL_INFO,
		"started server",
	)
	s = sServer{
		listener: l,
	}

	var connId int
	var connMap = &sync.Map{}

	// Get the context and a function to cancel it
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Client loop
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-ctrl:
					cancel()
					logger.Log(
						logger.LVL_INFO,
						err.Error(),
					)
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
			go func(ct context.Context) {
				defer done()
				requestHandler(ct, connId, conn, connMap)
			}(ctx)
		}
	}()
}

func Stop() {
	logger.Log(
		logger.LVL_INFO,
		"shutting down server",
	)

	// Signal close
	ctrl <- struct{}{}

	// Close listener
	s.listener.Close()

	// HACK
	os.Exit(0)

	// Wait for all connections to finish
	wait()

	// Exit message (Client loop has closed)
	logger.Log(
		logger.LVL_INFO,
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
		logger.LVL_INFO,
		"recieved sigterm",
	)

	// Stop TCP service
	Stop()

	// Finally stop our logger
	logger.Stop()
}
