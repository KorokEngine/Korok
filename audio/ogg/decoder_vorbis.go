package ogg

import (
	"golang.org/x/mobile/asset"
	"github.com/jfreymuth/oggvorbis"

	"unsafe"
	"io"
	"log"
)

type Decoder struct {
	numChannels   int32
	sampleRate    int32
	bitDepth      int32

	i16buffer []int16
	f32buffer []float32
	size int

	name string
	file asset.File
	reader *oggvorbis.Reader
	reachEnd bool
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
	return ((*[1<<31]byte)(unsafe.Pointer(&d.i16buffer[0])))[:d.size*2]
}

func (d *Decoder) ReachEnd() bool {
	return d.reachEnd
}

func (*Decoder) FullDecode(file asset.File) (data []byte, numChan, bitDepth, freq int32, err error) {
	floats, format, err := oggvorbis.ReadAll(file)
	if err != nil {
		return
	}
	defer file.Close()

	log.Println("data size:", len(floats), "format:", format)

	numChan = int32(format.Channels)
	bitDepth = 16
	freq = int32(format.SampleRate)

	i16s := make([]int16, len(floats))
	f216(floats, i16s)
	ptr := unsafe.Pointer(&i16s[0])
	data = ((*[1<<31]byte)(ptr))[:len(floats)*2]
	return
}

func (d *Decoder) Decode() int {
	size, err := d.reader.Read(d.f32buffer)
	if err != nil && err != io.EOF {
		log.Println("vorbis decode err:", err)
		return 0
	}
	if io.EOF == err {
		d.reachEnd = true
	}
	f216(d.f32buffer, d.i16buffer)
	d.size = size
	return size
}

func (d *Decoder) head() error {
	if f := d.file; f != nil {
		f.Close()
	}
	f, err := asset.Open(d.name)
	if err != nil {
		return err
	}
	r, err := oggvorbis.NewReader(f)
	if err != nil {
		return err
	}

	d.file = f
	d.reader = r
	d.numChannels = int32(r.Channels())
	d.sampleRate  = int32(r.SampleRate())
	d.bitDepth = 16
	d.reachEnd = false
	return nil
}

func (d *Decoder) Rewind() {
	d.head()
}

func NewVorbisDecoder(name string) (d *Decoder, err error){
	d = new(Decoder)
	d.name = name
	d.f32buffer = make([]float32, 16384)
	d.i16buffer = make([]int16, 16384)
	err = d.head()
	return
}

func f216(f32 []float32, i16 []int16) {
	for i, f := range f32 {
		i16[i] = int16(f*32767)
	}
}
