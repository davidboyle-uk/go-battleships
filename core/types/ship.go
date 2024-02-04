package types

import "strconv"

type Ship struct {
	Dimension int           `json:"ds"`
	Positions []Coordinates `json:"ps"`
}

// PrettyPrintShipInfo return Ship string info
//
//	[s:*Ship]	Ship point pointer
//	[return]	string
func (s *Ship) PrettyPrintShipInfo() string {
	ss := "\tShip dimensions: " + strconv.Itoa(s.Dimension) + "\n\t\t["
	for _, pv := range s.Positions {
		ss += pv.PrettyPrintCoordinatesInfo() + " "
	}
	ss += "]"
	return ss
}
