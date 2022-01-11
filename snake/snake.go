package snake

import (
	"errors"
	"log"
)

type Snake struct {
	Body []Coordinates
}

func (s *Snake) canMove(board *Board, direction int) bool {
	nextHeadPosition, err := s.nextHeadPosition(direction)

	// If current body contains next head position, then return false.
	for _, position := range s.Body {
		if nextHeadPosition == position {
			return false
		}
	}

	if err != nil {
		log.Fatalln(err.Error())
		return false
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

func (s *Snake) nextHeadPosition(direction int) (Coordinates, error) {
	var head Coordinates
	var err error

	switch direction {
	case Up:
		head = NewCoordinates(s.Body[0].x, s.Body[0].y-1)
	case Left:
		head = NewCoordinates(s.Body[0].x-1, s.Body[0].y)
	case Right:
		head = NewCoordinates(s.Body[0].x+1, s.Body[0].y)
	case Down:
		head = NewCoordinates(s.Body[0].x, s.Body[0].y+1)
	default:
		err = errors.New("error: invalid direction")
	}

	return head, err
}

func (s *Snake) Move(direction int) {
	newBody := make([]Coordinates, 0)

	for i := 0; i < len(s.Body); i++ {
		var c Coordinates
		var err error
		if i == 0 {
			c, err = s.nextHeadPosition(direction)

			if err != nil {
				log.Fatalln(err.Error())
				return
			}
		} else {
			c = NewCoordinates(s.Body[i-1].x, s.Body[i-1].y)
		}

		newBody = append(newBody, c)
	}

	s.Body = newBody
}

func NewSnake() *Snake {
	body := make([]Coordinates, 0)
	body = append(body, NewCoordinates(10, 7))
	body = append(body, NewCoordinates(10, 8))
	body = append(body, NewCoordinates(10, 9))
	body = append(body, NewCoordinates(10, 10))
	body = append(body, NewCoordinates(9, 10))

	return &Snake{Body: body}
}
