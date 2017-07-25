package codec

import (
	"encoding/binary"
	"io"
	"github.com/pkg/errors"
)

// Creator
type Decoder interface {
	Decode(r io.Reader) (Audio, error)
}

// A format holds an audio format's name, magic header and how to decode it.
type format struct {
	name, magic string
	decoder Decoder
}

var formats []format

func RegisterFormat(name, magic string, newDecoder Decoder) {
	formats = append(formats, format{name, magic, newDecoder})
}

// Format is a high level representation of the underlying data.
type Format struct {
	// NumChannels is the number of channels contained in the data
	NumChannels int
	// SampleRate is the sampling rate in Hz
	SampleRate int
	// BitDepth is the number of bits of data for each sample
	BitDepth int
	// Endianess indicate how the byte order of underlying bytes
	Endianness binary.ByteOrder
}

// Data Stream
type Stream interface {
	io.Reader
	io.Seeker
	io.Closer
}

// Audio
type Audio interface {
	// return audio file format
	Format() *Format

	// return Data Stream
	Stream() Stream
}

// decode a audio file
func Decode(hint string, r io.Reader) (Audio, error){
	var decoder Decoder
	for _, f := range formats {
		if f.magic == hint {
			decoder = f.decoder
		}
	}

	if decoder == nil {
		return nil, errors.New("unknown format")
	}

	return decoder.Decode(r)
}
