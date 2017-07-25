package ogg

import (
	"korok/audio/codec"
	"io"
	"github.com/jfreymuth/oggvorbis"
)

/**
	based on freymuth/oggvorbis
 */
type stream struct {
	reader *oggvorbis.Reader
}

func (s*stream)Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (s*stream)Close() error {
	return nil
}

func (s*stream)Read(b []byte) (int, error) {
	buf := make([]float32, len(b)/2)
	n, err := s.reader.Read(buf)
	if err != nil {
		return -1, err
	}
	for i := 0; i < n; i++ {
		v := int16(buf[i] * (1<<15 - 1))
		b[i*2 + 0] = uint8(v)
		b[i*2 + 1] = uint8(v>>8)
	}
	return n*2, nil
}


type Audio struct {
	f *codec.Format
	s *stream
}

func (a *Audio) Format() *codec.Format {
	return a.f
}

func (a *Audio) Stream() codec.Stream {
	return a.s
}

// impl Decode(r io.Reader) (codec.Audio, error)
type Decoder struct {
}

// Read returns raw wav data from an input reader
func (*Decoder)Decode(r io.Reader) (codec.Audio, error) {
	reader, err := oggvorbis.NewReader(r)
	if err != nil {
		return nil, err
	}

	f := codec.Format{
		NumChannels:reader.Channels(),
		SampleRate:reader.SampleRate(),
		BitDepth:reader.Bitrate().Nominal, // TODO
	}

	s := stream{reader}

	a := &Audio{
		f: &f,
		s: &s,
	}
	return a, nil
}

func init() {
	codec.RegisterFormat("ogg", "vorbis", &Decoder{})
}
