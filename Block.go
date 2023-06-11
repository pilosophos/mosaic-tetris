package main

type Block struct {
	Color string
	Tile  string
}

func NewBlock(color string, tile string) *Block {
	return &Block{color, tile}
}

func (block Block) String() string {
	return block.Color + block.Tile + "\033[0m"
}
