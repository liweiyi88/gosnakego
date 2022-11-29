package snake

type Coordinate struct {
	x, y int
}

// Create a new coordinate.
func newCoordinate(x, y int) Coordinate {
	return Coordinate{x, y}
}
