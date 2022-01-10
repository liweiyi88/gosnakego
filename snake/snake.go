package snake

type Snake struct {
	Body []Coordinates
}

func (s *Snake) MoveUp() {
	s.move(Up)
}

func (s *Snake)  MoveDown() {
	s.move(Down)
}

func (s *Snake)  MoveLeft() {
	s.move(Left)
}

func (s *Snake)  MoveRight() {
	s.move(Right)
}

func (s *Snake) move(direction int)  {
	if direction != Up && direction != Down && direction != Left && direction != Right {
		panic("invalid direction")
	}

	newBody := make([]Coordinates, 0)

	for i := 0; i < len(s.Body); i++ {
		var c Coordinates
		if i == 0 {
			switch direction {
			case Up:
				c = NewCoordinates(s.Body[i].x, s.Body[i].y - 1)
			case Left:
				c = NewCoordinates(s.Body[i].x - 1, s.Body[i].y)
			case Right:
				c = NewCoordinates(s.Body[i].x + 1, s.Body[i].y)
			case Down:
				c = NewCoordinates(s.Body[i].x, s.Body[i].y + 1)
			}
		} else {
			c = NewCoordinates(s.Body[i - 1].x, s.Body[i - 1].y)
		}

		newBody = append(newBody, c)
	}

	s.Body = newBody
}

func NewSnake() *Snake{
	body := make([]Coordinates, 0)
	body = append(body, NewCoordinates(10,7))
	body = append(body, NewCoordinates(10,8))
	body = append(body, NewCoordinates(10,9))
	body = append(body, NewCoordinates(10,10))
	body = append(body, NewCoordinates(9,10))

	return &Snake{Body: body}
}
