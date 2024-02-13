package types

import (
	"encoding/json"
	"strconv"

	bb_types "github.com/dbx123/battleships-board/types"
	"github.com/dbx123/go-battleships/palette"
	"github.com/dbx123/go-battleships/util"
)

var (
	GAME_GRID_BORDER = "|"
	STR_SHIP_OK      = palette.White("██")
	STR_SHIP_HIT     = palette.Red("██")
	STR_SHIP_SUNK    = palette.Yellow("██")
	STR_SEA_HIT      = palette.Cyan("~~")
	STR_SEA_OK       = "  "
	STR_STATUS_ERROR = "??"

	PRINT_PLAYER_BOARD   = 2
	PRINT_OPPONENT_BOARD = 1
)

var legend = map[int]string{
	0: "Legend:",
	1: "Ship (OK)        " + STR_SHIP_OK,
	2: "Ship (HIT)       " + STR_SHIP_HIT,
	3: "Ship (SUNK)      " + STR_SHIP_SUNK,
	4: "Sea (MISS)       " + STR_SEA_HIT,
}

type Game struct {
	Game      bb_types.Game
	GameType  *GameType
	LastMoves map[int]bb_types.Coord
}

type GameType int

const (
	_ GameType = iota
	ONE_PLAYER
	TWO_PLAYER
)

func (t GameType) IsValid() bool {
	return t >= ONE_PLAYER && t <= TWO_PLAYER
}

func (g *Game) ToJSON() string {
	json, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func (g *Game) LastShot(playerID int, playerName string) string {
	if lastMove, ok := g.LastMoves[playerID]; ok {
		return PrettyPrintCoord(lastMove) + " shot fired by " + playerName + "\n"
	}
	return ""
}

func (g *Game) PrettyPrintGame(currentPlayer int) string {
	// clean tty screen
	util.CleanScreen()

	pOneID := 0
	pTwoID := 1

	pOne := g.Game.Players[pOneID]
	pTwo := g.Game.Players[pTwoID]

	p1 := pOne.Name
	p2 := pTwo.Name

	dim := pOne.Board.Dim

	var gs string
	switch currentPlayer {
	case 1:
		gs += ">>> " + p1 + "'s Board <<<\n\n"
		gs += PrettyPrintBoard(pOne.Board.Moves, pTwo.Moves, dim, PRINT_PLAYER_BOARD)
		gs += "\n\n>>> " + p2 + "'s Board <<<\n\n"
		gs += PrettyPrintBoard(bb_types.Moves{}, pOne.Moves, dim, PRINT_OPPONENT_BOARD)
	case 2:
		gs += ">>> " + p2 + "'s Board <<<\n\n"
		gs += PrettyPrintBoard(pTwo.Board.Moves, pOne.Moves, dim, PRINT_PLAYER_BOARD)
		gs += "\n\n>>> " + p1 + "'s Board <<<\n\n"
		gs += PrettyPrintBoard(bb_types.Moves{}, pTwo.Moves, dim, PRINT_OPPONENT_BOARD)
	}

	gs += "\n"
	switch *g.GameType {
	case ONE_PLAYER:
		gs += g.LastShot(pTwoID, p2)
		gs += g.LastShot(pOneID, p1)
	case TWO_PLAYER:
		gs += g.LastShot(pOneID, p1)
		gs += g.LastShot(pTwoID, p2)
	}

	return gs
}

func PrettyPrintCoord(c bb_types.Coord) string {
	return "[" + util.IntToLetter(c.X) + strconv.Itoa(c.Y+1) + "]"
}

func PrettyPrintBoard(myBoardMoves bb_types.Moves, playerMoves bb_types.Moves, dim, h int) string {
	var s string

	// create column indicator line
	s = "   " + GAME_GRID_BORDER
	for r := 0; r < dim-1; r++ {
		s += "-" + util.IntToLetter(r) + "--" + GAME_GRID_BORDER
	}
	s += "-" + util.IntToLetter(dim-1) + "--" + GAME_GRID_BORDER + "\n"

	// create first separation line
	s += "   " + GAME_GRID_BORDER
	for r := 0; r < dim-1; r++ {
		s += "-----"
	}
	s += "----" + GAME_GRID_BORDER + "\n"

	// for each col
	var state string
	for y := 0; y < dim; y++ {
		// start with legend
		pad := "  "
		if y >= 9 {
			pad = " "
		}
		s += strconv.Itoa(y+1) + pad

		// add grid border
		s += GAME_GRID_BORDER

		// for each column
		for x := 0; x < dim; x++ {
			t := bb_types.Coord{X: x, Y: y}.String()

			state = STR_SEA_OK

			// if we are drawing caller's Board
			switch h {
			case PRINT_PLAYER_BOARD:
				if move, ok := myBoardMoves[t]; ok {
					state = ConvertState(move.State)
					if move.Ship != nil {
						state = STR_SHIP_OK
					}
				}
				if move, ok := playerMoves[t]; ok {
					state = ConvertState(move.State)
				}

			// if we are drawing opponent's Board
			case PRINT_OPPONENT_BOARD:
				if move, ok := playerMoves[t]; ok {
					state = ConvertState(move.State)
				}
			}

			s += " " + state + " " + GAME_GRID_BORDER
		}

		// create separation line
		s += "\n" + "   " + GAME_GRID_BORDER
		l := "" // legend text
		for c := 0; c < dim-1; c++ {
			s += "-----"
			if lt, ok := legend[y]; ok {
				l = lt
			}
		}
		s += "----" + GAME_GRID_BORDER + "   " + l + "\n"
	}

	return s
}

func ConvertState(s bb_types.State) string {
	switch s {
	case bb_types.SEA:
		return STR_SEA_OK
	case bb_types.HIT:
		return STR_SHIP_HIT
	case bb_types.MISS:
		return STR_SEA_HIT
	case bb_types.SUNK:
		return STR_SHIP_SUNK
	default:
		return STR_STATUS_ERROR
	}
}
