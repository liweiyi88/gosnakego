package snake

type Board struct {
	width  int
	height int
}

func NewBoard(width int, height int) *Board {
	return &Board{width, height}
}

// ToCoordinates transform width and height to two-dimensional slice.
func (b *Board) ToCoordinates() []Coordinate {
	var coordinates []Coordinate
	for i := 1; i < b.width; i++ {
		for j := 1; j < b.height; j++ {
			coordinates = append(coordinates, NewCoordinates(i, j))
		}
	}

	return coordinates
}
