package client

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dbx123/go-battleships/core"
	"github.com/dbx123/go-battleships/tcp"
	"github.com/dbx123/go-battleships/util"
)

func handleRequest(p tcp.Proto) tcp.Proto {
	switch p.Action {
	case core.GAMETYPE:
		fmt.Print(p.Body)
		reply(core.GAMETYPE, getGameType())

	case core.ASSIGN:
		i, err := strconv.Atoi(p.Body)
		if err != nil {
			_, _ = conn.Write([]byte(err.Error()))
			break
		}
		currentPlayer = i
		reply(core.ASSIGN, p.Body)

	case core.PLAYER_NAME:
		fmt.Print(p.Body)
		reply(core.PLAYER_NAME, getPlayerName())

	case core.AWAIT_OPPONENT:
		fmt.Print(p.Body)

	case core.DRAW_GAME_AWAIT:
		decodeGame(p.Body)
		fmt.Println(g.PrettyPrintGame(currentPlayer))
		awaitText(g)

	case core.DRAW_GAME_SHOOT:
		decodeGame(p.Body)
		fmt.Println(g.PrettyPrintGame(currentPlayer))
		fmt.Println(shootText())
		c := getCoordinates(g)
		DoGunShot(c)

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
		os.Exit(0)

	case core.LEFT:
		fmt.Print(leftText(opponentName(g)))
		reply(core.QUIT, "")

	case core.QUIT:
		fmt.Println(quitText())
		os.Exit(0)
	}

	return tcp.Proto{}
}
