package core

import (
	"encoding/json"
	"fmt"
	"time"

	"go-battleships/core/types"
	"go-battleships/util"
)

var (
	g types.Game
)

// Constants for default and game config
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
	QUIT            = "quit"

	OK    = "OK"
	ERROR = "ERR"

	SIMULATION_THINKING_TIME = 1 //milliseconds

	CPU_NAME = "CPU"
	CPU_GRID = 10
)

func SetGame(game types.Game) {
	g = game
}

// Game init
//
// [d:int]	grid dimension
// [sf:int] number of ships per player
func PrepareGame(d int, sf int) types.Game {
	return types.Game{
		GridSize:     d,
		ShipVolume:   CalcShipVolume(sf),
		FirstPlayer:  types.Player{Sea: PrepareSea(d, sf)},
		SecondPlayer: types.Player{Sea: PrepareSea(d, sf)},
	}
}

func CalcShipVolume(n int) int {
	return n * (n + 1) / 2
}

func CheckWinner() (int, int) {
	if len(g.FirstPlayer.Hits) == g.ShipVolume {
		return 1, 2
	}
	if len(g.SecondPlayer.Hits) == g.ShipVolume {
		return 2, 1
	}
	return 0, 0
}

// Sea init
//
// [d:int] grid dimension
// [s:int] number of ship
//
// TODO CREATE MORE EFFICIENT ALGORITHM FOR RANDOM GEN OF SHIPS
func PrepareSea(d int, n int) (s types.Sea) {
	// prepare array of Ship
	ss := make([]types.Ship, n)

	// create n Ship with incremental dimension
	for i := 0; i < n; i++ {
		// create Ship
		st := PrepareShip(i+1, d)

		// if it doesn't collide with other ships
		if !CheckCollisions(&st, ss) {
			// add to Sea
			ss[i] = st
		} else {
			// retry
			i--
		}
	}

	// create Sea
	s = types.Sea{Dimension: d, Ships: ss}
	return
}

// Ship init
//
// [sd:int] ship dimension
// [gd:int]	grid dimension
func PrepareShip(sd int, gd int) types.Ship {
	// choose if horizontal
	h := util.Random(0, 1) == 1

	// create Ship coordinates
	p := make([]types.Coordinates, sd)

	// if Ship dimension is 1
	if sd == 1 {
		// create Random coordinate
		x := util.Random(1, gd)
		y := util.Random(1, gd)
		// add unique Coordinate
		p[0] = types.Coordinates{Abscissa: x, Ordinate: y, Status: types.STATUS_SHIP_OK}

		return types.Ship{Dimension: sd, Positions: p}
	}

	// create x coordinate no more than grid dimension
	x := util.Random(1, gd-sd)
	// create y coordinate no more than grid dimension
	y := util.Random(1, gd)

	// create Coordinates
	for t := 0; t < sd; t++ {
		// offset on x
		if h {
			p[t] = types.Coordinates{Abscissa: x + t, Ordinate: y, Status: types.STATUS_SHIP_OK}
			// offset on y
		} else {
			p[t] = types.Coordinates{Abscissa: y, Ordinate: x + t, Status: types.STATUS_SHIP_OK}
		}
	}

	return types.Ship{Dimension: sd, Positions: p}
}

func Timeout() {
	time.Sleep(time.Second * SIMULATION_THINKING_TIME)
}

// GameDecoder decode game in request
func GameDecoder(r string) {
	// decode game
	err := json.Unmarshal([]byte(r), &g)
	if err != nil {
		panic(err)
	}
}

// CheckCollisions check if a collides with at least one of b ships
//
//	[a:*Ship]	ship pointer		[b:array of Ships]		array of Ships
func CheckCollisions(a *types.Ship, b []types.Ship) bool {
	for _, sb := range b {
		if CheckCollision(a, &sb) {
			return true
		}
	}
	return false
}

// CheckCollisions check if a collides with b
//
//	[a:*Ship]	a ship pointer		[b:*Ship]		b Ship pointer
func CheckCollision(a *types.Ship, b *types.Ship) bool {
	for _, av := range a.Positions {
		for _, bv := range b.Positions {
			if av.Abscissa == bv.Abscissa && av.Ordinate == bv.Ordinate {
				return true
			}
		}
	}
	return false
}

// CheckShotsFired check p coordinates in given Sea's Player
//
//	[p:*Coordinates]	Coordinate point pointer		[pp:*Player]		b Player pointer
//	[return]	bool (collision), ship index, coordinate index
func CheckShotsFired(p *types.Coordinates, pp *types.Player) (bool, int) {
	// for each shot fired
	for pi, pv := range pp.ShotsFired {
		// if coordinates == positions
		if p.Abscissa == pv.Abscissa && p.Ordinate == pv.Ordinate {
			// return true, positions index in Player ShotsFired
			return true, pi
		}
	}
	return false, -1
}

// ###########################################################################################################
// ######################################### TEST METHODS STRINGIFIER ########################################
// ###########################################################################################################

func main() {
	s := PrepareSea(10, 5)
	fmt.Println(s.PrettyPrintSeaInfo())
	g := PrepareGame(10, 5)
	fmt.Println(g.FirstPlayer.SeaToString(-1))
	fmt.Println(g.FirstPlayer.SeaToString(-1))
	g.FirstPlayer.GunShot(&g.SecondPlayer, &g.SecondPlayer.Sea.Ships[0].Positions[0])
	fmt.Println(g.FirstPlayer.SeaToString(-1))
}
