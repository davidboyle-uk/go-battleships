package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"go-battleships/core"
	"go-battleships/core/types"
	"go-battleships/tcp"
	"go-battleships/util"
)

var (
	clientPort string = "3333"
	clientHost string = "127.0.0.1"
)

var (
	conn          net.Conn
	currentPlayer int
)

const (
	INVALID_INPUT = "Invalid input, please try again..."
)

// Announce to the server that we are here
func Announce(g *types.Game) {
	reply(core.HELLO, "")
}

// DoGunShot make a request to gun shot and receive response
func DoGunShot(g *types.Game, c *types.Coordinates) {
	// gun shot
	myPlayer(g).GunShot(otherPlayer(g), c)

	// prepare JSON
	js, _ := json.Marshal(g)

	// send to socket
	reply(core.GUNSHOT, string(js))
}

// DoGunShot make an exit request
func ExitGame() {
	// make exit request
	reply(core.EXIT, "")

	// clean screen
	util.CleanScreen()

	// exit
	os.Exit(1)
}

func getCoordinates(g *types.Game) types.Coordinates {
	var move string
	var c types.Coordinates

	for {
		_, e1 := fmt.Scanf("%s", &move)
		if len(move) < 1 {
			fmt.Println("Input cannot be empty")
		}
		if len(move) < 2 {
			fmt.Println("Inproper input")
		}

		x := util.RuneToInt(rune(move[0]))
		y, e2 := strconv.Atoi(move[1:])
		c = types.Coordinates{Abscissa: x, Ordinate: y}

		switch {
		case alreadyFired(g, c):
			fmt.Printf("You have already shot coord %s, try again\n", move)
		case !c.Validate(g.GridSize):
			fmt.Printf("Out of bounds or invalid %s\n", move)
		case e1 == nil && e2 == nil && c.Validate(g.GridSize) && !alreadyFired(g, c):
			return c
		}
	}
}

func getPlayerName() string {
	var input string

	for {
		_, err := fmt.Scanf("%s", &input)
		if err == nil {
			return input
		}

		fmt.Println(INVALID_INPUT)
	}
}

func getGameType() string {
	var input string

	for {
		_, e1 := fmt.Scanf("%s", &input)
		gt, e2 := strconv.Atoi(input)
		v := types.GameType(gt)
		if e1 == nil && e2 == nil && v.IsValid() {
			return input
		}

		fmt.Println(INVALID_INPUT)
	}
}

func winnerText(winner string) string {
	return winner + " is the ..." + core.WINNER_TEXT
}

func loserText(loser string) string {
	return loser + " is the ..." + core.LOSER_TEXT
}

func shoot() {
	fmt.Println("Enter next move, eg: b10")
}

func getShotsFired(g *types.Game) []types.Coordinates {
	return myPlayer(g).ShotsFired
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

func await(g *types.Game) {
	fmt.Printf("Wait for %s to fire...\n", opponentName(g))
}

func reply(action, body string) {
	fmt.Fprintf(conn, fmt.Sprint(tcp.Proto{Action: action, Player: currentPlayer, Body: body}))
}

func main() {
	// *** Define and parse flags ***

	// Server host
	flag.StringVar(&clientHost, "host", clientHost, "host for our server")

	// Server port
	flag.StringVar(&clientPort, "port", clientPort, "port for our server")

	flag.Parse()

	// connect to the server
	var err error
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", clientHost, clientPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// setup
	util.CleanScreen()
	g := types.Game{}

	// auto start
	Announce(&g)

	notify := make(chan error)
	// run it
	go func() {
		play(notify, &g)
	}()

	for {
		select {
		case err := <-notify:
			if err == io.EOF {
				fmt.Println("connection to server was closed")
				return
			}
			break
		}
	}
}

func play(ch chan error, g *types.Game) {
	for {
		// listen for reply
		m, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			ch <- err
			if io.EOF == err {
				close(ch)
				return
			}
		}

		// Parse
		p, err := tcp.ParseMessage(m)
		if err != nil {
			os.Exit(1)
		}

		// Process the Action
		switch p.Action {
		case core.GAMETYPE:
			fmt.Print(p.Body)
			reply(core.GAMETYPE, getGameType())
		case core.ASSIGN:
			i, err := strconv.Atoi(p.Body)
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			currentPlayer = i
			reply(core.ASSIGN, p.Body)
		case core.PLAYER_NAME:
			fmt.Print(p.Body)
			reply(core.PLAYER_NAME, getPlayerName())
		case core.AWAIT_OPPONENT:
			fmt.Print(p.Body)
		case core.DRAW_GAME_AWAIT:
			json.Unmarshal([]byte(p.Body), &g)
			fmt.Println(g.PrettyPrintGame(currentPlayer))
			await(g)
		case core.DRAW_GAME_SHOOT:
			json.Unmarshal([]byte(p.Body), &g)
			fmt.Println(g.PrettyPrintGame(currentPlayer))
			shoot()
			// fire!!!
			c := getCoordinates(g)
			DoGunShot(g, &c)
		case core.DRAW_ENDSCREEN:
			util.CleanScreen()
			switch p.Body {
			case core.WINNER:
				fmt.Print(winnerText(myName(g)))
				fmt.Print(loserText(opponentName(g)))
			case core.LOSER:
				fmt.Print(loserText(myName(g)))
				fmt.Print(winnerText(opponentName(g)))
			}
			reply(core.QUIT, "")
		case core.QUIT:
			fmt.Printf("server shutting down, quitting")
			os.Exit(0)
		}
	}
}
