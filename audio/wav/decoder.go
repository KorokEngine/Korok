package wav

import (
	"encoding/binary"
	"io"
	"io/ioutil"

	"korok.io/korok/asset/res"
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

// Read returns raw wav data from an input reader
func decode(r io.Reader) (*header, error) {
	h := &header{}
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
	return h, nil
}

/**
impl:
type Decoder interface {
	FullDecode() (d []byte, numChan, freq int32, err error)

	Decode() int
	NumOfChan() int
	BitDepth() int
	SampleRate() int32
	Static() []byte
	ReachEnd() bool
}
*/
type Decoder struct {
	numChannels int32
	sampleRate  int32
	bitDepth    int32

	buffer   []byte
	size     int32
	offset   int32
	reachEnd bool

	file res.File
	name string
}

// DON'T change decoder state! pure-virtual function
func (*Decoder) FullDecode(file res.File) (data []byte, numChan, bitDepth, freq int32, err error) {
	h, err := decode(file)
	defer file.Close()

	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	numChan = int32(h.NumChannels)
	freq = int32(h.SampleRate)
	bitDepth = int32(h.BitsPerSample)
	data = buf
	return
}

// streamed from disc
func (d *Decoder) Decode() (decoded int) {
	n, err := io.ReadFull(d.file, d.buffer)
	if err == io.EOF {
		d.reachEnd = true
	}
	return n
}

func (d *Decoder) head() error {
	if d.file != nil {
		d.file.Close()
	}

	file, err := res.Open(d.name)
	if err != nil {
		return err
	}
	d.file = file
	h, err := decode(file)
	if err != nil {
		return err
	}
	d.numChannels = int32(h.NumChannels)
	d.sampleRate = int32(h.SampleRate)
	d.bitDepth = int32(h.BitsPerSample)
	d.buffer = make([]byte, 16384)
	d.reachEnd = false
	return nil
}

func (d *Decoder) NumOfChan() int32 {
	return d.numChannels
}

func (d *Decoder) BitDepth() int32 {
	return d.bitDepth
}

func (d *Decoder) SampleRate() int32 {
	return d.sampleRate
}

func (d *Decoder) Buffer() []byte {
	return d.buffer
}

func (d *Decoder) ReachEnd() bool {
	return d.reachEnd
}

func (d *Decoder) Rewind() {
	d.head()
}

func (d *Decoder) Close() {
	d.file.Close()
}

func NewDecoder(name string) (d *Decoder, err error) {
	d = new(Decoder)
	d.name = name
	err = d.head()
	return
}
