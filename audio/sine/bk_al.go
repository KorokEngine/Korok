//+build !android

package sine

/*
#cgo darwin   CFLAGS:  -DGOOS_darwin
#cgo linux    CFLAGS:  -DGOOS_linux
#cgo darwin   LDFLAGS: -framework OpenAL
#cgo linux    LDFLAGS: -lopenal


#ifdef GOOS_darwin
#include <stdlib.h>
#include <stdio.h>
#include <OpenAL/al.h>
#include <OpenAL/alc.h>
#endif

typedef struct SineEngine {
	ALCdevice *device;
	ALCcontext *context;
} SineEngine;

void Sine_init(SineEngine *eng) {
	eng->device = alcOpenDevice(0);
	eng->context = alcCreateContext(eng->device, 0);
}

void Sine_wake(SineEngine *eng) {
	alcMakeContextCurrent(eng->context);
	// alcProcessContext(engine.context);
}

void Sine_destroy(SineEngine *eng) {

}

typedef struct SineBufferPlayer {
	ALuint idSource;
} SineBufferPlayer;

void SineBufferPlayer_init(SineBufferPlayer *p) {
	ALuint id;
	alGenSources(1, &id);
	p->idSource = id;
}

void SineBufferPlayer_play(SineBufferPlayer *p, ALuint idBuffer) {
	alSourceQueueBuffers(p->idSource, 1, &idBuffer);
	alSourcePlay(p->idSource);
}

void SineBufferPlayer_pause(SineBufferPlayer *p) {
	alSourcePause(p->idSource);
}

void SineBufferPlayer_stop(SineBufferPlayer *p) {
	alSourceStop(p->idSource);
}

typedef struct SineStreamPlayer {
	ALuint idSource;

} SineStreamPlayer;

void SineStreamPlayer_init(SineStreamPlayer *p) {
	ALuint id;
	alGenSources(1, &id);
	p->idSource = id;
}

void SineStreamPlayer_tick() {

}

// stream play!!

*/
import "C"
import (
	"unsafe"
	"log"
)

// StaticData is small audio sampler, which will be load into memory directly.
type StaticData struct {
	idBuffer C.ALuint
	fmt uint32
	freq int32
}

func (d *StaticData) Create(fmt uint32, bits []byte, freq int32) {
	d.fmt = fmt
	d.freq = freq

	C.alGenBuffers(1, &d.idBuffer)
	C.alBufferData(d.idBuffer, C.ALenum(fmt), unsafe.Pointer(&bits[0]), C.ALsizei(len(bits)), C.ALsizei(freq))
}

// StreamData will decode pcm-data at runtime. It's used to play big audio files(like .ogg).
type StreamData struct {
	decoder Decoder
}

func (d *StreamData) Create(file string, ft FileType) {
	decoder, err := factory.NewDecoder(file, ft)
	if err != nil {
		log.Println(err);return
	}
	d.decoder = decoder
}

type Engine struct {
	engine C.SineEngine
}

func (eng *Engine) Initialize() {
	C.Sine_init(&eng.engine)
	C.Sine_wake(&eng.engine)
}

func (eng *Engine) Destroy() {
	C.Sine_destroy(&eng.engine)
}

// BufferPlayer can play audio loaded as StaticData.
type BufferPlayer struct {
	player        C.SineBufferPlayer
	playingBuffer C.ALuint
}

func (p *BufferPlayer) initialize(engine *Engine) {
	C.SineBufferPlayer_init(&p.player)
}

func (p *BufferPlayer) Play(data *StaticData) {
	C.SineBufferPlayer_play(&p.player, data.idBuffer)
}

func (p *BufferPlayer) Stop() {
	C.SineBufferPlayer_stop(&p.player)
}

func (p *BufferPlayer) Pause() {
	C.SineBufferPlayer_pause(&p.player)
}

func (p *BufferPlayer) Resume() {
	C.SineBufferPlayer_play(&p.player, p.playingBuffer)
}

func (p *BufferPlayer) State() {

}

// StreamPlayer can play audio loaded as StreamData.
type StreamPlayer struct {
	player C.SineStreamPlayer
}

func (p *StreamPlayer) initialize(engine *Engine) {
	C.SineStreamPlayer_init(&p.player)
}

func (p *StreamPlayer) Play(d *StreamData) {

}

func Initialize() {
	//C.hello()
}

const (
	FormatMono8    = 0x1100
	FormatMono16   = 0x1101
	FormatStereo8  = 0x1102
	FormatStereo16 = 0x1103
)

