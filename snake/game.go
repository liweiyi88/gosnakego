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

type Apple struct {
	Coordinate
}

type Game struct {
	sync.Mutex
	direction int
	speed time.Duration
	isStart   bool
	isOver    bool
	score int
	Apple  *Apple
	Board  *Board
	Snake  *Snake
	screen tcell.Screen
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newApple(x int, y int) *Apple {
	return &Apple{Coordinate{
		x,
		y,
	}}
}

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
		Board:  board,
		Snake:  NewSnake(),
		direction: Up,
		speed: time.Millisecond * 100,
		screen: screen,
	}

	game.setNewApplePosition()
	game.updateScreen()

	return game
}

func (g *Game) PollEvent() tcell.Event {
	return g.screen.PollEvent()
}

func (g *Game) Resize() {
	g.Lock()
	g.screen.Sync()
	g.Unlock()
}

func (g *Game) Exit() {
	g.Lock()
	g.screen.Fini()
	g.Unlock()
	os.Exit(0)
}

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

func (g *Game) drawLoading() {
	if !g.HasStarted() {
		g.drawText(g.Board.width/2-12, g.Board.height/2, g.Board.width/2+13, g.Board.height/2, "PRESS <ENTER> TO CONTINUE")
	}
}

func (g *Game) drawEnding() {
	if g.HasEnded() {
		g.drawText(g.Board.width/2-5, g.Board.height/2, g.Board.width/2+10, g.Board.height/2, "Game over")
	}
}

func (g *Game) HasEnded() bool {
	g.Lock()
	defer g.Unlock()
	return g.isOver
}

func (g *Game) over() {
	g.Lock()
	defer g.Unlock()
	g.isOver = true
}

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

func (g *Game) drawApple() {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
	g.screen.SetContent(g.Apple.x, g.Apple.y, '', nil, style)
}

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

func (g *Game) drawSnake() {
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorGreen)
	for _, coordinates := range g.Snake.Body {
		g.screen.SetContent(coordinates.x, coordinates.y, tcell.RuneCkBoard, nil, snakeStyle)
	}
}

func (g *Game) move() {
	if g.Snake.canMove(g.Board, g.direction) {
		g.Snake.Move(g.direction)

		if g.Snake.CanEat(g.Apple) {
			g.Snake.Eat(g.Apple)
			g.score++
			g.setNewApplePosition()
		}
	} else {
		g.over()
	}
}

func (g *Game) updateScreen() {
	g.screen.Clear()
	g.drawLoading()
	g.drawApple()
	g.drawBoard()
	g.drawSnake()
	g.drawEnding()
	g.screen.Show()
}

func (g *Game) Start() {
	g.Lock()
	defer g.Unlock()
	g.isStart = true
}

func (g *Game) HasStarted() bool {
	g.Lock()
	defer g.Unlock()
	return g.isStart
}

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
			if !g.HasEnded() && g.HasStarted() {
				g.move()
			}

			g.updateScreen()
		}
	}
}
