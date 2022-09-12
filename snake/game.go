package snake

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	Up = iota
	Left
	Right
	Down
)

type Apple Coordinate

type Game struct {
	sync.Mutex
	direction int
	speed     time.Duration
	isStart   bool
	isOver    bool
	score     int
	Apple     *Apple
	Board     *Board
	Snake     *Snake
	screen    tcell.Screen
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Create a new apple.
func newApple(x int, y int) *Apple {
	return &Apple{x, y}
}

// Resize screen if terminal changed.
func (g *Game) resizeScreen() {
	g.Lock()
	g.screen.Sync()
	g.Unlock()
}

// Exit the game.
func (g *Game) exit() {
	g.Lock()
	g.screen.Fini()
	g.Unlock()
	os.Exit(0)
}

// Set the new apple's position in the board.
func (g *Game) setNewApplePosition() {
	var availableCoordinates []Coordinate

	for _, coordinate := range g.Board.ToCoordinates() {
		if !g.Snake.Contains(coordinate) {
			availableCoordinates = append(availableCoordinates, coordinate)
		}
	}

	applePosition := availableCoordinates[rand.Intn(len(availableCoordinates))]
	g.Apple = newApple(applePosition.x, applePosition.y)
}

// Check if we need to change the direction based on the new direction
func (g *Game) shouldUpdateDirection(direction int) bool {
	if g.direction == direction {
		return false
	}

	if g.direction == Left && direction != Right {
		return true
	}

	if g.direction == Up && direction != Down {
		return true
	}

	if g.direction == Down && direction != Up {
		return true
	}

	if g.direction == Right && direction != Left {
		return true
	}

	return false
}

// Display the loading screen.
func (g *Game) drawLoading() {
	if !g.hasStarted() {
		g.drawText(g.Board.width/2-12, g.Board.height/2, g.Board.width/2+13, g.Board.height/2, "PRESS <ENTER> TO CONTINUE")
	}
}

// Display the ending screen.
func (g *Game) drawEnding() {
	if g.hasEnded() {
		g.drawText(g.Board.width/2-5, g.Board.height/2, g.Board.width/2+10, g.Board.height/2, "Game over")
	}
}

// Check if the game has ended.
func (g *Game) hasEnded() bool {
	g.Lock()
	defer g.Unlock()
	return g.isOver
}

// Update the game's state to over.
func (g *Game) over() {
	g.Lock()
	defer g.Unlock()
	g.isOver = true
}

// Display text in terminal
func (g *Game) drawText(x1, y1, x2, y2 int, text string) {
	row := y1
	col := x1
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	for _, r := range text {
		g.screen.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

// Display the apple in the board.
func (g *Game) drawApple() {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
	g.screen.SetContent(g.Apple.x, g.Apple.y, '', nil, style)
}

// Display the game board.
func (g *Game) drawBoard() {
	width, height := g.Board.width, g.Board.height

	boardStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.screen.SetContent(0, 0, tcell.RuneULCorner, nil, boardStyle)
	for i := 1; i < width; i++ {
		g.screen.SetContent(i, 0, tcell.RuneHLine, nil, boardStyle)
	}
	g.screen.SetContent(width, 0, tcell.RuneURCorner, nil, boardStyle)

	for i := 1; i < height; i++ {
		g.screen.SetContent(0, i, tcell.RuneVLine, nil, boardStyle)
	}

	g.screen.SetContent(0, height, tcell.RuneLLCorner, nil, boardStyle)

	for i := 1; i < height; i++ {
		g.screen.SetContent(width, i, tcell.RuneVLine, nil, boardStyle)
	}

	g.screen.SetContent(width, height, tcell.RuneLRCorner, nil, boardStyle)

	for i := 1; i < width; i++ {
		g.screen.SetContent(i, height, tcell.RuneHLine, nil, boardStyle)
	}

	g.drawText(1, height+1, width, height+10, fmt.Sprintf("Score:%d", g.score))
	g.drawText(1, height+3, width, height+10, "Press ESC or Ctrl+C to quit")
	g.drawText(1, height+4, width, height+10, "Press arrow keys to control direction")
}

// Display the snake
func (g *Game) drawSnake() {
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorGreen)
	for _, coordinates := range *g.Snake {
		g.screen.SetContent(coordinates.x, coordinates.y, tcell.RuneCkBoard, nil, snakeStyle)
	}
}

// Move the snake and set apple in the board.
func (g *Game) move() {
	if g.Snake.canMove(g.Board, g.direction) {
		g.Snake.
		Move(g.direction)

		if g.Snake.CanEat(g.Apple) {
			g.Snake.Eat(g.Apple)
			g.score++
			g.setNewApplePosition()
		}
	} else {
		g.over()
	}
}

// Update the game screen.
func (g *Game) updateScreen() {
	g.screen.Clear()
	g.drawLoading()
	g.drawApple()
	g.drawBoard()
	g.drawSnake()
	g.drawEnding()
	g.screen.Show()
}

// Update game's state to start.
func (g *Game) start() {
	g.Lock()
	defer g.Unlock()
	g.isStart = true
}

// Check if the game has started.
func (g *Game) hasStarted() bool {
	g.Lock()
	defer g.Unlock()
	return g.isStart
}

// Start the game.
func (g *Game) Run(directionChan chan int) {
	ticker := time.NewTicker(g.speed)
	defer ticker.Stop()

	for {
		select {
		case newDirection := <-directionChan:
			if g.shouldUpdateDirection(newDirection) {
				g.Lock()
				g.direction = newDirection
				g.Unlock()
			}
		case <-ticker.C:
			if !g.hasEnded() && g.hasStarted() {
				g.move()
			}

			g.updateScreen()
		}
	}
}

// Update game state based on keyboard events.
func (g *Game) ReactToKeyBoardEvents(directionChan chan int) {
	defer close(directionChan)

	for {
		switch event := g.screen.PollEvent().(type) {
		case *tcell.EventResize:
			g.resizeScreen()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				g.exit()
			}

			if !g.hasStarted() && event.Key() == tcell.KeyEnter {
				g.start()
			}

			if !g.hasEnded() {
				if event.Key() == tcell.KeyLeft {
					directionChan <- Left
				}

				if event.Key() == tcell.KeyRight {
					directionChan <- Right
				}

				if event.Key() == tcell.KeyDown {
					directionChan <- Down
				}

				if event.Key() == tcell.KeyUp {
					directionChan <- Up
				}
			}
		}
	}
}

// Create a new game
func NewGame(board *Board) *Game {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	game := &Game{
		Board:     board,
		Snake:     NewSnake(),
		direction: Up,
		speed:     time.Millisecond * 100,
		screen:    screen,
	}

	game.setNewApplePosition()
	game.updateScreen()

	return game
}
