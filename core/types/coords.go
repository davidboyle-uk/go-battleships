package types

import (
	"fmt"
	"strconv"

	"github.com/dbx123/go-battleships/util"
)

type Coordinates struct {
	Abscissa int `json:"x"`
	Ordinate int `json:"y"`
	Status   int `json:"s"`
}

func (c Coordinates) Validate(gridSize int) bool {
	if c.Abscissa < 1 || c.Abscissa > gridSize {
		return false
	}
	if c.Ordinate < 1 || c.Ordinate > gridSize {
		return false
	}
	return true
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%v:%v", c.Abscissa, c.Ordinate)
}

// PrettyPrintCoordinatesInfo return String representation of Coordinates
//
//	[p:*Coordinates]	Coordinates point pointer
//	[return]	string
func (p *Coordinates) PrettyPrintCoordinatesInfo() string {
	return "[" + util.IntToLetter(p.Abscissa-1) + strconv.Itoa(p.Ordinate) + "]"
}
