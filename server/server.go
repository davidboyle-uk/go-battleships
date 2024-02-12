package server

import (
	"github.com/dbx123/go-battleships/core"
	"github.com/dbx123/go-battleships/tcp"
)

func Run(host, port string) {
	// Create game
	core.PrepareGame(core.CPU_GRID)

	// Start tcp
	tcp.Start(host, port, core.HandleRequest)

	// Plan for a graceful exit
	tcp.GracefulExit()
}
