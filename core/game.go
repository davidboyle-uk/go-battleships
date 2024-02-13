package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dbx123/battleships-board/pkg/ai"
	"github.com/dbx123/battleships-board/pkg/game"
	bb_types "github.com/dbx123/battleships-board/types"
	"github.com/dbx123/go-battleships/tcp"
	"github.com/dbx123/go-battleships/types"

	"github.com/rockwell-uk/go-logger/logger"
)

func ProcessRequest(p tcp.Proto) ([]tcp.Proto, error) {
	logger.Log(
		logger.LVL_INTERNAL,
		fmt.Sprintf("process request: %s", p),
	)
	// Run the operation
	switch p.Action {
	// Initial connection
	case HELLO:
		return Hello()
	case GAMETYPE:
		return SetGameType(p)
	case ASSIGN:
		return Assigned(p)
	case PLAYER_NAME:
		return SetPlayerName(p)
	case GUNSHOT:
		return DoGunShot(p)
	case LEFT:
		return DoLeft(p)
	case QUIT:
		return DoQuit()
	default:
		logger.Log(
			logger.LVL_ERROR,
			fmt.Sprintf("unknown action %s", p.Action),
		)
		return []tcp.Proto{}, fmt.Errorf("unknown action %s", p.Action)
	}
}

func DoLeft(p tcp.Proto) ([]tcp.Proto, error) {
	return []tcp.Proto{
		{
			Action: LEFT,
			Player: 0,
		},
		{
			Action: QUIT,
			Player: 0,
		},
	}, nil
}

func DoQuit() ([]tcp.Proto, error) {
	return []tcp.Proto{
		{
			Action: QUIT,
			Player: 0,
		},
	}, nil
}

func Hello() ([]tcp.Proto, error) {
	// Has the game type been set?
	if g.GameType == nil {
		return []tcp.Proto{
			{
				Action: GAMETYPE,
				Body:   "Enter game type - 1 = vs CPU, 2 = vs Human: ",
			},
		}, nil
	}

	return []tcp.Proto{
		{
			Action: ASSIGN,
			Player: 2,
			Body:   "2",
		},
	}, nil
}

func SetGameType(p tcp.Proto) ([]tcp.Proto, error) {
	logger.Log(
		logger.LVL_DEBUG,
		fmt.Sprintf("setting game type - %v", p.Body),
	)
	n, _ := strconv.Atoi(p.Body)
	gt := types.GameType(n)
	g.GameType = &gt

	var body string
	switch *g.GameType {
	case types.ONE_PLAYER:
		body = "2"
	case types.TWO_PLAYER:
		body = "1"
	}

	return []tcp.Proto{
		{
			Action: ASSIGN,
			Player: 1,
			Body:   body,
		},
	}, nil
}

func Assigned(p tcp.Proto) ([]tcp.Proto, error) {
	var player int
	switch *g.GameType {
	case types.ONE_PLAYER:
		player = 0
	case types.TWO_PLAYER:
		player = p.Player
	}

	return []tcp.Proto{
		{
			Action: PLAYER_NAME,
			Player: player,
			Body:   "Enter player name: ",
		},
	}, nil
}

func SetPlayerName(p tcp.Proto) ([]tcp.Proto, error) {
	logger.Log(
		logger.LVL_DEBUG,
		fmt.Sprintf("setting p%v name - %s", p.Player, p.Body),
	)
	switch *g.GameType {
	case types.ONE_PLAYER:
		g.Game.Players[0].Name = "CPU"
		g.Game.Players[1].Name = p.Body
		return []tcp.Proto{
			{
				Action: DRAW_GAME_SHOOT,
				Player: 1,
				Body:   getGame(1),
			},
		}, nil
	case types.TWO_PLAYER:
		switch p.Player {
		case 1:
			g.Game.Players[0].Name = p.Body
			return []tcp.Proto{
				{
					Action: AWAIT_OPPONENT,
					Player: 1,
					Body:   "Waiting for opponent to join...",
				},
			}, nil
		case 2:
			g.Game.Players[1].Name = p.Body
			return []tcp.Proto{
				{
					Action: DRAW_GAME_SHOOT,
					Player: 1,
					Body:   getGame(1),
				},
				{
					Action: DRAW_GAME_AWAIT,
					Player: 2,
					Body:   getGame(2),
				},
			}, nil
		}
	}

	return []tcp.Proto{}, errors.New("already have 2 players")
}

