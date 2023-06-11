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

func (tq *TetrominoQueue) Next() (nextTetromino *UnplacedTetromino) {
	if len(tq.Queue) == 0 {
		tq.RefreshQueue()
	}
	nextTetromino, tq.Queue = tq.Queue[0], tq.Queue[1:]
	return nextTetromino
}

func (tq *TetrominoQueue) RefreshQueue() {
	tetrominoShapes := [][][2]int{
		{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, // I
		{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, // O
		{{0, 0}, {1, 0}, {2, 0}, {1, 1}}, // T
		{{1, 0}, {1, 1}, {1, 2}, {0, 2}}, // J
		{{0, 0}, {0, 1}, {1, 2}, {0, 2}}, // L
		{{0, 1}, {1, 1}, {1, 0}, {2, 0}}, // S
		{{0, 0}, {1, 0}, {1, 1}, {2, 1}}, // Z
	}

	// shuffle tetromino shapes randomly
	for i := range tetrominoShapes {
		j := rand.Intn(i + 1)
		tetrominoShapes[i], tetrominoShapes[j] = tetrominoShapes[j], tetrominoShapes[i]
	}

	for _, shape := range tetrominoShapes {
		tetromino := NewUnplacedTetromino(
			shape,
			[2]int{5, 10},
			4,
			"\033[1;34m",
		)
		tq.Queue = append(tq.Queue, tetromino)
	}
}
