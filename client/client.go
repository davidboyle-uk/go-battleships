package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	"github.com/dbx123/go-battleships/core"
	"github.com/dbx123/go-battleships/core/types"
	"github.com/dbx123/go-battleships/tcp"
	"github.com/dbx123/go-battleships/util"
)

var (
	conn          net.Conn
	g             *types.Game
	currentPlayer int
)

func Run(host, port string) {
	// connect to the server
	var err error
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// plan for graceful exit
	go GracefulExit()

	// setup
	util.CleanScreen()

	// auto start
	Announce()

	ops := make(chan tcp.Proto)

	go func() {
		reader := bufio.NewReader(conn)
		for {
			m, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
			}
			// Parse
			p, err := tcp.ParseMessage(m)
			if err != nil {
				os.Exit(1)
			}

			ops <- p
		}
	}()

	processRequests(ops)
}

func GracefulExit() {
	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	// Let the server know we quit
	reply(core.LEFT, "")

	// Close connection
	conn.Close()
}

// Announce to the server that we are here.
func Announce() {
	reply(core.HELLO, "")
}

// DoGunShot make a request to gun shot and receive response.
func DoGunShot(c *types.Coordinates) {
	// gun shot
	myPlayer(g).GunShot(otherPlayer(g), c)

	// prepare JSON
	js, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}

	// send to socket
	reply(core.GUNSHOT, string(js))
}

func processRequests(c1 chan tcp.Proto) {
	var wg sync.WaitGroup
	for r := range c1 {
		wg.Add(1)
		go func(r tcp.Proto) {
			handleRequest(r)
			wg.Done()
		}(r)
	}
	wg.Wait()
}

func reply(action, body string) {
	fmt.Fprint(conn, fmt.Sprint(tcp.Proto{Action: action, Player: currentPlayer, Body: body}))
}

func decodeGame(r string) {
	err := json.Unmarshal([]byte(r), &g)
	if err != nil {
		panic(err)
	}
}

func readUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}

func getGameType() string {
	for {
		input, e1 := readUserInput()
		gt, e2 := strconv.Atoi(input)

		v := types.GameType(gt)
		if e1 == nil && e2 == nil && v.IsValid() {
			return input
		}

		fmt.Println(invalidInputText())
	}
}

func getPlayerName() string {
	for {
		input, err := readUserInput()
		if err == nil {
			return input
		}

		fmt.Println(invalidInputText())
	}
}

func getCoordinates(g *types.Game) types.Coordinates {
	for {
		move, e1 := readUserInput()
		if len(move) < 1 {
			fmt.Println("Input cannot be empty")
			continue
		}
		if len(move) < 2 {
			fmt.Println("Inproper input")
			continue
		}

		x := util.RuneToInt(rune(move[0]))
		y, e2 := strconv.Atoi(move[1:])
		c := types.Coordinates{Abscissa: x, Ordinate: y}

		switch {
		case alreadyFired(g, c):
			fmt.Printf("You have already shot coord [%s], try again\n", move)
		case !c.Validate(g.GridSize):
			fmt.Printf("Out of bounds or invalid %s\n", move)
		case e1 == nil && e2 == nil && c.Validate(g.GridSize) && !alreadyFired(g, c):
			return c
		}
	}
}

func awaitText(g *types.Game) {
	fmt.Printf("Wait for %s to fire...\n", opponentName(g))
}

func winnerText(winner string) string {
	return winner + " is the ..." + core.WINNER_TEXT
}

func loserText(loser string) string {
	return loser + " is the ..." + core.LOSER_TEXT
}

func shootText() string {
	return "Enter next move, eg: b10"
}

func invalidInputText() string {
	return "Invalid input, please try again..."
}

func quitText() string {
	return "\nServer shutting down, quitting"
}

func leftText(opponent string) string {
	return fmt.Sprintf("%s left the game", opponent)
}

func alreadyFired(g *types.Game, c types.Coordinates) bool {
	pp, _ := core.CheckShotsFired(&c, myPlayer(g))
	return pp
}

func myPlayer(g *types.Game) *types.Player {
	var p types.Player
	switch *g.GameType {
	case types.ONE_PLAYER:
		return &g.SecondPlayer
	case types.TWO_PLAYER:
		if currentPlayer == 1 {
			return &g.FirstPlayer
		}
		return &g.SecondPlayer
	}
	return &p
}

func otherPlayer(g *types.Game) *types.Player {
	var p types.Player
	switch *g.GameType {
	case types.ONE_PLAYER:
		return &g.FirstPlayer
	case types.TWO_PLAYER:
		if currentPlayer == 1 {
			return &g.SecondPlayer
		}
		return &g.FirstPlayer
	}
	return &p
}

func myName(g *types.Game) string {
	return myPlayer(g).Name
}

func opponentName(g *types.Game) string {
	return otherPlayer(g).Name
}