func CheckWinner() (int, int) {
	if g.Game.Players[0].Hits == g.Game.Players[0].Board.ShipTot {
		return 1, 2
	}
	if g.Game.Players[1].Hits == g.Game.Players[1].Board.ShipTot {
		return 2, 1
	}
	return 0, 0
}

func CombineMoves(a, b bb_types.Moves) bb_types.Moves {
	new := make(bb_types.Moves)
	for k, v := range a {
		new[k] = v
	}
	for k, v := range b {
		new[k] = v
	}
	return new
}

func CoordFromString(s string) (bb_types.Coord, error) {
	bits := strings.Split(s, " ")
	x, err := strconv.Atoi(bits[0])
	if err != nil {
		return bb_types.Coord{}, err
	}
	y, err := strconv.Atoi(bits[1])
	if err != nil {
		return bb_types.Coord{}, err
	}

	return bb_types.Coord{
		X: x,
		Y: y,
	}, nil
}

func DoGunShot(p tcp.Proto) ([]tcp.Proto, error) {
	// get bearings
	var myID int
	var us, them *bb_types.Player
	switch p.Player {
	case 1:
		myID = 0
		us = g.Game.Players[0]
		them = g.Game.Players[1]
	case 2:
		myID = 1
		us = g.Game.Players[1]
		them = g.Game.Players[0]
	}

	// decode coord
	c, err := CoordFromString(p.Body)
	if err != nil {
		return []tcp.Proto{}, err
	}

	// gun shot
	game.TakeShot(us, them, c)
	g.LastMoves[myID] = c

	switch *g.GameType {
	case types.ONE_PLAYER:
		winner, _ := CheckWinner()
		if winner == 2 {
			return []tcp.Proto{
				{
					Action: DRAW_ENDSCREEN,
					Player: 0,
					Body:   WINNER,
				},
			}, nil
		}

		move := ai.CalculateMove(bb_types.Board{
			Dim:   g.Game.Players[1].Board.Dim,
			Moves: g.Game.Players[0].Moves,
		})
		logger.Log(
			logger.LVL_INTERNAL,
			fmt.Sprintf("CPU shot %s", move),
		)
		// ai gun shot
		game.TakeShot(g.Game.Players[0], g.Game.Players[1], move)
		g.LastMoves[0] = move

		winner, _ = CheckWinner()
		if winner == 1 {
			return []tcp.Proto{
				{
					Action: DRAW_ENDSCREEN,
					Player: 0,
					Body:   LOSER,
				},
			}, nil
		}

		// send back to client
		return []tcp.Proto{
			{
				Action: DRAW_GAME_SHOOT,
				Player: 0,
				Body:   getGame(1),
			},
		}, nil

	case types.TWO_PLAYER:
		var opponent int
		switch p.Player {
		case 1:
			opponent = 2
		case 2:
			opponent = 1
		}

		winner, loser := CheckWinner()
		if winner > 0 {
			return []tcp.Proto{
				{
					Action: DRAW_ENDSCREEN,
					Player: winner,
					Body:   WINNER,
				},
				{
					Action: DRAW_ENDSCREEN,
					Player: loser,
					Body:   LOSER,
				},
			}, nil
		}

		return []tcp.Proto{
			{
				Action: DRAW_GAME_AWAIT,
				Player: p.Player,
				Body:   getGame(p.Player),
			},
			{
				Action: DRAW_GAME_SHOOT,
				Player: opponent,
				Body:   getGame(p.Player),
			},
		}, nil
	}

	return []tcp.Proto{}, nil
}

func getGame(playerID int) string {
	p1 := *g.Game.Players[0]
	p2 := *g.Game.Players[1]
	gc := types.Game{
		GameType: g.GameType,
		Game: bb_types.Game{
			Players: []*bb_types.Player{
				&p1,
				&p2,
			},
		},
		LastMoves: g.LastMoves,
	}

	switch playerID {
	case 1:
		switch *g.GameType {
		case types.ONE_PLAYER:
			p1.Board.Moves = bb_types.Moves{}
		case types.TWO_PLAYER:
			p2.Board.Moves = bb_types.Moves{}
		}
	case 2:
		switch *g.GameType {
		case types.TWO_PLAYER:
			p1.Board.Moves = bb_types.Moves{}
		}
	}

	gc.Game.Players[0] = &p1
	gc.Game.Players[1] = &p2

	res, err := json.Marshal(gc)
	if err != nil {
		panic(err)
	}

	return string(res)
}
