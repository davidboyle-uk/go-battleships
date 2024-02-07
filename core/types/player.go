package types

import (
	"strconv"

	"github.com/dbx123/go-battleships/core/palette"
	"github.com/dbx123/go-battleships/util"
)

const (
	STATUS_SHIP_OK = iota
	STATUS_SHIP_HIT
	STATUS_SEA_OK
	STATUS_SEA_HIT
)

var (
	GAME_GRID_BORDER = "|"
	STR_SHIP_OK      = palette.White("██")
	STR_SHIP_HIT     = palette.Red("██")
	STR_SEA_OK       = "  "
	STR_SEA_HIT      = palette.Cyan("~~")
	STR_STATUS_ERROR = "??"

	PRINT_CALLERSEA_MODE   = 2
	PRINT_OPPONENTSEA_MODE = 1
)

var legend = map[int]string{
	0: "Legend:",
	1: "Ship (OK)        " + STR_SHIP_OK,
	2: "Ship (HIT)       " + STR_SHIP_HIT,
	3: "Sea (MISS)       " + STR_SEA_HIT,
}

type Player struct {
	Name       string        `json:"id"`
	Sea        Sea           `json:"sea"`
	ShotsFired []Coordinates `json:"sht"`
	Suffered   []Coordinates `json:"rcv"`
	Hits       []Coordinates `json:"hit"`
}

// SeaToString print Player Sea.
//
//	[p:*Player]	Player
//	[return]	string
func (p *Player) SeaToString(h int) (ss string) {
	// create column indicator line
	ss = "   " + GAME_GRID_BORDER
	for r := 0; r < p.Sea.Dimension-1; r++ {
		ss += "-" + util.IntToLetter(r) + "--" + GAME_GRID_BORDER
	}
	ss += "-" + util.IntToLetter(p.Sea.Dimension-1) + "--" + GAME_GRID_BORDER + "\n"

	// create first separation line
	ss += "   " + GAME_GRID_BORDER
	for r := 0; r < p.Sea.Dimension-1; r++ {
		ss += "-----"
	}
	ss += "----" + GAME_GRID_BORDER + "\n"

	// for each row
	for r := 0; r < p.Sea.Dimension; r++ {
		// start with legend
		pad := "  "
		if r >= 9 {
			pad = " "
		}
		ss += strconv.Itoa(r+1) + pad

		// add grid border
		ss += GAME_GRID_BORDER

		// for each column
		for c := 0; c < p.Sea.Dimension; c++ {
			// if we are drawing caller's Sea
			if h == PRINT_CALLERSEA_MODE {
				// check ShipPosition in Sea
				rp, si, ci := p.Sea.CheckShipPosition(&Coordinates{Abscissa: c + 1, Ordinate: r + 1})

				// if there's a Sea in position
				if rp {
					// add correct status representation
					ss += " " + StatusToString(p.Sea.Ships[si].Positions[ci].Status) + " " + GAME_GRID_BORDER
				} else {
					// check SufferedMoves in Sea
					pp, pi := p.CheckSufferedMoves(&Coordinates{Abscissa: c + 1, Ordinate: r + 1})

					// if opponent shot in the cell
					if pp {
						// add correct status representation
						ss += " " + StatusToString(p.Suffered[pi].Status) + " " + GAME_GRID_BORDER
					} else {
						ss += " " + STR_SEA_OK + " " + GAME_GRID_BORDER
					}
				}
			}

			// if we are drawing opponent's Sea
			if h == PRINT_OPPONENTSEA_MODE {
				// check SufferedMoves in Sea
				pp, pi := p.CheckSufferedMoves(&Coordinates{Abscissa: c + 1, Ordinate: r + 1})

				// if opponent shot in the cell
				if pp {
					// add correct status representation
					ss += " " + StatusToString(p.Suffered[pi].Status) + " " + GAME_GRID_BORDER
				} else {
					ss += " " + STR_SEA_OK + " " + GAME_GRID_BORDER
				}
			}
		}

		// create separation line
		ss += "\n" + "   " + GAME_GRID_BORDER
		l := "" // legend text
		for c := 0; c < p.Sea.Dimension-1; c++ {
			ss += "-----"
			if lt, ok := legend[r]; ok {
				l = lt
			}
		}
		ss += "----" + GAME_GRID_BORDER + "   " + l + "\n"
	}

	return ss
}

// StatusToString print Status of Coordinates.
func StatusToString(s int) string {
	// check Ship status in specific position
	switch s {
	// ship && sea status
	case STATUS_SHIP_HIT:
		return STR_SHIP_HIT
	case STATUS_SHIP_OK:
		return STR_SHIP_OK
	case STATUS_SEA_HIT:
		return STR_SEA_HIT
	case STATUS_SEA_OK:
		return STR_SEA_OK
	}

	return STR_STATUS_ERROR
}

// CheckSufferedMoves check p coordinates in given Sea's Player.
//
//	[p:*Coordinates]	Coordinate point pointer		[pp:*Player]		b Player pointer
//	[return]	bool (collision), ship index, coordinate index
func (pp *Player) CheckSufferedMoves(p *Coordinates) (bool, int) {
	// for each suffered
	for pi, pv := range pp.Suffered {
		// if coordinates == positions
		if p.Abscissa == pv.Abscissa && p.Ordinate == pv.Ordinate {
			// return true, positions index in Player Suffered Moves
			return true, pi
		}
	}
	return false, -1
}

// CheckSufferedMoves check p coordinates in given Sea's Player.
//
//	[p:*Coordinates]	Coordinate point pointer		[pp:*Player]		b Player pointer
//	[return]	bool (collision), ship index, coordinate index
func (pp *Player) CheckHitMoves(p *Coordinates) (bool, int) {
	// for each suffered
	for pi, pv := range pp.Hits {
		// if coordinates == positions
		if p.Abscissa == pv.Abscissa && p.Ordinate == pv.Ordinate {
			// return true, positions index in Player Suffered Moves
			return true, pi
		}
	}
	return false, -1
}

// GunShot from p Player to t Player in p Coordinates.
//
//	[f:*Player]			from Player	pointer		//	[t:*Player]			to Player pointer
//	[p:*Coordinates]	Coordinate point pointer
func (f *Player) GunShot(t *Player, p *Coordinates) {
	// check if f player hit t player in position p
	rs, si, ci := t.Sea.CheckShipPosition(p)

	// record the shot
	f.ShotsFired = append(f.ShotsFired, *p)

	var np Coordinates
	// if Ship hit
	if rs {
		// t player ship stricken
		t.Sea.Ships[si].Positions[ci].Status = STATUS_SHIP_HIT
		np = Coordinates{Abscissa: p.Abscissa, Ordinate: p.Ordinate, Status: STATUS_SHIP_HIT}
		if h, _ := f.CheckHitMoves(p); !h {
			f.Hits = append(f.Hits, np)
		}
	} else {
		np = Coordinates{Abscissa: p.Abscissa, Ordinate: p.Ordinate, Status: STATUS_SEA_HIT}
	}

	// add suffered move
	t.Suffered = append(t.Suffered, np)
}
