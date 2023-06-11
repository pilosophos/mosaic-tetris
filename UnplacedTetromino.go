package main

import (
	"strconv"
	"strings"
)

type UnplacedTetromino struct {
	BlockRelativeXYs [][2]int
	TopLeftXY        [2]int
	TimeLeft         int
	Color            string
}

func NewUnplacedTetromino(blocksRelativeXY [][2]int, topLeftXY [2]int, timeLeft int, color string) *UnplacedTetromino {
	return &UnplacedTetromino{
		blocksRelativeXY,
		topLeftXY,
		timeLeft,
		color,
	}
}

func (tetromino *UnplacedTetromino) Tick() int {
	tetromino.TimeLeft--
	return tetromino.TimeLeft
}

// Get the global coordinates of this tetromino's blocks
func (tetromino UnplacedTetromino) BlockGlobalXYs() (globalXYs [][2]int) {
	for _, blockXY := range tetromino.BlockRelativeXYs {
		blockGlobalXY := [2]int{
			blockXY[0] + tetromino.TopLeftXY[0],
			blockXY[1] + tetromino.TopLeftXY[1],
		}
		globalXYs = append(globalXYs, blockGlobalXY)
	}
	return globalXYs
}

func (tetromino *UnplacedTetromino) Translate(dx int, dy int) {
	tetromino.TopLeftXY[0] += dx
	tetromino.TopLeftXY[1] += dy
}

func (tetromino UnplacedTetromino) BlockString() string {
	return tetromino.Color + strconv.Itoa(tetromino.TimeLeft) + "\033[0m"
}

func (tetromino UnplacedTetromino) String() string {
	squares := make([][4]string, 4)

	for x := range squares {
		for y := range squares[x] {
			squares[y][x] = " "
		}
	}

	for _, blockXY := range tetromino.BlockRelativeXYs {
		x, y := blockXY[0], blockXY[1]
		squares[y][x] = tetromino.BlockString()
	}

	rows := make([]string, 4)
	for x := range squares {
		rows[x] = strings.Join(squares[x][:], "")
	}

	return strings.Join(rows, "\n")
}
