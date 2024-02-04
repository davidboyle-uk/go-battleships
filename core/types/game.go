package types

import (
	"encoding/json"
	"go-battleships/util"
)

type Game struct {
	GridSize     int       `json:"gs"`
	ShipVolume   int       `json:"sv"`
	GameType     *GameType `json:"gt"`
	FirstPlayer  Player    `json:"p1"`
	SecondPlayer Player    `json:"p2"`
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

func (g *Game) LastShotInToString() string {
	var player Player
	switch *g.GameType {
	case ONE_PLAYER:
		player = g.FirstPlayer
	case TWO_PLAYER:
		player = g.SecondPlayer
	}

	if len(player.ShotsFired) > 0 {
		c := player.ShotsFired[len(player.ShotsFired)-1]
		return c.PrettyPrintCoordinatesInfo() + " shot fired by " + player.Name + "\n"
	}
	return ""
}

func (g *Game) LastShotOutToString() string {
	var player Player
	switch *g.GameType {
	case ONE_PLAYER:
		player = g.SecondPlayer
	case TWO_PLAYER:
		player = g.FirstPlayer
	}

	if len(player.ShotsFired) > 0 {
		c := player.ShotsFired[len(player.ShotsFired)-1]
		return c.PrettyPrintCoordinatesInfo() + " shot fired by " + player.Name + "\n"
	}
	return ""
}

// PrettyPrintGame from p Player to t Player in p Coordinates
//
//	[g:*Game]			Game pointer		[m:int]	game mode 0 1:PC 1:1
func (g *Game) PrettyPrintGame(currentPlayer int) string {
	// clean tty screen
	util.CleanScreen()

	var p1, p2 string
	switch *g.GameType {
	case ONE_PLAYER:
		p1 = g.SecondPlayer.Name
		p2 = g.FirstPlayer.Name
	case TWO_PLAYER:
		p1 = g.FirstPlayer.Name
		p2 = g.SecondPlayer.Name
	}

	var gs string
	switch currentPlayer {
	case 1:
		gs += ">>> " + p1 + "'s Board <<<\n\n"
		gs += g.FirstPlayer.SeaToString(PRINT_CALLERSEA_MODE)
		gs += "\n\n>>> " + p2 + "'s Board <<<\n\n"
		gs += g.SecondPlayer.SeaToString(PRINT_OPPONENTSEA_MODE)
	case 2:
		gs += ">>> " + p2 + "'s Board <<<\n\n"
		gs += g.SecondPlayer.SeaToString(PRINT_CALLERSEA_MODE)
		gs += "\n\n>>> " + p1 + "'s Board <<<\n\n"
		gs += g.FirstPlayer.SeaToString(PRINT_OPPONENTSEA_MODE)
	}

	gs += "\n"
	gs += g.LastShotOutToString()
	gs += g.LastShotInToString()

	return gs
}

func (g *Game) DebugGame() string {
	jg, _ := json.MarshalIndent(g, "", "    ")
	return string(jg)
}
