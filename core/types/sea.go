package types

import "strconv"

type Sea struct {
	Dimension int    `json:"dim"`
	Ships     []Ship `json:"shp"`
}

// PrettyPrintSeaInfo return Sea string info
//
//	[s:*Ship]	Ship point pointer
//	[return]	string
func (s *Sea) PrettyPrintSeaInfo() string {
	ss := "Sea dimensions: " + strconv.Itoa(s.Dimension) + "\n"
	for _, sv := range s.Ships {
		if sv.Dimension != 0 {
			ss += sv.PrettyPrintShipInfo() + "\n"
		}
	}
	return ss
}

// CheckShipPosition check if in p coordinates in given Sea there's a Ship
//
//	[p:*Coordinates]	Coordinate point pointer		[s:*Sea]		b Sea pointer
//	[return]	bool (collision), ship index, coordinate index
func (s *Sea) CheckShipPosition(p *Coordinates) (bool, int, int) {
	// for each Ships
	for si, sv := range s.Ships {
		// for each Positions occuped by
		for ci, cv := range sv.Positions {
			// if coordinates == positions
			if p.Abscissa == cv.Abscissa && p.Ordinate == cv.Ordinate {
				// return true, ship index, positions index in ship struct
				return true, si, ci
			}
		}
	}
	return false, -1, -1
}
