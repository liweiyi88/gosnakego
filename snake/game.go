package snake

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type GameMode int64

const (
	EasyMode GameMode = iota
	NormalMode
	HardMode
)

const minSpeed = 50

const (
	Up = iota
	Left
	Right
	Down
)

type Apple Position

func (apple *Apple) Position() Position {
	return newPosition(apple.x, apple.y)
}

type Game struct {
	mu             sync.Mutex
	direction      int
	isStart        bool
	isOver         bool
	score          int
	canSpeedChange bool
	mode           GameMode
	Apple          *Apple
	Board          *Board
	Snake          *Snake
	screen         tcell.Screen
}

// Create a new apple.
func newApple(x int, y int) *Apple {
	return &Apple{x, y}
}

func (g *Game) speed() time.Duration {
	switch g.mode {
	case EasyMode:
		return time.Millisecond * 500
	case NormalMode:
		return time.Millisecond * 100
	case HardMode:
		return time.Millisecond * time.Duration(math.Max(minSpeed, float64(100-3*g.score)))
	default:
		return time.Millisecond * 100
	}
}

// Should we update speed when in hard mode.
func (g *Game) shouldUpdateSpeed() bool {
	if g.mode != HardMode {
		return false
	}

	if !g.canSpeedChange {
		return false
	}

	if g.score%3 == 0 {
		g.mu.Lock()
		g.canSpeedChange = false
		g.mu.Unlock()
		return true
	}

	return false
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
	var availablePositions []Position

	for _, position := range g.Board.area {
		if !g.Snake.contains(position) {
			availablePositions = append(availablePositions, position)
		}
	}

	applePosition := availablePositions[rand.Intn(len(availablePositions))]
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
		g.drawText(g.Board.width/2-8, g.Board.height/2-4, g.Board.width/2+17, g.Board.height/2-4, "PRESS <ENTER> TO CONTINUE")

		easy := "  Easy"
		if g.mode == EasyMode {
			easy = "> Easy"
		}

		g.drawText(g.Board.width/2-10, g.Board.height/2-2, g.Board.width/2+13, g.Board.height/2-2, easy)

		normal := "  Normal"
		if g.mode == NormalMode {
			normal = "> Normal"
		}

		g.drawText(g.Board.width/2-10, g.Board.height/2-1, g.Board.width/2+13, g.Board.height/2-1, normal)

		hard := "  Hard"
		if g.mode == HardMode {
			hard = "> Hard"
		}
		g.drawText(g.Board.width/2-10, g.Board.height/2, g.Board.width/2+13, g.Board.height/2, hard)
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

	current := g.Snake.head

	for current != nil {
		g.screen.SetContent(current.x, current.y, tcell.RuneCkBoard, nil, snakeStyle)

		current = current.next
	}
}

// Update game items' (snake and apple) state.
func (g *Game) updateItemState() {
	if g.Snake.canMove(g.Board, g.direction) {
		g.Snake.move(g.direction)

		if g.Snake.CanEat(g.Apple) {
			g.Snake.Eat(g.Apple)
			g.score++
			g.canSpeedChange = true
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

func (g *Game) SetMode(mode GameMode) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.mode = mode
}

// Check if the game has started.
func (g *Game) hasStarted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.isStart
}

// Run the game.
func (g *Game) run(directionChan chan int, gameMode chan GameMode) {
	if !g.hasStarted() {
		for newMode := range gameMode {
			g.SetMode(newMode)
			g.updateScreen()
		}
	}

	g.updateScreen()
	ticker := time.NewTicker(g.speed())
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

			if g.shouldUpdateSpeed() {
				ticker.Reset(g.speed())
			}

			g.updateScreen()
		}
	}
}

// Update game state based on keyboard events.
func (g *Game) handleKeyBoardEvents(directionChan chan int, gameMode chan GameMode) {
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
				close(gameMode)
			}

			if !g.hasStarted() {
				if event.Key() == tcell.KeyUp {
					if g.mode == HardMode {
						gameMode <- NormalMode
					} else if g.mode == NormalMode {
						gameMode <- EasyMode
					}
				} else if event.Key() == tcell.KeyDown {
					if g.mode == EasyMode {
						gameMode <- NormalMode
					} else if g.mode == NormalMode {
						gameMode <- HardMode
					}
				}
			}

			if g.hasStarted() && !g.hasEnded() {
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
func newGame(board *Board) *Game {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	game := &Game{
		Board:     board,
		Snake:     NewSnake(),
		direction: Right,
		mode:      NormalMode,
		screen:    screen,
	}

	game.setNewApplePosition()
	game.updateScreen()

	return game
}

// Start the snake game.
func StartGame() {
	directionChan := make(chan int, 10)
	gameMode := make(chan GameMode, 1)
	game := newGame(newBoard(50, 20))

	go game.run(directionChan, gameMode)
	game.handleKeyBoardEvents(directionChan, gameMode)
}
