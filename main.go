package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type slot struct {
	value    int
	hasMine  bool
	isFaceUp bool
	posX     int
	posY     int
}

var colorYellow string = "\033[33m\033[1m"
var colorReset string = "\033[0m"

func buildBoard(sizeX int, sizeY int) [][]slot {
	board := make([][]slot, sizeX)
	fmt.Printf("Creating a %dx%d board\n", sizeX, sizeY)
	for i := 0; i < sizeX; i++ {
		board[i] = make([]slot, sizeY)

		for j := 0; j < sizeY; j++ {
			board[i][j] = slot{
				posX: i,
				posY: j,
			}
		}
	}

	return board
}

func getRandomNumberFromSlice(s []int) (int, []int) {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(s))
	ret := s[index]

	return ret, append(s[:index], s[index+1:]...)

}

// Whenerver ew add a mine, lets increase the vaue of the neighbors. Can probably be done in a smarte way
func increaseSlotValues(x int, y int, board [][]slot, sizeX, sizeY int) {

	// start at top left and move clockwise
	xToTest := x - 1
	yToTest := y - 1

	if xToTest >= 0 && yToTest >= 0 && xToTest < sizeX && yToTest < sizeY && board[xToTest][yToTest].hasMine == false {
		board[xToTest][yToTest].value += 1
	}

	xToTest = x
	yToTest = y - 1

	if xToTest >= 0 && yToTest >= 0 && xToTest < sizeX && yToTest < sizeY && board[xToTest][yToTest].hasMine == false {
		board[xToTest][yToTest].value += 1
	}

	xToTest = x + 1
	yToTest = y - 1

	if xToTest >= 0 && yToTest >= 0 && xToTest < sizeX && yToTest < sizeY && board[xToTest][yToTest].hasMine == false {
		board[xToTest][yToTest].value += 1
	}

	xToTest = x + 1
	yToTest = y

	if xToTest >= 0 && yToTest >= 0 && xToTest < sizeX && yToTest < sizeY && board[xToTest][yToTest].hasMine == false {
		board[xToTest][yToTest].value += 1
	}

	xToTest = x + 1
	yToTest = y + 1

	if xToTest >= 0 && yToTest >= 0 && xToTest < sizeX && yToTest < sizeY && board[xToTest][yToTest].hasMine == false {
		board[xToTest][yToTest].value += 1
	}

	xToTest = x
	yToTest = y + 1

	if xToTest >= 0 && yToTest >= 0 && xToTest < sizeX && yToTest < sizeY && board[xToTest][yToTest].hasMine == false {
		board[xToTest][yToTest].value += 1
	}

	xToTest = x - 1
	yToTest = y + 1

	if xToTest >= 0 && yToTest >= 0 && xToTest < sizeX && yToTest < sizeY && board[xToTest][yToTest].hasMine == false {
		board[xToTest][yToTest].value += 1
	}

	xToTest = x - 1
	yToTest = y

	if xToTest >= 0 && yToTest >= 0 && xToTest < sizeX && yToTest < sizeY && board[xToTest][yToTest].hasMine == false {
		board[xToTest][yToTest].value += 1
	}
}

// Pick ana vailable number from set, flag the proper slot and remove from set
func addMines(totalMines int, board [][]slot, sizeX int, sizeY int) [][]slot {
	var index int
	availableNumbers := make([]int, sizeX*sizeY)
	for i := 0; i < sizeX*sizeY; i++ {
		availableNumbers[i] = i
	}
	for totalMines > 0 {

		index, availableNumbers = getRandomNumberFromSlice(availableNumbers)

		x := index / sizeY
		y := index % sizeY

		board[x][y].hasMine = true
		board[x][y].value = 0

		increaseSlotValues(x, y, board, sizeX, sizeY)
		totalMines--
	}

	return board
}

