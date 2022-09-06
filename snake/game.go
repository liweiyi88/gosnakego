package snake

import (
	"fmt"
	"log"
	"math/rand"
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

type State struct {
	sync.Mutex
	isStart   bool
	isOver    bool
	Direction int
	Speed     time.Duration
	Score     int
}

type Game struct {
	sync.Mutex
	State  State
	Apple  *Apple
	Board  *Board
	Snake  *Snake
	Screen tcell.Screen
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
		State:  State{Direction: Up, Speed: time.Millisecond * 100},
		Screen: screen,
	}

	game.setNewApplePosition()
	game.updateScreen()

	return game
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
	if g.State.Direction == direction {
		return false
	}

	if g.State.Direction == Left && direction != Right {
		return true
	}

	if g.State.Direction == Up && direction != Down {
		return true
	}

	if g.State.Direction == Down && direction != Up {
		return true
	}

	if g.State.Direction == Right && direction != Left {
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
	return g.State.isOver
}

func (g *Game) over() {
	g.Lock()
	defer g.Unlock()
	g.State.isOver = true
}

func (g *Game) drawText(x1, y1, x2, y2 int, text string) {
	row := y1
	col := x1
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	for _, r := range text {
		g.Screen.SetContent(col, row, r, nil, style)
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
	g.Screen.SetContent(g.Apple.x, g.Apple.y, '', nil, style)
}

func (g *Game) drawBoard() {
	width, height := g.Board.width, g.Board.height

	boardStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.Screen.SetContent(0, 0, tcell.RuneULCorner, nil, boardStyle)
	for i := 1; i < width; i++ {
		g.Screen.SetContent(i, 0, tcell.RuneHLine, nil, boardStyle)
	}
	g.Screen.SetContent(width, 0, tcell.RuneURCorner, nil, boardStyle)

	for i := 1; i < height; i++ {
		g.Screen.SetContent(0, i, tcell.RuneVLine, nil, boardStyle)
	}

	g.Screen.SetContent(0, height, tcell.RuneLLCorner, nil, boardStyle)

	for i := 1; i < height; i++ {
		g.Screen.SetContent(width, i, tcell.RuneVLine, nil, boardStyle)
	}

	g.Screen.SetContent(width, height, tcell.RuneLRCorner, nil, boardStyle)

	for i := 1; i < width; i++ {
		g.Screen.SetContent(i, height, tcell.RuneHLine, nil, boardStyle)
	}

	g.drawText(1, height+1, width, height+10, fmt.Sprintf("Score:%d", g.State.Score))
	g.drawText(1, height+3, width, height+10, "Press ESC or Ctrl+C to quit")
	g.drawText(1, height+4, width, height+10, "Press arrow keys to control direction")
}

func (g *Game) drawSnake() {
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorGreen)
	for _, coordinates := range g.Snake.Body {
		g.Screen.SetContent(coordinates.x, coordinates.y, tcell.RuneCkBoard, nil, snakeStyle)
	}
}

func (g *Game) move() {
	if g.Snake.canMove(g.Board, g.State.Direction) {
		g.Snake.Move(g.State.Direction)

		if g.Snake.CanEat(g.Apple) {
			g.Snake.Eat(g.Apple)
			g.State.Score++
			g.setNewApplePosition()
		}
	} else {
		g.over()
	}
}

func (g *Game) updateScreen() {
	g.Screen.Clear()
	g.drawLoading()
	g.drawApple()
	g.drawBoard()
	g.drawSnake()
	g.drawEnding()
	g.Screen.Show()
}

func (g *Game) Start() {
	g.State.Lock()
	defer g.State.Unlock()
	g.State.isStart = true
}

func (g *Game) HasStarted() bool {
	g.State.Lock()
	defer g.State.Unlock()
	return g.State.isStart
}

func (g *Game) Run(directionChan chan int) {
	ticker := time.NewTicker(g.State.Speed)
	defer ticker.Stop()

	for {
		select {
		case newDirection := <-directionChan:
			if g.shouldUpdateDirection(newDirection) {
				g.State.Direction = newDirection
			}
		case <-ticker.C:
			if !g.HasEnded() && g.HasStarted() {
				g.move()
			}

			g.updateScreen()
		}
	}
}
