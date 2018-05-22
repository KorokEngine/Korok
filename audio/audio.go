package audio

import (
	"korok.io/korok/audio/sine"
)

type MusicPlayer struct {
	*sine.StreamPlayer
}

type ChanId int

var (
	music *sine.StreamPlayer
	effects *sine.SoundPool
)

func Init() (err error) {
	sine.Init(DefaultDecoderFactory)
	music   = sine.NewStreamPlayer()
	effects = sine.NewSoundPool()
	return
}

func Destroy() {
	sine.Destroy()
}

func AdvanceFrame() {
	music.Tick()
	effects.Tick()
}

////////////////////// Music ////////////////////

func PlayMusic(id uint16) (sp MusicPlayer, ook bool){
	if sound, ok := sine.R.Sound(id); ok {
		if d, ok := sound.Data.(*sine.StreamData); ok {
			music.Play(d)
			sp, ook = MusicPlayer{music}, true
		}
	}
	return
}

func PauseMusic() {
	music.Pause()
}

func ResumeMusic() {
	music.Resume()
}

func StopMusic() {
	music.Stop()
}

func SetMusicVolume(v float32) {
	music.SetVolume(v)
}

func MusicVolume() float32 {
	return music.Volume()
}

////////////////////// Effect ////////////////////

func PlayEffect(id uint16, priority int) (cid ChanId){
	return ChanId(effects.Play(id, priority))
}

func PauseEffect(cid ChanId) {
	effects.PauseChan(int(cid))
}

func ResumeEffect(cid ChanId) {
	effects.ResumeChan(int(cid))
}

func StopEffect(cid ChanId) {
	effects.StopChan(int(cid))
}

func EffectChannelVolume(cid ChanId) (v float32, ok bool){
	return effects.GetChanVolume(int(cid))
}

func SetEffectChannelVolume(cid ChanId, v float32) {
	effects.SetChanVolume(int(cid), v)
}

// Overall volume setting.
func EffectVolume() float32 {
	return effects.Volume()
}

func SetEffectVolume(v float32) {
	effects.SetVolume(v)
}


// default Audio-File-Decoder
var DefaultDecoderFactory = &decoderFactory{}


