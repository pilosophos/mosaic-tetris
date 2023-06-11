package main

type Block struct {
	Color string
	Tile  rune
}

func NewBlock(color string, tile rune) *Block {
	return &Block{color, tile}
}

func (block Block) Rune() rune {
	return block.Tile
}
