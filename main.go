package main

import (
	"flag"

	"github.com/dbx123/go-battleships/client"
	"github.com/dbx123/go-battleships/server"
	"github.com/dbx123/go-battleships/util"

	"github.com/rockwell-uk/go-logger/logger"
)

type runType string

const (
	RUNTYPE_SERVER runType = "server"
	RUNTYPE_CLIENT runType = "client"
)

func (r runType) Validate() bool {
	if r == RUNTYPE_SERVER {
		return true
	}
	if r == RUNTYPE_CLIENT {
		return true
	}
	return false
}

var (
	clientOrServer string = "client"
	host           string = "localhost"
	port           string = "3333"
	numShips       int    = 5
)

func main() {
	// Start fresh
	util.CleanScreen()

	// *** Define and parse flags ***

	// Client / Server
	flag.StringVar(&clientOrServer, "type", clientOrServer, "run type client/server")

	// Server port
	flag.StringVar(&port, "port", port, "port for our server")

	// Server port
	flag.StringVar(&host, "host", host, "ip or host for our server")

	// Number of ships
	flag.IntVar(&numShips, "ships", numShips, "number of ships")

	// Log level
	var v, vv bool
	flag.BoolVar(&v, "v", false, "DEBUG level log verbosity override")
	flag.BoolVar(&vv, "vv", false, "INTERNAL level log verbosity override")

	flag.Parse()

	// Start our logger
	var vbs logger.LogLvl = logger.LVL_APP
	switch {
	case vv:
		vbs = logger.LVL_INTERNAL
	case v:
		vbs = logger.LVL_DEBUG
	}
	logger.Start(vbs)

	clientOrServer := runType(clientOrServer)
	if !clientOrServer.Validate() {
		logger.Log(
			logger.LVL_FATAL,
			"Type must be 'client' or 'server'",
		)
		return
	}

	switch clientOrServer {
	case RUNTYPE_SERVER:
		server.Run(host, port, numShips)
	case RUNTYPE_CLIENT:
		client.Run(host, port)
	}
}
