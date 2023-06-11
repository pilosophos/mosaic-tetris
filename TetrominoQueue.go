package main

import "math/rand"

type TetrominoQueue struct {
	Queue []*UnplacedTetromino
}

func NewTetrominoQueue() *TetrominoQueue {
	return &TetrominoQueue{
		[]*UnplacedTetromino{},
	}
}

func (tq *TetrominoQueue) Peek() (nextTetromino *UnplacedTetromino) {
	if len(tq.Queue) == 0 {
		tq.RefreshQueue()
	}
	return tq.Queue[0]
}

func (tq *TetrominoQueue) Pop() (nextTetromino *UnplacedTetromino) {
	if len(tq.Queue) == 0 {
		tq.RefreshQueue()
	}
	nextTetromino, tq.Queue = tq.Queue[0], tq.Queue[1:]
	return nextTetromino
}

func (tq *TetrominoQueue) RefreshQueue() {
	startLoc := [2]int{5, 10}
	startTime := 4

	tetrominoShapes := []*UnplacedTetromino{
		NewUnplacedTetromino([][2]int{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, startLoc, startTime, "\033[1;36m"), // I (cyan)
		NewUnplacedTetromino([][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, startLoc, startTime, "\033[1;37m"), // O (white)
		NewUnplacedTetromino([][2]int{{0, 0}, {1, 0}, {2, 0}, {1, 1}}, startLoc, startTime, "\033[1;35m"), // T (magenta)
		NewUnplacedTetromino([][2]int{{1, 0}, {1, 1}, {1, 2}, {0, 2}}, startLoc, startTime, "\033[1;34m"), // J (blue)
		NewUnplacedTetromino([][2]int{{0, 0}, {0, 1}, {1, 2}, {0, 2}}, startLoc, startTime, "\033[1;33m"), // L (yellow)
		NewUnplacedTetromino([][2]int{{0, 1}, {1, 1}, {1, 0}, {2, 0}}, startLoc, startTime, "\033[1;32m"), // S (green)
		NewUnplacedTetromino([][2]int{{0, 0}, {1, 0}, {1, 1}, {2, 1}}, startLoc, startTime, "\033[1;31m"), // Z (red)
	}

	// shuffle tetromino shapes randomly
	for i := range tetrominoShapes {
		j := rand.Intn(i + 1)
		tetrominoShapes[i], tetrominoShapes[j] = tetrominoShapes[j], tetrominoShapes[i]
	}

	// rotate them a random amount of times
	for _, tetromino := range tetrominoShapes {
		tetromino.Rotate(90 * rand.Intn(4))
	}

	tq.Queue = tetrominoShapes
}
