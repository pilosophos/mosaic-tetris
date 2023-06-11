package main

import (
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Board struct {
	Blocks            [][]*Block         // 2D array of the blocks on the board [0]=height/y/row, [1]=width/x/col
	HoveringTetromino *UnplacedTetromino // The tetromino currently hovering over the board
	IllegalBlocks     [][2]int           // Coordinates the hovering tetromino is above that already have blocks on them
	Width             int                // Width of the board
	LinesCleared      int                // Number of lines cleared
	Score             int                // Guess what this might represent
	Message           string             // Message to display after clearing one or more lines
}

// Create a new Mosaic Tetris board
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
		"",
	}
}

// Clear any rows that are full of blocks and update Score/LinesCleared/Message accordingly
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
	board.Message = []string{"         ", "Single", "Double", "Triple", "Tetris!"}[linesCleared]
}

// Try placing the tetromino on the board, returning true if it succeeded
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

// Hover a tetromino over the board and calculate positions where it would
// collide with an already placed block
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

// Get a string representation of the board
func (board Board) String() string {
	rows := []string{}

	for y, row := range board.Blocks {
		rows = append(rows, "")

		for x, block := range row {
			if slices.Contains(board.IllegalBlocks, [2]int{x, y}) {
				// rows[y] += "\033[1;31m" + "X" + "\033[0m"
				rows[y] += "X"
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

	rows = append(rows, "Lines\t"+strconv.Itoa(board.LinesCleared))
	rows = append(rows, "Score\t"+strconv.Itoa(board.Score))
	rows = append(rows, board.Message)

	return strings.Join(rows, "\n")
}
