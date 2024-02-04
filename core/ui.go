package core

import (
	"go-battleships/core/pallete"
)

var (
	WINNER_TEXT = pallete.Red(`
	__      ___ _ __  _ __   ___ _ __ 
	\ \ /\ / / | '_ \| '_ \ / _ \ '__|
	 \ V  V /| | | | | | | |  __/ |   
	  \_/\_/ |_|_| |_|_| |_|\___|_|

`)
	LOSER_TEXT = pallete.Cyan(`
	 _                     
	| |                    
	| | ___  ___  ___ _ __ 
	| |/ _ \/ __|/ _ \ '__|
	| | (_) \__ \  __/ |   
	|_|\___/|___/\___|_|

`)
)
