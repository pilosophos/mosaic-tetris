package main

import (
	"math/rand"
	"time"
)

type TetrominoQueue struct {
	Queue []*UnplacedTetromino
	RNG   *rand.Rand
}

// Create a new tetromino queue
func NewTetrominoQueue() *TetrominoQueue {
	return &TetrominoQueue{
		[]*UnplacedTetromino{},
		rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// Peek at the next tetromino without removing it from the queue,
// generating more tetrominos is the queue is empty
func (tq *TetrominoQueue) Peek() (nextTetromino *UnplacedTetromino) {
	if len(tq.Queue) == 0 {
		tq.RefreshQueue()
	}
	return tq.Queue[0]
}

// Pop the next tetromino from the queue,
// generating more tetrominos is the queue is empty
func (tq *TetrominoQueue) Pop() (nextTetromino *UnplacedTetromino) {
	if len(tq.Queue) == 0 {
		tq.RefreshQueue()
	}
	nextTetromino, tq.Queue = tq.Queue[0], tq.Queue[1:]
	return nextTetromino
}

// Generates more tetrominoes to add to the queue
// All 7 tetrominoes are shuffled, then rotated a random amount before going in the queue
func (tq *TetrominoQueue) RefreshQueue() {
	startLoc := [2]int{5, 10}
	startTime := 4

	tetrominoShapes := []*UnplacedTetromino{
		NewUnplacedTetromino([][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, startLoc, startTime, "cyan"),    // I (cyan)
		NewUnplacedTetromino([][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, startLoc, startTime, "white"),   // O (white)
		NewUnplacedTetromino([][2]int{{0, 0}, {1, 0}, {2, 0}, {1, 1}}, startLoc, startTime, "magenta"), // T (magenta)
		NewUnplacedTetromino([][2]int{{1, 0}, {1, 1}, {1, 2}, {0, 2}}, startLoc, startTime, "blue"),    // J (blue)
		NewUnplacedTetromino([][2]int{{0, 0}, {0, 1}, {1, 2}, {0, 2}}, startLoc, startTime, "yellow"),  // L (yellow)
		NewUnplacedTetromino([][2]int{{0, 1}, {1, 1}, {1, 0}, {2, 0}}, startLoc, startTime, "green"),   // S (green)
		NewUnplacedTetromino([][2]int{{0, 0}, {1, 0}, {1, 1}, {2, 1}}, startLoc, startTime, "red"),     // Z (red)
	}

	// shuffle tetromino shapes randomly
	for i := range tetrominoShapes {
		j := tq.RNG.Intn(i + 1)
		tetrominoShapes[i], tetrominoShapes[j] = tetrominoShapes[j], tetrominoShapes[i]
	}

	// rotate them a random amount of times
	for _, tetromino := range tetrominoShapes {
		tetromino.Rotate(90 * tq.RNG.Intn(4))
	}

	tq.Queue = tetrominoShapes
}
