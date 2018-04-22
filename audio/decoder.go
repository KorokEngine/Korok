package audio

import (
	"korok.io/korok/audio/sine"
	"korok.io/korok/audio/wav"
	"korok.io/korok/audio/ogg"

	"fmt"
)

type decoderFactory struct {
}

func (df *decoderFactory) NewDecoder(name string, fileType sine.FileType) (sine.Decoder, error) {
	switch fileType {
	case sine.WAV:
		return NewWavDecoder(name)
	case sine.VORB:
		return NewVorbisDecoder(name)
	}

	return nil, fmt.Errorf("not support file type: %d", fileType)
}

func NewWavDecoder(name string) (sine.Decoder, error) {
	return wav.NewDecoder(name)
}

func NewVorbisDecoder(name string) (sine.Decoder, error) {
	return ogg.NewVorbisDecoder(name)
}
