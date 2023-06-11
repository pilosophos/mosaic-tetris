package main

import (
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/inancgumus/screen"
)

const BoardSizeW = 10
const BoardSizeH = 20

func main() {
	keysEvents, err := keyboard.GetKeys(1)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	// make a test tetromino
	tet := NewUnplacedTetromino(
		[][2]int{
			{0, 0},
			{0, 1},
			{0, 2},
			{1, 0},
		},
		[2]int{0, 0},
		4,
		"\033[1;34m",
	)

	board := NewBoard(BoardSizeW, BoardSizeH)

	screen.Clear()

	// start global tick timer
	tickTimer := make(chan bool)
	go tickGameForever(tickTimer)

	for {
		screen.MoveTopLeft()
		board.HoverTetromino(tet)
		fmt.Println(board)

		select {
		case event := <-keysEvents:
			if event.Err != nil {
				panic(event.Err)
			}
			handleKeypress(event.Key, tet, board)
		case tick := <-tickTimer:
			if tick {
				tet.Tick()
			}
		}
	}
}

func tickGameForever(tick chan bool) {
	for {
		time.Sleep(1 * time.Second)
		tick <- true
	}
}

func handleKeypress(key keyboard.Key, hoveringTetromino *UnplacedTetromino, board *Board) {
	switch key {
	case keyboard.KeyEsc:
		os.Exit(0)
	case keyboard.KeyArrowLeft:
		hoveringTetromino.Translate(-1, 0)
	case keyboard.KeyArrowRight:
		hoveringTetromino.Translate(1, 0)
	case keyboard.KeyArrowUp:
		hoveringTetromino.Translate(0, -1)
	case keyboard.KeyArrowDown:
		hoveringTetromino.Translate(0, 1)
	case keyboard.KeySpace:
		board.PlaceTetromino(hoveringTetromino)
	}
}
