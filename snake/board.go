package snake

type Board struct {
	width  int
	height int
}

func NewBoard(width int, height int) *Board {
	return &Board{width, height}
}
