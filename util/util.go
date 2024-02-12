package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	ROW_NUMBER int = 1024
	COL_NUMBER int = 1024

	BLANK_SPACE string = " "
	PAUSE_MEX   string = ">>> press ENTER to go on..."
)

func Random(min, max int) int {
	if min == max {
		return min
	}
	// calculate the max we will be using
	bg := big.NewInt(int64(max - min))
	n, _ := rand.Int(rand.Reader, bg)
	return int(n.Int64()) + min
}

// a = 0.
func RuneToInt(letter rune) int {
	char := int(letter)
	char -= 65 + 32
	return char
}

// 0 = a.
func intToRune(number int) rune {
	char := number + 65 + 32
	return rune(char)
}

func IntToLetter(number int) string {
	return string(intToRune(number))
}

func Search(a int, b []int) bool {
	for _, v := range b {
		if v == a {
			return true
		}
	}
	return false
}

func CleanScreen() {
	r, c := 0, 0
	for r < ROW_NUMBER {
		for c < COL_NUMBER {
			fmt.Print(BLANK_SPACE)
			c++
		}
		fmt.Println(BLANK_SPACE)
		r++
	}
	fmt.Print("\033[0;0H")
}
