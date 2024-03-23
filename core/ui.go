package core

import (
	"github.com/davidboyle-uk/go-battleships/palette"
)

//nolint:dupword
var (
	WINNER_TEXT = palette.Red(`
	__      ___ _ __  _ __   ___ _ __ 
	\ \ /\ / / | '_ \| '_ \ / _ \ '__|
	 \ V  V /| | | | | | | |  __/ |   
	  \_/\_/ |_|_| |_|_| |_|\___|_|

`)
	LOSER_TEXT = palette.Cyan(`
	 _                     
	| |                    
	| | ___  ___  ___ _ __ 
	| |/ _ \/ __|/ _ \ '__|
	| | (_) \__ \  __/ |   
	|_|\___/|___/\___|_|

`)
)
