package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const BoardSizeW = 10
const BoardSizeH = 20

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.Clear()

	tetrominoQueue := NewTetrominoQueue()
	hoveringTetromino := tetrominoQueue.Pop()

	board := NewBoard(BoardSizeW, BoardSizeH)

	// start global tick timer
	tickTimer := make(chan bool)
	go tickGameForever(tickTimer)

	termEvents := make(chan tcell.Event)
	tcellQuit := make(chan struct{})
	go s.ChannelEvents(termEvents, tcellQuit)

	for {
		board.HoverTetromino(hoveringTetromino)

		s.Show()
		drawText(s, (BoardSizeW*2)+2, 1, defStyle, "NEXT")
		drawText(s, (BoardSizeW*2)+2, 2, defStyle, tetrominoQueue.Peek().String())
		board.Render(s, defStyle, drawText, 0, 0)

		select {
		case <-tickTimer:
			timeleft := hoveringTetromino.Tick()
			if timeleft == 0 {
				tetrominoPlaced := board.PlaceTetromino(hoveringTetromino)
				if !tetrominoPlaced {
					fmt.Println("You lose!")
					quit()
				}
				hoveringTetromino = tetrominoQueue.Pop()
			}
		case ev := <-termEvents:
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				placed := handleKeypress(ev, quit, hoveringTetromino, board)
				if placed {
					hoveringTetromino = tetrominoQueue.Pop()
				}
			}
		default: // pass
		}
	}
}

func tickGameForever(tick chan bool) {
	for {
		time.Sleep(1 * time.Second)
		tick <- true
	}
}

func handleKeypress(eventKey *tcell.EventKey, quit func(), hoveringTetromino *UnplacedTetromino, board *Board) (placed bool) {
	specialKeys := map[tcell.Key]string{
		tcell.KeyLeft:   "left",
		tcell.KeyRight:  "right",
		tcell.KeyUp:     "up",
		tcell.KeyDown:   "down",
		tcell.KeyEscape: "quit",
		tcell.KeyCtrlC:  "quit",
	}

	runeKeys := map[rune]string{
		rune(' '): "harddrop",
		rune('w'): "up",
		rune('a'): "left",
		rune('s'): "down",
		rune('d'): "right",
	}

	action, actionFound := specialKeys[eventKey.Key()]

	if !actionFound {
		action = runeKeys[eventKey.Rune()]
	}

	switch action {
	case "quit":
		quit()
	case "left":
		hoveringTetromino.Translate(-1, 0, BoardSizeW, BoardSizeH)
	case "right":
		hoveringTetromino.Translate(1, 0, BoardSizeW, BoardSizeH)
	case "up":
		hoveringTetromino.Translate(0, -1, BoardSizeW, BoardSizeH)
	case "down":
		hoveringTetromino.Translate(0, 1, BoardSizeW, BoardSizeH)
	case "harddrop":
		return board.PlaceTetromino(hoveringTetromino)
	}
	return false
}

func drawText(s tcell.Screen, x1, y1 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if r == rune('\n') {
			row++
			col = x1
		}
	}
}
