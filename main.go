package main

import (
	"flag"

	"github.com/liweiyi88/gosnakego/snake"
)

func main() {
	silent := flag.Bool("silent", false, "do not play sound")
	flag.Parse()
	snake.StartGame(*silent)
}
