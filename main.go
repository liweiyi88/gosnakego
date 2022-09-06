package main

import (
	"github.com/liweiyi88/gosnakego/snake"
)

func main() {
	directionChan := make(chan int, 10)
	game := snake.NewGame(snake.NewBoard(50, 20))

	go game.Run(directionChan)
	game.ReactToKeyBoardEvents(directionChan)
}
