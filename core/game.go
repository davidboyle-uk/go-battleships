package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"go-battleships/core/types"
	"go-battleships/logger"
	"go-battleships/tcp"
	"go-battleships/util"
)

func ProcessRequest(p tcp.Proto) ([]tcp.Proto, error) {
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
	case QUIT:
		os.Exit(1)
		return Quit()
	default:
		logger.Log(
			logger.LVL_ERROR,
			fmt.Sprintf("unknown action %s", p.Action),
		)
		return []tcp.Proto{}, fmt.Errorf("unknown action %s", p.Action)
	}
}

func Quit() ([]tcp.Proto, error) {
	return []tcp.Proto{{
		Action: QUIT,
	}}, nil
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
		g.FirstPlayer.Name = "CPU"
		g.SecondPlayer.Name = p.Body
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
			g.FirstPlayer.Name = p.Body
			return []tcp.Proto{
				{
					Action: AWAIT_OPPONENT,
					Player: 1,
					Body:   "Waiting for opponent to join...",
				},
			}, nil
		case 2:
			g.SecondPlayer.Name = p.Body
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

func SleepRequest() {
	time.Sleep(SIMULATION_THINKING_TIME * time.Millisecond)
}

// DoGunShot from p Player to t Player in p Coordinates
func DoGunShot(p tcp.Proto) ([]tcp.Proto, error) {
	// decode received game
	GameDecoder(p.Body)

	switch *g.GameType {
	case types.ONE_PLAYER:

		winner, _ := CheckWinner()
		if winner == 2 {
			return []tcp.Proto{
				{
					Action: DRAW_ENDSCREEN,
					Player: 2,
					Body:   WINNER,
				},
			}, nil
		}

		// debug print of json
		//fmt.Println(DebugGame(&g))

		// debug pause
		//util.ConsolePause(util.PAUSE_MEX)
		Timeout()

		// TODO: IMPLEMENT STRATEGY
		s := util.Random(0, len(g.SecondPlayer.Sea.Ships)-1)
		k := util.Random(0, len(g.SecondPlayer.Sea.Ships[s].Positions)-1)
		g.FirstPlayer.GunShot(&g.SecondPlayer, &g.SecondPlayer.Sea.Ships[s].Positions[k])

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

		// debug pause
		//util.ConsolePause(util.PAUSE_MEX)
		Timeout()

		// print game
		//fmt.Println(PrettyPrintGame(&g))

		// send back to client
		return []tcp.Proto{
			{
				Action: DRAW_GAME_SHOOT,
				Player: 2,
				Body:   getGame(2),
			},
		}, nil
	case types.TWO_PLAYER:

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

		var opponent int
		switch p.Player {
		case 1:
			opponent = 2
		case 2:
			opponent = 1
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
				Body:   getGame(opponent),
			},
		}, nil
	}

	return []tcp.Proto{}, nil
}

func getGame(p int) string {
	res, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}
	return string(res)
}
