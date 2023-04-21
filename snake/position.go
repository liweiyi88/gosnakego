package snake

type Position struct {
	x, y int
}

// Create a new position.
func newPosition(x, y int) Position {
	return Position{x, y}
}
