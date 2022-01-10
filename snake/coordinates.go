package snake

type Coordinates struct {
	x int
	y int
}

func NewCoordinates(x, y int) Coordinates {
	return Coordinates{x, y}
}
