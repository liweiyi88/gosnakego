package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/liweiyi88/gosnakego/snake"
)

func eventLoop(game *snake.Game, directionChan chan int) {
	defer close(directionChan)

	for {
		switch event := game.PollEvent().(type) {
		case *tcell.EventResize:
			game.Resize()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				game.Exit()
			}

			if !game.HasStarted() && event.Key() == tcell.KeyEnter {
				game.Start()
			}

			if !game.HasEnded() {
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
