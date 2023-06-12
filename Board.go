package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/exp/slices"
)

type Board struct {
	Blocks            [][]*Block         // 2D array of the blocks on the board [0]=height/y/row, [1]=width/x/col
	HoveringTetromino *UnplacedTetromino // The tetromino currently hovering over the board
	IllegalBlocks     [][2]int           // Coordinates the hovering tetromino is above that already have blocks on them
	Width             int                // Width of the board
	Height            int                // Height of the board
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
		height,
		0,
		0,
		"",
	}
}

// Clear any rows that are full of blocks and update Score/LinesCleared/Message accordingly
func (board *Board) ClearFullRows() {
	// detect full rows
	clearRows := []int{}
	for row, rowBlocks := range board.Blocks {
		if !slices.Contains(rowBlocks, nil) {
			clearRows = append(clearRows, row)
		}
	}

	// detect full columns
	clearCols := []int{}
	for col := 0; col < board.Width; col++ {
		verticalClear := true
		for row := 0; row < board.Height; row++ {
			if board.Blocks[row][col] == nil {
				verticalClear = false
			}
		}
		if verticalClear {
			clearCols = append(clearCols, col)
		}
	}

	// clear full rows/columns
	for _, col := range clearCols {
		for row := 0; row < board.Height; row++ {
			board.Blocks[row][col] = nil
		}
	}
	for _, row := range clearRows {
		board.Blocks[row] = make([]*Block, board.Width)
	}

	linesCleared := len(clearRows) + len(clearCols)
	board.LinesCleared += linesCleared
	board.Score += []int{0, 100, 300, 500, 800, 1100, 1400}[linesCleared]
	board.Message = []string{
		"            ",
		"Single      ",
		"Double!     ",
		"Triple!!    ",
		"Tetris!!!   ",
		"MOSAIC!!!!  ",
		"SUPER MOSAIC"}[linesCleared]
}

// Try placing the tetromino on the board, returning true if it succeeded
func (board *Board) PlaceTetromino(tetromino *UnplacedTetromino) bool {
	if len(board.IllegalBlocks) == 0 {
		for _, blockXY := range tetromino.BlockGlobalXYs() {
			x, y := blockXY[0], blockXY[1]
			board.Blocks[y][x] = NewBlock(tetromino.Color, '▄')
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

type drawTextFunc func(tcell.Screen, int, int, tcell.Style, string)

// Render the board on the screen
func (board Board) Render(s tcell.Screen, defaultStyle tcell.Style, drawText drawTextFunc, topLeftX, topLeftY int) {
	colors := map[string]tcell.Style{
		"cyan":    tcell.StyleDefault.Foreground(tcell.ColorAqua),
		"white":   tcell.StyleDefault.Foreground(tcell.ColorWhite),
		"magenta": tcell.StyleDefault.Foreground(tcell.ColorFuchsia),
		"blue":    tcell.StyleDefault.Foreground(tcell.ColorBlue),
		"yellow":  tcell.StyleDefault.Foreground(tcell.ColorYellow),
		"green":   tcell.StyleDefault.Foreground(tcell.ColorGreen),
		"red":     tcell.StyleDefault.Foreground(tcell.ColorRed),
	}

	cursorY := topLeftY
	for y, row := range board.Blocks {
		s.SetContent(topLeftX, y+topLeftY, ' ', nil, defaultStyle)

		for x, block := range row {
			if slices.Contains(board.IllegalBlocks, [2]int{x, y}) {
				s.SetContent((x*2)+topLeftX, y+topLeftY, 'X', nil, colors["red"])
			} else if board.HoveringTetromino != nil && slices.Contains(board.HoveringTetromino.BlockGlobalXYs(), [2]int{x, y}) {
				s.SetContent((x*2)+topLeftX, y+topLeftY, board.HoveringTetromino.BlockRune(), nil, colors[board.HoveringTetromino.Color])
			} else if block == nil {
				s.SetContent((x*2)+topLeftX, y+topLeftY, '.', nil, defaultStyle)
			} else {
				s.SetContent((x*2)+topLeftX, y+topLeftY, block.Rune(), nil, colors[block.Color])
			}
		}
		cursorY = y + topLeftY
	}

	drawText(s, 0, cursorY+1, defaultStyle, "Lines\t"+strconv.Itoa(board.LinesCleared))
	drawText(s, 0, cursorY+2, defaultStyle, "Score\t"+strconv.Itoa(board.Score))
}
