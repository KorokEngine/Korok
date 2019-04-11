// +build js
package sine

import (
	"log"

	"korok.io/korok/webav"
)

const (
	FormatMono8    = 0x1100
	FormatMono16   = 0x1101
	FormatStereo8  = 0x1102
	FormatStereo16 = 0x1103
)

// StaticData is small audio sampler, which will be load into memory directly.
type StaticData struct {
	bits []byte
	fmt  uint32
	freq int32
}

func (d *StaticData) Create(fmt uint32, bits []byte, freq int32) {
	d.fmt = fmt
	d.bits = bits
	d.freq = freq
}

// StreamData will decode pcm-data at runtime. It's used to play big audio files(like .ogg).
type StreamData struct {
	decoder Decoder
}

func (d *StreamData) Create(file string, ft FileType) {
	decoder, err := factory.NewDecoder(file, ft)
	if err != nil {
		log.Println(err)
		return
	}
	d.decoder = decoder
}

type Engine struct {
	engine *webav.AudioContext
}

func (eng *Engine) Initialize() {
	var err error
	eng.engine, err = webav.AudioNewContext(nil)
	if err != nil {
		log.Println("fail to initialize Web Audio API")
	}
}

func (eng *Engine) Destroy() {
	eng.engine.Close()
}

// BufferPlayer can play audio loaded as StaticData.
type BufferPlayer struct {
	source *webav.AudioBufferSource
	status int
}

func (p *BufferPlayer) initialize(engine *Engine) {
	p.source = engine.engine.CreateBufferSource()
	p.status = Stopped
}

func (p *BufferPlayer) destroy() {

}

func (p *BufferPlayer) Play(d *StaticData) {
	p.source.Start()
}

func (p *BufferPlayer) Stop() {
	if p.status == Stopped {
		return
	}
	p.source.Stop()
	p.status = Stopped
}

func (p *BufferPlayer) Pause() {

}

func (p *BufferPlayer) Resume() {

}

// OpenSL's state is different from OpenAL. In OpenAL, if buffer
// queue exhausted, OpenAL will issue a 'Stop' state. But in SL,
// buffer-queue has nothing to do with player state. It' still
// playing even though queue exhausted.

func (p *BufferPlayer) State() int {

	return 0
}

func (p *BufferPlayer) Volume() float32 {
	return 0
}

func (p *BufferPlayer) SetVolume(v float32) {

}

func (p *BufferPlayer) SetLoop(n int) {

}

// StreamPlayer can play audio loaded as StreamData.
type StreamPlayer struct {
	decoder Decoder
	source  *webav.AudioBufferSource
	status  int
}

func (p *StreamPlayer) initialize(engine *Engine) {
	p.source = engine.engine.CreateBufferSource()
	p.status = Stopped
}

func (p *StreamPlayer) Play(d *StreamData) {
	p.decoder = d.decoder
	p.source.Start()
}

func (p *StreamPlayer) Stop() {
	if p.status == Stopped {
		return
	}
	p.source.Stop()
	p.status = Stopped
}

func (p *StreamPlayer) Pause() {

}

func (p *StreamPlayer) Resume() {

}

func (p *StreamPlayer) State() int {
	return 0
}

func (p *StreamPlayer) Volume() float32 {
	return 0
}

func (p *StreamPlayer) SetVolume(v float32) {
}

func (p *StreamPlayer) Tick() {

}

const (
	Initial = 3
	Playing = 2
	Paused  = 1
	Stopped = 0
)
