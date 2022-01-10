package main

import (
	"github.com/gdamore/tcell/v2"
	"gosnakego/snake"
	"os"
)

func main() {
	game := snake.NewGame(snake.NewBoard(50, 20))
	game.Init()

	quit := func() {
		game.End()
		os.Exit(0)
	}

	for {
		game.Run()
		event := game.Screen.PollEvent()

		switch event := event.(type) {
		case *tcell.EventResize:
			game.Screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				quit()
			}

			if event.Key() == tcell.KeyLeft {
				game.MoveSnakeLeft()
			}

			if event.Key() == tcell.KeyRight {
				game.MoveSnakeRight()
			}

			if event.Key() == tcell.KeyDown {
				game.MoveSnakeDown()
			}

			if event.Key() == tcell.KeyUp {
				game.MoveSnakeUp()
			}
		}
	}
}