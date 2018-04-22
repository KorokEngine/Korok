package sine

import (
	"golang.org/x/mobile/exp/audio/al"
)

// Data Model:
// Static & Source


type SourceState uint8

const (
	READY SourceState = iota
	PLAYING
	STOP
	PAUSED
)

type BufferType uint8
const (
// TODO
)

type SourceType uint8
const (
	Static SourceType = iota
	Stream
)

type BufferAL struct {
	al.Buffer
}

func (buf *BufferAL) Create(format uint32, data []byte, freq int32) error {
	array := al.GenBuffers(1)
	buf.Buffer = array[0]

	b := buf.Buffer
	b.BufferData(format, data, freq)

	return nil
}

func (buf *BufferAL) CreateEmpty() {
	array := al.GenBuffers(1)
	buf.Buffer = array[0]
}

func (buf *BufferAL) Destroy() {
	al.DeleteBuffers(buf.Buffer)
}

type StreamBuffer struct {
	Buffer [MaxStreamBuffer]al.Buffer
}

func (stream *StreamBuffer) Create() {
	array := al.GenBuffers(MaxStreamBuffer)
	copy(stream.Buffer[:], array)
}

func (stream *StreamBuffer) Destroy() {
	al.DeleteBuffers(stream.Buffer[:]...)
}
