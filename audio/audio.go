package audio

import (
	"korok/audio/ap"
)

// 组件设计
type SourceComp struct {
	// 音频资源
	Id uint16
	// 优先级
	P  uint16

	// 静音
	Mute bool

	// 循环
	Loop bool

	// 音量 [0, 1]
	Volume float32

	// system
	as *AudioSystem
}

func (s *SourceComp) Play() {

}

// 系统设计
type AudioSystem struct {
	comps []SourceComp

	// player
	sPlayer SamplerPlayer
	mPlayer MusicPlayer
}

func NewAudioSystem() (*AudioSystem, error) {
	ap.Init()

	sys := &AudioSystem{}
	g_sPlayer = &sys.sPlayer
	g_mPlayer = &sys.mPlayer

	return sys, nil
}

// 必须通过更新方法，来检测音频的状态
func (sys *AudioSystem) Update(dt float32) {
	ap.NextFrame()
}

//////////////// static & global field

// 便捷的方法，无视优先级系统直接进入播放
func Play(id uint16) (ok bool){
	// TODO!
	return
}

// default Audio-File-Decoder
var DefaultDecoderFactory = &HaDecoderFactory{}

// shared player
var g_sPlayer *SamplerPlayer
var g_mPlayer *MusicPlayer

func init() {
}



