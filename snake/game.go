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
	mu        sync.Mutex
	direction int
	speed     time.Duration
	isStart   bool
	isOver    bool
	score     int
	Apple     *Apple
	Board     *Board
	Snake     *Snake
	screen    tcell.Screen
	sound     Sound
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
	g.mu.Lock()
	g.screen.Sync()
	g.mu.Unlock()
}

// Exit the game.
func (g *Game) exit() {
	g.mu.Lock()
	g.screen.Fini()
	g.mu.Unlock()
	os.Exit(0)
}

// Set the new apple's position in the board.
func (g *Game) setNewApplePosition() {
	var availableCoordinates []Coordinate

	for _, coordinate := range g.Board.area {
		if !g.Snake.contains(coordinate) {
			availableCoordinates = append(availableCoordinates, coordinate)
		}
	}

	applePosition := availableCoordinates[rand.Intn(len(availableCoordinates))]
	g.Apple = newApple(applePosition.x, applePosition.y)
}

// Check if we need to change the direction based on the new direction.
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

// Determine should the game continue.
func (g *Game) shouldContinue() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return !g.isOver && g.isStart
}

// Check if the game has ended.
func (g *Game) hasEnded() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.isOver
}

// Update the game's state to over.
func (g *Game) over() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.isOver = true
	g.sound.GameOver()
}

// Display text in terminal.
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

// Display the snake.
func (g *Game) drawSnake() {
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorGreen)
	for _, coordinates := range *g.Snake {
		g.screen.SetContent(coordinates.x, coordinates.y, tcell.RuneCkBoard, nil, snakeStyle)
	}
}

// Update game items' (snake and apple) state.
func (g *Game) updateItemState() {
	if g.Snake.canMove(g.Board, g.direction) {
		g.Snake.move(g.direction)

		if g.Snake.CanEat(g.Apple) {
			g.sound.Hiss()
			g.Snake.eat(g.Apple)
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
	g.mu.Lock()
	defer g.mu.Unlock()
	g.isStart = true
}

// Check if the game has started.
func (g *Game) hasStarted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.isStart
}

// Run the game.
func (g *Game) run(directionChan chan int) {
	ticker := time.NewTicker(g.speed)
	defer ticker.Stop()

	for {
		select {
		case newDirection := <-directionChan:
			if g.shouldUpdateDirection(newDirection) {
				g.mu.Lock()
				g.direction = newDirection
				g.mu.Unlock()
			}
		case <-ticker.C:
			if g.shouldContinue() {
				g.updateItemState()
			}
			g.updateScreen()
		}
	}
}

// Update game state based on keyboard events.
func (g *Game) handleKeyBoardEvents(directionChan chan int) {
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

// Create a new game.
func newGame(board *Board, silent bool) *Game {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
	sound, err := NewSound(silent)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	game := &Game{
		Board:     board,
		Snake:     newSnake(),
		direction: Up,
		speed:     time.Millisecond * 100,
		screen:    screen,
		sound:     sound,
	}

	game.setNewApplePosition()
	game.updateScreen()

	return game
}

// Start the snake game.
func StartGame(silent bool) {
	directionChan := make(chan int, 10)
	game := newGame(newBoard(50, 20), silent)

	go game.run(directionChan)
	game.handleKeyBoardEvents(directionChan)
}
