package main

import (
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Board struct {
	Blocks            [][]*Block // height/y/row, width/x/col
	HoveringTetromino *UnplacedTetromino
	IllegalBlocks     [][2]int // coordinates the hovering tetromino above that already have blocks on them
	Width             int
	LinesCleared      int
	Score             int
}

func NewBoard(width int, height int) *Board {
	blocks := make([][]*Block, height)
	for i := range blocks {
		blocks[i] = make([]*Block, width)
	}

	return &Board{
		blocks,
		nil,
		make([][2]int, 0),
		width,
		0,
		0,
	}
}

func (board *Board) ClearFullRows() {
	linesCleared := 0
	for row, rowBlocks := range board.Blocks {
		if !slices.Contains(rowBlocks, nil) {
			board.Blocks[row] = make([]*Block, board.Width)
			linesCleared++
		}
	}

	board.LinesCleared += linesCleared
	board.Score += []int{0, 100, 300, 500, 800}[linesCleared]
}

func (board *Board) PlaceTetromino(tetromino *UnplacedTetromino) bool {
	if len(board.IllegalBlocks) == 0 {
		for _, blockXY := range tetromino.BlockGlobalXYs() {
			x, y := blockXY[0], blockXY[1]
			board.Blocks[y][x] = NewBlock(tetromino.Color, "▅")
		}
		board.Score += 2 * (tetromino.TimeLeft + 1)
		board.ClearFullRows()
		return true
	}
	return false
}

func (board *Board) HoverTetromino(tetromino *UnplacedTetromino) {
	board.IllegalBlocks = make([][2]int, 0)
	board.HoveringTetromino = tetromino

	for _, hoverBlockXY := range tetromino.BlockGlobalXYs() {
		hoverX, hoverY := hoverBlockXY[0], hoverBlockXY[1]
		if board.Blocks[hoverY][hoverX] != nil {
			board.IllegalBlocks = append(board.IllegalBlocks, [2]int{hoverX, hoverY})
		}
	}
}

func (board Board) String() string {
	rows := []string{}

	for y, row := range board.Blocks {
		rows = append(rows, "")

		for x, block := range row {
			if slices.Contains(board.IllegalBlocks, [2]int{x, y}) {
				rows[y] += "\033[1;31m" + "X" + "\033[0m"
			} else if board.HoveringTetromino != nil && slices.Contains(board.HoveringTetromino.BlockGlobalXYs(), [2]int{x, y}) {
				rows[y] += board.HoveringTetromino.BlockString()
			} else if block == nil {
				rows[y] += "."
			} else {
				rows[y] += block.String()
			}

			rows[y] += " "
		}
	}

	return strings.Join(rows, "\n") + "\nLines cleared: " + strconv.Itoa(board.LinesCleared) + "\nScore: " + strconv.Itoa(board.Score)
}