func printBoard(board [][]slot, allUncovered bool, sizeX, sizeY, userX, userY int) {
	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			s := board[i][j]
			var prefix string
			if i == userX && j == userY {
				prefix = colorYellow
			} else {
				prefix = colorReset
			}
			if s.isFaceUp || allUncovered {
				if s.hasMine {
					fmt.Printf("%s* ", prefix)
				} else {
					fmt.Printf("%s%d ", prefix, s.value)
				}
			} else {
				fmt.Printf("%s? ", prefix)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Print("\n")
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// DFS to clear out all the adjacent zeroes
func clearZeroesFromPosition(board [][]slot, x, y, boardSizeX, boardSizeY int) {
	if x < 0 || y < 0 || x >= boardSizeX || y >= boardSizeY || board[x][y].isFaceUp || board[x][y].value != 0 || board[x][y].hasMine {
		return
	}

	board[x][y].isFaceUp = true
	clearZeroesFromPosition(board, x-1, y, boardSizeX, boardSizeY)
	clearZeroesFromPosition(board, x, y-1, boardSizeX, boardSizeY)
	clearZeroesFromPosition(board, x+1, y, boardSizeX, boardSizeY)
	clearZeroesFromPosition(board, x, y+1, boardSizeX, boardSizeY)
}

func main() {
	var totalMines int
	var boardSizeX int
	var boardSizeY int
	var faceUpSlots int = 0

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Difficulty (1..3): ")
	text, _ := reader.ReadString('\n')
	diff, err := strconv.ParseInt(strings.Trim(text, "\n"), 10, 0)
	for err != nil {
		fmt.Printf("Invalid difficulty %s, please choose a number between 1 and 3\n", text)
		fmt.Print("Enter Difficulty (1..3): ")
		text, _ := reader.ReadString('\n')
		diff, err = strconv.ParseInt(strings.Trim(text, "\n"), 10, 0)
	}

	switch diff {
	case 1:
		totalMines = 10
		boardSizeX = 10
		boardSizeY = 10
	case 2:
		totalMines = 40
		boardSizeX = 16
		boardSizeY = 16
	case 3:
		totalMines = 99
		boardSizeX = 30
		boardSizeY = 16
	}

	board := buildBoard(boardSizeX, boardSizeY)
	addMines(totalMines, board, boardSizeX, boardSizeY)

	userRow := 0
	userColumn := 0
	for {
		clearScreen()
		printBoard(board, false, boardSizeX, boardSizeY, userRow, userColumn)

		_, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyArrowLeft:
			if userColumn > 0 {
				userColumn -= 1
			}
		case keyboard.KeyArrowDown:
			if userRow < boardSizeX-1 {
				userRow += 1
			}
		case keyboard.KeyArrowRight:
			if userColumn < boardSizeY-1 {
				userColumn += 1
			}
		case keyboard.KeyArrowUp:
			if userRow > 0 {
				userRow -= 1
			}
		case keyboard.KeyEsc:
			fmt.Println("Bye bye!")
			return
		case keyboard.KeyEnter:
			if board[userRow][userColumn].isFaceUp {
				continue
			}

			if board[userRow][userColumn].hasMine {
				printBoard(board, true, boardSizeX, boardSizeY, userRow, userColumn)
				fmt.Print("You lost!")
				return
			}

			if !board[userRow][userColumn].isFaceUp {
				if board[userRow][userColumn].value == 0 {
					clearZeroesFromPosition(board, userRow, userColumn, boardSizeX, boardSizeY)
					continue
				} else {
					board[userRow][userColumn].isFaceUp = true
				}

				faceUpSlots = countFaceUpSlots(board)

				if faceUpSlots == (boardSizeX*boardSizeY)-totalMines {
					printBoard(board, true, boardSizeX, boardSizeY, userRow, userColumn)
					fmt.Print("You won!")
					return
				}
			}
		}
	}
}

// Count all the the slots that haven't been flipped. We can probbaly keep a running count instead of re counting all the time
func countFaceUpSlots(board [][]slot) int {
	faceUp := 0

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			if board[i][j].isFaceUp {
				faceUp += 1
			}
		}
	}

	return faceUp
}
