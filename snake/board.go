package snake

type Board struct {
	width  int
	height int
}

// Create a new board.
func NewBoard(width int, height int) *Board {
	return &Board{width, height}
}

// Transform width and height to two-dimensional slice.
func (b *Board) ToCoordinates() []Coordinate {
	var coordinates []Coordinate
	for i := 1; i < b.width; i++ {
		for j := 1; j < b.height; j++ {
			coordinates = append(coordinates, NewCoordinate(i, j))
		}
	}

	return coordinates
}
