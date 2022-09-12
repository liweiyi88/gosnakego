package snake

import (
	"errors"
	"log"
)

// A snake is made of a slice of coordinate.
type Snake []Coordinate

// Check if the snake can move to the given direction.
func (s *Snake) canMove(board *Board, direction int) bool {
	nextHeadPosition, err := s.nextHeadPosition(direction)

	// If current body contains next head position, then return false.
	for _, position := range *s {
		if nextHeadPosition == position {
			return false
		}
	}

	if err != nil {
		log.Fatalln(err.Error())
	}

	switch direction {
	case Up:
		return nextHeadPosition.y > 0
	case Left:
		return nextHeadPosition.x > 0
	case Right:
		return nextHeadPosition.x < board.width
	case Down:
		return nextHeadPosition.y < board.height
	}

	return true
}

// The next snake head position in the board.
func (s *Snake) nextHeadPosition(direction int) (Coordinate, error) {
	var head Coordinate
	var err error

	switch direction {
	case Up:
		head = NewCoordinate((*s)[0].x, (*s)[0].y-1)
	case Left:
		head = NewCoordinate((*s)[0].x-1, (*s)[0].y)
	case Right:
		head = NewCoordinate((*s)[0].x+1, (*s)[0].y)
	case Down:
		head = NewCoordinate((*s)[0].x, (*s)[0].y+1)
	default:
		err = errors.New("error: invalid direction") // In reality, It shouldn't reach to this line.
	}

	return head, err
}

// Check if the snake has already had the coordinate.
func (s *Snake) Contains(coordinate Coordinate) bool {
	for _, body := range *s {
		if coordinate == body {
			return true
		}
	}

	return false
}

// Check if the snake can eat the apple.
func (s *Snake) CanEat(apple *Apple) bool {
	headPosition := (*s)[0]

	return headPosition.x == apple.x && headPosition.y == apple.y
}

// The snake eat the apple and add apple's coordinate to snake body slice.
func (s *Snake) Eat(apple *Apple) {
	coordinate := NewCoordinate(apple.x, apple.y)
	(*s) = append([]Coordinate{coordinate}, (*s)...)
}

// Move the snake based on the direction.
func (s *Snake) Move(direction int) {
	newBody := make([]Coordinate, 0)

	for i := 0; i < len((*s)); i++ {
		var coordinates Coordinate
		var err error
		if i == 0 {
			coordinates, err = s.nextHeadPosition(direction)

			if err != nil {
				log.Fatalln(err.Error())
				return
			}
		} else {
			coordinates = NewCoordinate((*s)[i-1].x, (*s)[i-1].y)
		}

		newBody = append(newBody, coordinates)
	}

	*s = newBody
}

// Create a new snake with default position and length.
func NewSnake() *Snake {
	body := make(Snake, 0)
	body = append(body, NewCoordinate(10, 7))
	body = append(body, NewCoordinate(10, 8))
	body = append(body, NewCoordinate(10, 9))
	body = append(body, NewCoordinate(10, 10))
	body = append(body, NewCoordinate(9, 10))

	return &body
}
