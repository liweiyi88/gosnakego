package snake

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/liweiyi88/gosnakego/assets"
)

const fileGameOver = "gameOver.mp3"
const fileHiss = "hiss.mp3"

type Sound struct {
	gameOver *beep.Buffer
	hiss     *beep.Buffer
}

func NewSound() (Sound, error) {
	hiss, _, err := load(fileHiss)
	if err != nil {
		return Sound{}, err
	}
	gameOver, format, err := load(fileGameOver)
	if err != nil {
		return Sound{}, err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	return Sound{
		gameOver: gameOver,
		hiss:     hiss,
	}, nil
}

func load(file string) (*beep.Buffer, beep.Format, error) {
	f, err := assets.Assets.Open(file)
	if err != nil {
		return nil, beep.Format{}, err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, beep.Format{}, err
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	return buffer, format, nil
}

func (s Sound) GameOver() {
	play(s.gameOver)
}

func (s Sound) Hiss() {
	play(s.hiss)
}

func play(buffer *beep.Buffer) {
	streamer := buffer.Streamer(0, buffer.Len())
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
