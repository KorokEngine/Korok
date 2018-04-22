package ap

import (
	"golang.org/x/mobile/exp/audio/al"
	"golang.org/x/mobile/asset"

	"log"
)

type FileType uint8
const (
	None FileType = iota
	WAV
	VORB

	// NOT IMPLEMENT YET
	OPUS
	MP3
	FLAC
	WMV
)

type FormatEnum uint8
const (
	FormatNone FormatEnum = iota
	Mono8
	Mono16
	Stereo8
	Stereo16
)


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

func (am *AudioManger) LoadSoundBack() {

}

func (am *AudioManger) UnloadSoundBack() {

}

/// 加载数据，得到 Sound 实例
/// 此时应该得出, 采样率，是否Stream等，
func (am *AudioManger) LoadSound(name string, fType FileType, sType SourceType) (id uint16, sound *Sound){
	id, sound = am.indexPool, &am.soundPool[am.indexPool]
	am.indexPool ++
	sound.Type = sType

	d, err := factory.NewDecoder(name, fType)
	if err != nil {
		log.Println("fail to init decoder, ", err)
	}
	if sType == Static {
		file, err := asset.Open(name)
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
		fc = al.FormatMono16
		_, sd := am.allocStaticData(fc, data, freq)
		sound.Data = sd

		log.Println("alloc sound id:", id, " sound:", sound)
	} else {
		_, sd := am.allocStreamData(d)
		sound.Data = sd
	}
	return
}

func (am *AudioManger) allocStaticData(fmt uint32, bits []byte, freq int32) (id uint16, data *StaticData) {
	id, data = am.indexStatic, &am.staticData[am.indexStatic]
	am.indexStatic ++
	data.Static.Create(fmt, bits, freq)
	return
}

func (am *AudioManger) allocStreamData(d Decoder) (id uint16, data *StreamData) {
	id, data = am.indexStream, &am.streamData[am.indexStream]
	am.indexStream ++
	data.Stream.Create()
	data.Decoder = d
	return
}

// TODO!
func (am *AudioManger) freeStatic(id uint16) {
	rear := am.indexStatic - 1
	if id < rear  {

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
	al.FormatMono8,
	al.FormatMono16,
	al.FormatStereo8,
	al.FormatStereo16,
}
