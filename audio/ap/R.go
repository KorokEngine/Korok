package ap

import (
	"golang.org/x/mobile/exp/audio/al"

	"log"
)

type FileType uint8
const (
	WAV 	FileType = iota
	VORBIS

	// NOT IMPLEMENT YET
	OPUS
	MP3
	FLAC
	WMV
)

type FormatEnum uint8
const (
	Mono8 FormatEnum = iota
	Mono16
	Stereo8
	Stereo16

	FORMAT_END
)


const MAX_SOUND_POOL_SIZE = 128 // 128=96+32
const MAX_STATIC_DATA = 96
const MAX_STREAM_DATA = 32

type AudioManger struct {
	// sound array
	soundPool [MAX_SOUND_POOL_SIZE]Sound

	// data pool
	staticData [MAX_STATIC_DATA]StaticData
	streamData [MAX_STREAM_DATA]StreamData

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

	d, err := g_df.NewDecoder(name, fType)
	if err != nil {
		log.Println("fail to init decoder, ", err)
	}
	if sType == Static {
		data, numChan, bitDepth, freq, err := d.FullDecode()
		if err != nil {
			log.Println("fail to full decode audio data")
			return
		}
		format := getFormat(numChan, bitDepth)
		if format == FORMAT_END {
			log.Println("invalid audio format")
			return
		}

		fc := formatCodes[format]
		fc = al.FormatMono16
		_, sd := am.allocStaticData()
		buff := &sd.Buffer
		buff.Create(fc, data, freq)
		sound.Data = sd

		log.Println("alloc sound id:", id, " sound:", sound)
	} else {
		sound.Data = &StreamData{d}
	}
	return
}

func (am *AudioManger) allocStaticData() (id uint16, data *StaticData) {
	id, data = am.indexStatic, &am.staticData[am.indexStatic]
	am.indexStatic ++
	return
}

func (am *AudioManger) allocStreamData() (id uint16, data *StreamData) {
	id, data = am.indexStream, &am.streamData[am.indexStream]
	am.indexStream ++
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

func (am *AudioManger) Sound(id uint16) (ok bool, sound *Sound) {
	if id >= MAX_SOUND_POOL_SIZE {
		return false, nil
	}
	return false, &am.soundPool[id]
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
		format = FORMAT_END
	}

	log.Println("input params, chann:", channels, " depth:", depth, " format:", format)

	return format
}

///////// static and global field
var formatCodes = []uint32{
	al.FormatMono8,
	al.FormatMono16,
	al.FormatStereo8,
	al.FormatStereo16,
}
