package main

import (
	"flag"

	"go-battleships/core"
	"go-battleships/logger"
	"go-battleships/tcp"
	"go-battleships/util"
)

// ###########################################################################################################
// ############################################# SERVER LOGIC ################################################
// ###########################################################################################################

var (
	serverHost string = "localhost"
	serverPort string = "3333"
	numShips   int    = 5
)

func main() {
	// Start fresh
	util.CleanScreen()

	// *** Define and parse flags ***

	// Server port
	flag.StringVar(&serverPort, "port", serverPort, "port for our server")

	// Server port
	flag.StringVar(&serverHost, "host", serverHost, "ip or host for our server")

	// Number of ships
	flag.IntVar(&numShips, "ships", numShips, "number of ships")

	// Log level
	var v, vv, vvv bool
	flag.IntVar(&logger.Vbs, "verbosity", logger.Vbs, "logging level, default ERROR")
	flag.BoolVar(&v, "v", false, "WARN log verbosity override")
	flag.BoolVar(&vv, "vv", false, "INFO log verbosity override")
	flag.BoolVar(&vvv, "vvv", false, "DEBUG log verbosity override")

	flag.Parse()

	// Start our logger
	var vbs int
	switch {
	case vvv:
		vbs = logger.LVL_DEBUG
	case vv:
		vbs = logger.LVL_INFO
	case v:
		vbs = logger.LVL_WARN
	}
	logger.Start(vbs)

	core.SetGame(core.PrepareGame(core.CPU_GRID, numShips))

	// Start tcp
	tcp.Start(serverHost, serverPort, core.HandleRequest)

	// Plan for a graceful exit
	tcp.GracefulExit()
}
