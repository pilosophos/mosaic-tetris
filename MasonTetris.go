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

	tetrominoQueue := NewTetrominoQueue()
	hoveringTetromino := tetrominoQueue.Pop()

	board := NewBoard(BoardSizeW, BoardSizeH)

	screen.Clear()

	// start global tick timer
	tickTimer := make(chan bool)
	go tickGameForever(tickTimer)

	for {
		board.HoverTetromino(hoveringTetromino)

		screen.MoveTopLeft()
		// screen.Clear()
		fmt.Println("NEXT")
		fmt.Println(tetrominoQueue.Peek())
		fmt.Println(board)

		select {
		case event := <-keysEvents:
			if event.Err != nil {
				panic(event.Err)
			}
			tetrominoPlaced := handleKeypress(event.Key, hoveringTetromino, board)
			if tetrominoPlaced {
				hoveringTetromino = tetrominoQueue.Pop()
			}
		case <-tickTimer:
			timeleft := hoveringTetromino.Tick()
			if timeleft == 0 {
				tetrominoPlaced := board.PlaceTetromino(hoveringTetromino)
				if !tetrominoPlaced {
					fmt.Println("You lose!")
					os.Exit(0)
				}
				hoveringTetromino = tetrominoQueue.Pop()
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

func handleKeypress(key keyboard.Key, hoveringTetromino *UnplacedTetromino, board *Board) (placed bool) {
	switch key {
	case keyboard.KeyEsc:
		os.Exit(0)
	case keyboard.KeyArrowLeft:
		hoveringTetromino.Translate(-1, 0, BoardSizeW, BoardSizeH)
	case keyboard.KeyArrowRight:
		hoveringTetromino.Translate(1, 0, BoardSizeW, BoardSizeH)
	case keyboard.KeyArrowUp:
		hoveringTetromino.Translate(0, -1, BoardSizeW, BoardSizeH)
	case keyboard.KeyArrowDown:
		hoveringTetromino.Translate(0, 1, BoardSizeW, BoardSizeH)
	case keyboard.KeySpace:
		return board.PlaceTetromino(hoveringTetromino)
	}
	return false
}
