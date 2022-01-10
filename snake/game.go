package snake

import (
	"github.com/gdamore/tcell/v2"
	"log"
)

const (
	Up = iota
	Left
	Right
	Down
)

type State struct {
	Direction int
}

type Game struct {
	State State
	Board  *Board
	Snake *Snake
	Screen tcell.Screen
}

func NewGame(board *Board) *Game{
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	return &Game{
		Board: board,
		Snake: NewSnake(),
		State: State{Direction: Up},
		Screen: screen,
	}
}

func (g *Game) createScreen() {
	screen, err := tcell.NewScreen()
	g.Screen = screen

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := g.Screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.Screen.SetStyle(defStyle)
}

func (g *Game) Init()  {
	g.createScreen()
}

func (g *Game) Run() {
	g.Screen.Clear()
	g.drawBoard()
	g.drawSnake()
	g.Screen.Show()
}

func (g *Game) Update() {
	g.Screen.Clear()
	g.drawBoard()
	g.drawSnake()
	g.Screen.Show()
}

func (g *Game) MoveSnakeLeft()  {
	g.Snake.MoveLeft()
	g.Update()
}

func (g *Game) MoveSnakeRight()  {
	g.Snake.MoveRight()
	g.Update()
}

func (g *Game) MoveSnakeUp()  {
	g.Snake.MoveUp()
	g.Update()
}

func (g *Game) MoveSnakeDown()  {
	g.Snake.MoveDown()
	g.Update()
}

func (g *Game) drawSnake() {
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorGreen)
	for _, coordinates := range g.Snake.Body {
		g.Screen.SetContent(coordinates.x, coordinates.y, tcell.RuneCkBoard, nil, snakeStyle)
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
}

func (g *Game) End()  {
	g.Screen.Fini()
}
