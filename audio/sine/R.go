package sine

import (
	"log"

	"korok.io/korok/asset/res"
)

type FileType uint8

const (
	None FileType = iota
	WAV
	VORB

	// NOT IMPLEMENT YET
	OPUS
	FLAC
)

type SourceType uint8

const (
	Static SourceType = iota
	Stream
)

type FormatEnum uint8

const (
	FormatNone FormatEnum = iota
	Mono8
	Mono16
	Stereo8
	Stereo16
)

// sound represent a audio segment
type Sound struct {
	Type     SourceType
	Priority uint16
	Data     interface{}
}

const MaxSoundPoolSize = 128 // 128=96+32
const MaxStaticData = 96
const MaxStreamData = 32

type AudioManger struct {
	// sound array
	soundPool [MaxSoundPoolSize]Sound

	// data pool
	staticData [MaxStaticData]StaticData
	streamData [MaxStreamData]StreamData

	indexPool   uint16
	indexStatic uint16
	indexStream uint16
	padding     uint16
}

func NewAudioManager() *AudioManger {
	return new(AudioManger)
}

/// 加载数据，得到 Sound 实例
/// 此时应该得出, 采样率，是否Stream等，
func (am *AudioManger) LoadSound(name string, ft FileType, sType SourceType) (id uint16, sound *Sound) {
	id, sound = am.indexPool, &am.soundPool[am.indexPool]
	am.indexPool++
	sound.Type = sType

	var d interface{}
	if sType == Static {
		_, d = am.LoadStatic(name, ft)
	} else {
		_, d = am.LoadStream(name, ft)
	}
	sound.Data = d
	return
}

func (am *AudioManger) LoadStatic(name string, ft FileType) (id uint16, sd *StaticData) {
	d, err := factory.NewDecoder(name, ft)
	if err != nil {
		log.Println(err)
		return
	}
	file, err := res.Open(name)
	if err != nil {
		log.Println(err)
		return
	}
	data, numChan, bitDepth, freq, err := d.FullDecode(file)
	if err != nil {
		log.Println("fail to full decode audio data")
		return
	}
	format := getFormat(numChan, bitDepth)
	if format == FormatNone {
		log.Println("invalid audio format")
		return
	}

	fc := formatCodes[format]
	fc = FormatMono16
	id, sd = am.allocStaticData(fc, data, freq)
	return
}

func (am *AudioManger) LoadStream(name string, ft FileType) (id uint16, data *StreamData) {
	return am.allocStreamData(name, ft)
}

func (am *AudioManger) allocStaticData(fmt uint32, bits []byte, freq int32) (id uint16, data *StaticData) {
	id, data = am.indexStatic, &am.staticData[am.indexStatic]
	am.indexStatic++
	data.Create(fmt, bits, freq)
	return
}

func (am *AudioManger) allocStreamData(name string, ft FileType) (id uint16, data *StreamData) {
	id, data = am.indexStream, &am.streamData[am.indexStream]
	am.indexStream++
	data.Create(name, ft)
	return
}

// TODO!
func (am *AudioManger) freeStatic(id uint16) {
	rear := am.indexStatic - 1
	if id < rear {

	}
}

func (am *AudioManger) UnloadSound(id uint16) {

}

func (am *AudioManger) Sound(id uint16) (sound *Sound, ok bool) {
	if id >= MaxSoundPoolSize {
		return nil, false
	}
	return &am.soundPool[id], true
}

func getFormat(channels, depth int32) FormatEnum {
	var format FormatEnum
	switch {
	case channels == 1 && depth == 8:
		format = Mono8
	case channels == 1 && depth == 16:
		format = Mono16
	case channels == 2 && depth == 8:
		format = Stereo8
	case channels == 2 && depth == 16:
		format = Stereo16
	default:
		format = FormatNone
	}
	return format
}

///////// static and global field
var formatCodes = []uint32{
	0, // none
	FormatMono8,
	FormatMono16,
	FormatStereo8,
	FormatStereo16,
}
