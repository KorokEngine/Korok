package audio

import (
	"korok.io/korok/audio/ap"
	"korok.io/korok/audio/wav"
	"korok.io/korok/audio/ogg"

	"fmt"
)

type HaDecoderFactory struct {
}

func (df *HaDecoderFactory) NewDecoder(name string, fileType ap.FileType) (ap.Decoder, error) {
	switch fileType {
	case ap.WAV:
		return NewWavDecoder(name)
	case ap.VORB:
		return NewVorbisDecoder(name)
	}

	return nil, fmt.Errorf("not support file type: %d", fileType)
}

func NewWavDecoder(name string) (ap.Decoder, error) {
	return wav.NewDecoder(name)
}

func NewVorbisDecoder(name string) (ap.Decoder, error) {
	return ogg.NewVorbisDecoder(name)
}
