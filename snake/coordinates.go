package snake

type Coordinate struct {
	x, y int
}

func NewCoordinates(x, y int) Coordinate {
	return Coordinate{x, y}
}
