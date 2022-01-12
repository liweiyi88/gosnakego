package main

import (
	"github.com/gdamore/tcell/v2"
	"gosnakego/snake"
	"os"
)

func eventLoop(game *snake.Game, directionChan chan int) {
	defer close(directionChan)

	for {
		switch event := game.Screen.PollEvent().(type) {
		case *tcell.EventResize:
			game.Screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				game.Screen.Fini()
				os.Exit(0)
			}

			if !game.State.IsStart && event.Key() == tcell.KeyEnter {
				game.Start()
			}

			if !game.State.IsOver {
				if event.Key() == tcell.KeyLeft {
					directionChan <- snake.Left
				}

				if event.Key() == tcell.KeyRight {
					directionChan <- snake.Right
				}

				if event.Key() == tcell.KeyDown {
					directionChan <- snake.Down
				}

				if event.Key() == tcell.KeyUp {
					directionChan <- snake.Up
				}
			}
		}
	}
}

func main() {
	directionChan := make(chan int, 10)
	game := snake.NewGame(snake.NewBoard(50, 20))

	go game.Run(directionChan)
	eventLoop(game, directionChan)
}
