package main

import (
	"log"
	"os"

	"github.com/liweiyi88/gosnakego/snake"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("USAGES: <PATH_TO_ASSETS_DIRECTORY>")
	}
	snake.StartGame(os.Args[1])
}
