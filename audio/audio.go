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

func SetMusicVolume(leftVolume, rightVolume float32) {
	music.SetVolume(leftVolume, rightVolume)
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

func SetEffectVolume(cid ChanId, leftVolume, rightVolume float32) {
	effects.SetVolume(int(cid), leftVolume, rightVolume)
}

func SetEffectCallback() {

}




// default Audio-File-Decoder
var DefaultDecoderFactory = &decoderFactory{}


