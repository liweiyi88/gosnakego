package snake

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"time"
)

const (
	Up = iota
	Left
	Right
	Down
)

type State struct {
	IsStart   bool
	IsOver    bool
	Direction int
	Speed time.Duration
}

type Game struct {
	State  State
	Board  *Board
	Snake  *Snake
	Screen tcell.Screen
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
		State:  State{Direction: Up, Speed: time.Millisecond * 200},
		Screen: screen,
	}

	game.updateScreen()

	return game
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
	if !g.State.IsStart {
		g.drawText(g.Board.width / 2 - 12, g.Board.height / 2, g.Board.width / 2 + 13, g.Board.height / 2, "PRESS <ENTER> TO CONTINUE")
	}
}

func (g *Game) drawEnding() {
	if g.State.IsOver {
		g.drawText(g.Board.width / 2 - 5, g.Board.height / 2, g.Board.width / 2 + 10, g.Board.height / 2, "Game over")
	}
}

func (g *Game) over() {
	g.State.IsOver = true
}

func (g *Game) drawText(x1, y1, x2, y2 int, text string) {
	row := y1
	col := x1
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	for _, r := range []rune(text) {
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

	g.drawText(1, height + 1, width, height + 10, fmt.Sprintf("Score:%d", 100))
	g.drawText(1, height + 3, width, height + 10, "Press ESC or Ctrl+C to quit")
	g.drawText(1, height + 4, width, height + 10, "Press arrow keys to control direction")
}

func (g *Game) drawSnake() {
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorGreen)
	for _, coordinates := range g.Snake.Body {
		g.Screen.SetContent(coordinates.x, coordinates.y, tcell.RuneCkBoard, nil, snakeStyle)
	}
}

func (g *Game) Start() {
	g.State.IsStart = true
}

func (g *Game) updateScreen() {
	g.Screen.Clear()
	g.drawLoading()
	g.drawBoard()
	g.drawSnake()
	g.drawEnding()
	g.Screen.Show()
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
			if !g.State.IsOver && g.State.IsStart {
				g.move()
			}

			g.updateScreen()
		}
	}
}

func (g *Game) move() {
	if g.Snake.canMove(g.Board, g.State.Direction) {
		g.Snake.Move(g.State.Direction)
	} else {
		g.over()
	}
}
