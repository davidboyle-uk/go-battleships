package util

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	ROW_NUMBER = 1024
	COL_NUMBER = 1024

	BLANK_SPACE = " "
	PAUSE_MEX   = ">>> press ENTER to go on..."
)

func Random(min, max int) int {
	max = max + 1
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// a = 1
func RuneToInt(letter rune) int {
	char := int(letter)
	char -= 65 + 31
	return char
}

// 1 = a
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

func ConsolePause(m string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(m)
	reader.ReadString('\n')
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

func Exit(w http.ResponseWriter, r *http.Request) {
	CleanScreen()
	os.Exit(1)
}
