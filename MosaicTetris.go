package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
)

const BoardSizeW = 10
const BoardSizeH = 20

type Highscore struct {
	Score int
	Lines int
	Date  int64
	Name  string
}

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

	var highscorePath string
	if home, err := os.UserHomeDir(); err == nil {
		highscorePath = filepath.Join(home, ".mosaic-tetris-highscore.json")
	}

	// set up TUI
	defStyle := tcell.StyleDefault.Background(tcell.Color16.TrueColor()).Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.Clear()

	termEvents := make(chan tcell.Event)
	tcellQuit := make(chan struct{})
	go s.ChannelEvents(termEvents, tcellQuit)

	// start global tick timer
	tickTimer := make(chan bool)
	go tickGameForever(tickTimer)

	// set up the game
	tetrominoQueue := NewTetrominoQueue()
	hoveringTetromino := tetrominoQueue.Pop()
	board := NewBoard(BoardSizeW, BoardSizeH)

	drawText(s, (BoardSizeW*2)+2, 1, defStyle, "NEXT")
	drawText(s, 0, BoardSizeH+3, defStyle, "Move = WASD/Arrow keys")
	drawText(s, 0, BoardSizeH+4, defStyle, "Hard drop = Space")
	drawText(s, 0, BoardSizeH+5, defStyle, "Quit = Esc/Ctrl+C/q")
	drawText(s, 0, BoardSizeH+7, defStyle, "HOW TO PLAY:")
	drawText(s, 0, BoardSizeH+8, defStyle, "Tetris pieces come randomly rotated in the center")
	drawText(s, 0, BoardSizeH+9, defStyle, "You can't rotate them, but you can put them anywhere and they don't fall")
	drawText(s, 0, BoardSizeH+10, defStyle, "Clear horizontal (or vertical) lines for more points")

	for {
		board.HoverTetromino(hoveringTetromino)

		s.Show()
		board.Render(s, defStyle, drawText, 0, 0)
		drawText(s, (BoardSizeW*2)+2, 2, defStyle, tetrominoQueue.Peek().String())
		drawText(s, (BoardSizeW*2)+2, BoardSizeH/2, defStyle, board.Message)

		select {
		case <-tickTimer:
			timeleft := hoveringTetromino.Tick()
			if timeleft == 0 {
				tetrominoPlaced := board.PlaceTetromino(hoveringTetromino)
				s.Beep()
				if !tetrominoPlaced {
					gameOver(board.Score, board.LinesCleared, highscorePath, s, quit)
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
					s.Beep()
					hoveringTetromino = tetrominoQueue.Pop()
				}
			}
		default: // pass
		}
	}
}

// Check if the current score is a high score, and if so, save the file
// Returns true if a new high score was saved
func saveHighscore(score, lines int, highscorePath string) bool {
	shouldWriteHighscore := false
	var highscore Highscore

	if highscorePath != "" {
		if _, err := os.Stat(highscorePath); errors.Is(err, os.ErrNotExist) {
			shouldWriteHighscore = true
		} else if content, err := os.ReadFile(highscorePath); err == nil {
			if err := json.Unmarshal(content, &highscore); err == nil {
				shouldWriteHighscore = score > highscore.Score
			}
		}
	}

	if shouldWriteHighscore {
		highscore.Score = score
		highscore.Lines = lines
		highscore.Date = time.Now().Unix()

		if currentUser, err := user.Current(); err == nil {
			highscore.Name = currentUser.Username
		} else {
			highscore.Name = "Player"
		}

		if fileData, err := json.Marshal(highscore); err == nil {
			if f, err := os.Create(highscorePath); err == nil {
				f.Write(fileData)
				return true
			}
		}
	}

	return false
}

// Enter the game over state, writing the highscore if needed and prompting the user for a quit
func gameOver(score, lines int, highscorePath string, s tcell.Screen, quit func()) {
	highscoreSaved := saveHighscore(score, lines, highscorePath)
	if highscoreSaved {
		fmt.Println("New highscore!")
	} else {
		fmt.Println("You lose!")
	}
	fmt.Println("Press q to quit!")
	waitForQuit(s, quit)
}

// Run the global tick timer
func tickGameForever(tick chan bool) {
	for {
		time.Sleep(1 * time.Second)
		tick <- true
	}
}

// Hang the program until a quit key is pressed
func waitForQuit(s tcell.Screen, quit func()) {
	for {
		ev := s.PollEvent()
		if ev, ok := ev.(*tcell.EventKey); ok {
			if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEsc || ev.Rune() == 'q' {
				quit()
			}
		}
	}
}

// Handle keypresses and perform game actions accordingly
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
		rune('q'): "quit",
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

// Draw text at the specified position
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
