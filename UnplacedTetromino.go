package main

import (
	"math"
	"strconv"
	"strings"
)

type UnplacedTetromino struct {
	BlockRelativeXYs [][2]int
	TopLeftXY        [2]int
	TimeLeft         int
	Color            string
	Width            int
	Height           int
}

func NewUnplacedTetromino(blocksRelativeXY [][2]int, topLeftXY [2]int, timeLeft int, color string) (tetromino *UnplacedTetromino) {
	tetromino = &UnplacedTetromino{
		blocksRelativeXY,
		topLeftXY,
		timeLeft,
		color,
		0, 0,
	}
	tetromino.ComputeDimensions()
	return tetromino
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

func (tetromino *UnplacedTetromino) Translate(dx int, dy int, boardWidth int, boardHeight int) {
	xlimit := boardWidth - tetromino.Width
	ylimit := boardHeight - tetromino.Height

	newx := tetromino.TopLeftXY[0] + dx
	if newx < 0 {
		newx = 0
	} else if newx > xlimit {
		newx = xlimit
	}
	tetromino.TopLeftXY[0] = newx

	newy := tetromino.TopLeftXY[1] + dy
	if newy < 0 {
		newy = 0
	} else if newy > ylimit {
		newy = ylimit
	}
	tetromino.TopLeftXY[1] = newy
}

func (tetromino UnplacedTetromino) BlockString() string {
	return tetromino.Color + strconv.Itoa(tetromino.TimeLeft) + "\033[0m"
}

func (tetromino *UnplacedTetromino) Rotate(angleDegrees int) {
	angle := float64(angleDegrees) * math.Pi / 180
	xmin := 0
	ymin := 0
	for i, xy := range tetromino.BlockRelativeXYs {
		newx := int(math.Cos(angle))*xy[0] - int(math.Sin(angle))*xy[1]
		newy := int(math.Sin(angle))*xy[0] + int(math.Cos(angle))*xy[1]
		tetromino.BlockRelativeXYs[i] = [2]int{newx, newy}

		if newx < xmin {
			xmin = newx
		}
		if newy < ymin {
			ymin = newy
		}
	}

	// keep all coordinates positive
	for i := range tetromino.BlockRelativeXYs {
		if xmin < 0 {
			tetromino.BlockRelativeXYs[i][0] += -xmin
		}
		if ymin < 0 {
			tetromino.BlockRelativeXYs[i][1] += -ymin
		}
	}

	tetromino.ComputeDimensions()
}

func (tetromino *UnplacedTetromino) ComputeDimensions() {
	xmax := 0
	ymax := 0

	for _, xy := range tetromino.BlockRelativeXYs {
		if xy[0] > xmax {
			xmax = xy[0]
		}
		if xy[1] > ymax {
			ymax = xy[1]
		}
	}
	tetromino.Width = xmax + 1
	tetromino.Height = ymax + 1
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
		rows[x] = strings.Join(squares[x][:], " ")
	}

	return strings.Join(rows, "\n")
}
