package audio

import (
	"korok.io/korok/audio/sine"
)

func Init() (err error) {
	sine.Init(DefaultDecoderFactory);
	return
}

func Destroy() {
	sine.Destroy()
}

func AdvanceFrame() {
	sine.Tick()
}

// play
func Play(id uint16, p uint16) {
	sine.Play(id)
}

// default Audio-File-Decoder
var DefaultDecoderFactory = &decoderFactory{}


