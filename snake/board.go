package snake

type Board struct {
	width, height int
	area          []Coordinate
}

// Create a new board.
func newBoard(width int, height int) *Board {
	var area []Coordinate
	for i := 1; i < width; i++ {
		for j := 1; j < height; j++ {
			area = append(area, newCoordinate(i, j))
		}
	}

	return &Board{width, height, area}
}
