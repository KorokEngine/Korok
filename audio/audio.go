package audio

import (
	"korok.io/korok/audio/sine"
	"log"
)

func Init() (err error) {
	if err = sine.Init(); err != nil {
		log.Println("audio:", err)
	} else {
		sine.SetDecoderFactory(DefaultDecoderFactory)
	}
	return
}

func Destroy() {
	sine.Destroy()
}

func AdvanceFrame() {
	sine.NextFrame()
}

// play
func Play(id uint16, p uint16) {
	sine.Play(id, p)
}

// default Audio-File-Decoder
var DefaultDecoderFactory = &decoderFactory{}


