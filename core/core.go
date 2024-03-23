package core

import (
	"fmt"

	"github.com/davidboyle-uk/battleships-board/pkg/game"
	bb_types "github.com/davidboyle-uk/battleships-board/types"
	"github.com/davidboyle-uk/go-battleships/types"
	"github.com/rockwell-uk/go-logger/logger"
)

var (
	g types.Game
)

// Constants for default and game config.
const (
	HELLO           = "hello"
	GAMETYPE        = "gametype"
	ASSIGN          = "assign"
	PLAYER_NAME     = "player"
	AWAIT_OPPONENT  = "awaitopponent"
	DRAW_GAME_AWAIT = "takecover"
	DRAW_GAME_SHOOT = "shoot"
	GUNSHOT         = "gunshot"
	EXIT            = "exit"
	DRAW_ENDSCREEN  = "gameover"
	WINNER          = "winner"
	LOSER           = "loser"
	LEFT            = "leftgame"
	QUIT            = "quit"

	CPU_NAME = "CPU"
	CPU_GRID = 10
)

func PrepareGame(boardSize int) {
	g = types.Game{
		Game:      game.Initialise(boardSize),
		LastMoves: make(map[int]bb_types.Coord),
	}
	logger.Log(
		logger.LVL_INTERNAL,
		fmt.Sprintf("Initialised game %v", g.ToJSON()),
	)
}
