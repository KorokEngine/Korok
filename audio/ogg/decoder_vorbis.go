package ogg

import (
	"io/ioutil"
	"os"

	"vorbis"
	"unsafe"
	"log"
)

/// 此处是双通道的代码示例！！
//func Read(b []byte) (int, error) {
//	buf := make([]float32, len(b)/2)
//	n, err := s.reader.Read(buf)
//	if err != nil {
//		return -1, err
//	}
//	for i := 0; i < n; i++ {
//		v := int16(buf[i] * (1<<15 - 1))
//		b[i*2 + 0] = uint8(v)
//		b[i*2 + 1] = uint8(v>>8)
//	}
//	return n*2, nil
//}

/** impl:
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
	numChannels   int32
	sampleRate    int32
	bitDepth      int32

	buffer []byte
	name string
	vorb *vorbis.Vorbis
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
	return false
}

func (*Decoder) FullDecode(file *os.File) (data []byte, numChan, bitDepth, freq int32, err error) {
	b, err := ioutil.ReadAll(file)
	defer file.Close()

	if err != nil {
		return
	}

	_data, _numChan, _freq, err := vorbis.Decode(b)

	if err != nil {
		return
	}
	ptr := unsafe.Pointer(&_data[0])
	data = ((*[1<<20]byte)(ptr))[:len(_data)*2]
	numChan = int32(_numChan)
	bitDepth = 16
	freq = int32(_freq)
	return
}

func (d *Decoder) Decode() int {
	if d.vorb == nil {
		f, err := os.Open(d.name)
		if err != nil {
			return 0
		}
		v, err := vorbis.New(f)
		if err != nil {
			return 0
		}
		d.vorb = v
		d.numChannels = int32(v.Channels)
		d.sampleRate  = int32(v.SampleRate)
		d.bitDepth = 16
	}

	data, err := d.vorb.Decode()
	if err != nil {
		return 0
	}
	log.Println("decode data:", len(data))

	ptr := unsafe.Pointer(&data[0])
	d.buffer = ((*[1<<20]byte)(ptr))[:len(data)*4]
	return len(data) * 4
}

func (d *Decoder) Close() {
	if v := d.vorb; v != nil {
		v.Close()
	}
}

func NewVorbisDecoder(name string) (d *Decoder, err error){
	d = new(Decoder)
	d.name = name
	return
}
