package snake

type node struct {
	Position
	next *node
}

type Snake struct {
	head *node
}

// Check if the snake can move to the given direction.
func (s *Snake) canMove(board *Board, direction int) bool {
	position := s.nextHeadPosition(direction)

	// If current body contains next head position, then return false.
	if s.contains(position) {
		return false
	}

	switch direction {
	case Up:
		return position.y > 0
	case Left:
		return position.x > 0
	case Right:
		return position.x < board.width
	case Down:
		return position.y < board.height
	}

	return true
}

// The next snake head position in the board.
func (s *Snake) nextHeadPosition(direction int) Position {
	switch direction {
	case Up:
		return newPosition(s.head.x, s.head.y-1)
	case Left:
		return newPosition(s.head.x-1, s.head.y)
	case Right:
		return newPosition(s.head.x+1, s.head.y)
	case Down:
		return newPosition(s.head.x, s.head.y+1)
	default:
		panic("error: invalid direction") // In reality, It shouldn't reach to this line.
	}
}

// Check if the snake has already had the position.
func (s *Snake) contains(position Position) bool {
	current := s.head

	for current != nil {
		if current.x == position.x && current.y == position.y {
			return true
		}

		current = current.next
	}

	return false
}

// Check if the snake can eat the apple.
func (s *Snake) CanEat(apple *Apple) bool {
	return s.head.x == apple.x && s.head.y == apple.y
}

// The snake Eat the apple and add apple's position to snake body slice.
func (s *Snake) Eat(apple *Apple) {
	s.add(apple.Position())
}

// Add new position to the head.
func (s *Snake) add(position Position) *Snake {
	s.head = &node{
		Position: position,
		next:     s.head,
	}

	return s
}

// Move the snake based on the direction.
func (s *Snake) move(direction int) {
	position := s.nextHeadPosition(direction)
	s.add(position)

	if s.head.next.next != nil {
		current := s.head

		for ; current.next.next != nil; current = current.next {
		}

		current.next = nil
	}
}

// Create a new snake with default position and length.
func NewSnake() *Snake {
	snake := &Snake{}

	return snake.add(newPosition(10, 12)).
		add(newPosition(11, 12)).
		add(newPosition(12, 12)).
		add(newPosition(13, 12)).
		add(newPosition(14, 12))
}
