package snake

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

const fileGameOver = "gameOver.mp3"

type Sound struct {
	gameOver *beep.Buffer
}

func NewSound(base string) (Sound, error) {
	f, err := os.Open(base + "/" + fileGameOver)
	if err != nil {
		return Sound{}, err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return Sound{}, err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	return Sound{
		gameOver: buffer,
	}, nil
}

func (s Sound) Play() {
	streamer := s.gameOver.Streamer(0, s.gameOver.Len())
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
