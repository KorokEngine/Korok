package wav

import (
	"io"
	"korok/audio/codec"
	"encoding/binary"
)

const (
	// Data format codes

	// PCM
	wave_FORMAT_PCM = 0x0001

	// IEEE float
	wave_FORMAT_IEEE_FLOAT = 0x0003

	// 8-bit ITU-T G.711 A-law
	wave_FORMAT_ALAW = 0x0006

	// 8-bit ITU-T G.711 µ-law
	wave_FORMAT_MULAW = 0x0007

	// Determined by SubFormat
	wave_FORMAT_EXTENSIBLE = 0xFFFE
)

// header
type header struct {
	bChunkID  [4]byte // BUF
	ChunkSize uint32  // L
	bFormat   [4]byte // BUF

	bSubchunk1ID  [4]byte // BUF
	Subchunk1Size uint32  // L

	AudioFormat   uint16 // L
	NumChannels   uint16 // L
	SampleRate    uint32 // L
	ByteRate      uint32 // L
	BlockAlign    uint16 // L
	BitsPerSample uint16 // L

	bSubchunk2ID  [4]byte // BUF
	Subchunk2Size uint32  // L
	Data          []byte  // L
}

// pcm stream
type stream struct {
	offset int
	reader io.Reader
}

func (s*stream)Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (s*stream)Close() error {
	return nil
}

func (s*stream)Read(b []byte) (int, error) {
	return s.reader.Read(b)
}

type Audio struct {
	header
	stream
}

func (a *Audio) Format() *codec.Format {
	return &codec.Format{
		NumChannels: int(a.header.NumChannels),
		SampleRate: int(a.header.SampleRate),
		BitDepth: int(a.header.BitsPerSample),
	}
}

func (d *Audio) Stream() codec.Stream {
	return &d.stream
}

// impl Decode(r io.Reader) (codec.Audio, error)
type Decoder struct {
}

// Read returns raw wav data from an input reader
func (*Decoder)Decode(r io.Reader) (codec.Audio, error) {
	h := header{}

	// header
	err := binary.Read(r, binary.BigEndian, &h.bChunkID)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.ChunkSize)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.BigEndian, &h.bFormat)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &h.bSubchunk1ID)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.Subchunk1Size)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.AudioFormat)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.NumChannels)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.SampleRate)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.ByteRate)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.BlockAlign)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.BitsPerSample)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &h.bSubchunk2ID)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r, binary.LittleEndian, &h.Subchunk2Size)
	if err != nil {
		return nil, err
	}

	// pcm stream
	s := stream{
		offset:44,
		reader:r,
	}

	// return audio
	return &Audio{h, s}, nil
}


func init() {
	/// Register to codec.formats
	codec.RegisterFormat("wav", "wav", &Decoder{})
}


