//+build !android,!js,!windows

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

ALenum Sine_init(SineEngine *eng) {
	alGetError(); // clear error code
	eng->device = alcOpenDevice(0);
	eng->context = alcCreateContext(eng->device, 0);

	return alGetError();
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

ALenum SineBufferPlayer_init(SineBufferPlayer *p) {
	ALuint id;
	alGenSources(1, &id);
	p->idSource = id;
	return alGetError();
}

ALenum SineBufferPlayer_play(SineBufferPlayer *p, ALuint idBuffer) {
	ALuint source = p->idSource;
	alSourcei(source, AL_BUFFER, 0); // do need reset?
	alSourcei(source, AL_BUFFER, idBuffer);
	alSourcePlay(source);
	return alGetError();
}

ALenum SineBufferPlayer_pause(SineBufferPlayer *p) {
	alSourcePause(p->idSource);
	return alGetError();
}

ALenum SineBufferPlayer_stop(SineBufferPlayer *p) {
	alSourceStop(p->idSource);
	return alGetError();
}

// ignore error check for state-checking
ALenum SineBufferPlayer_state(SineBufferPlayer *p) {
 	ALenum state;
 	alGetSourcei(p->idSource, AL_SOURCE_STATE, &state);
 	return state;
}

ALfloat SineBufferPlayer_getVolume(SineBufferPlayer *p) {
	ALfloat v;
	alGetSourcef(p->idSource, AL_GAIN, &v);
	return v;
}

void SineBufferPlayer_setVolume(SineBufferPlayer *p, ALfloat v) {
	alSourcef(p->idSource, AL_GAIN, v);
}

typedef struct SineStreamPlayer {
	ALuint idSource;
	ALuint buffers[8]; //max buffer size
	ALuint numBuffers; //buffer in use
	ALuint freed;      //buffer active
} SineStreamPlayer;

ALenum SineStreamPlayer_init(SineStreamPlayer *p, ALuint numBuffers) {
	ALuint id;
	alGenSources(1, &id);
	p->idSource = id;
	alSourcei(id, AL_LOOPING, AL_FALSE); // stop looping

	p->numBuffers = numBuffers;
	p->freed = numBuffers;
	alGenBuffers(numBuffers, &(p->buffers[0]));

	return alGetError();
}

ALenum SineStreamPlayer_play(SineStreamPlayer *p) {
	alSourcePlay(p->idSource);
	return alGetError();
}


ALenum SineStreamPlayer_stop(SineStreamPlayer *p) {
	alSourceStop(p->idSource);
	return alGetError();
}


ALenum SineStreamPlayer_pause(SineStreamPlayer *p) {
	alSourcePause(p->idSource);
	return alGetError();
}

ALuint SineStreamPlayer_freeBuffer(SineStreamPlayer *p) {
	ALuint source = p->idSource;
	ALuint freed  = p->freed;

	ALint bp;
	alGetSourcei(p->idSource, AL_BUFFERS_PROCESSED, &bp);

	if (bp > 0 ){
		alSourceUnqueueBuffers(source, bp, &(p->buffers[freed]));
		freed += bp;
		p->freed = freed;
	}

	return freed;
}

void SineStreamPlayer_feed(SineStreamPlayer *p, void *data, ALsizei size, ALenum format, ALsizei freq) {
	ALuint source = p->idSource;
	ALuint freed  = p->freed;

	ALuint buffer = p->buffers[freed-1];
	alBufferData(buffer, format, data, size, freq);
	alSourceQueueBuffers(source, 1, &buffer);
	p->freed = freed-1;
}

ALenum SineStreamPlayer_state(SineStreamPlayer *p) {
 	ALenum state;
 	alGetSourcei(p->idSource, AL_SOURCE_STATE, &state);
 	return state;
}


ALfloat SineStreamPlayer_getVolume(SineStreamPlayer *p) {
	ALfloat v;
	alGetSourcef(p->idSource, AL_GAIN, &v);
	return v;
}

void SineStreamPlayer_setVolume(SineStreamPlayer *p, ALfloat v) {
	alSourcef(p->idSource, AL_GAIN, v);
}

// stream play!!

*/
import "C"
import (
	"log"
	"unsafe"
)

// StaticData is small audio sampler, which will be load into memory directly.
type StaticData struct {
	idBuffer C.ALuint
	fmt      uint32
	freq     int32
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
		log.Println(err)
		return
	}
	d.decoder = decoder
}

func (d *StreamData) format() FormatEnum {
	return getFormat(d.decoder.NumOfChan(), d.decoder.BitDepth())
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
	if ret := C.SineBufferPlayer_init(&p.player); ret != NoError {
		log.Println("buffer-player init err:", errString(ret))
	}
}

func (p *BufferPlayer) Play(data *StaticData) {
	if ret := C.SineBufferPlayer_play(&p.player, data.idBuffer); ret != NoError {
		log.Println("buffer-player play err:", errString(ret))
	}
}

func (p *BufferPlayer) Stop() {
	if ret := C.SineBufferPlayer_stop(&p.player); ret != NoError {
		log.Println("buffer-player stop err:", errString(ret))
	}
}

func (p *BufferPlayer) Pause() {
	if ret := C.SineBufferPlayer_pause(&p.player); ret != NoError {
		log.Println("buffer-player pause err:", errString(ret))
	}
}

func (p *BufferPlayer) Resume() {
	if ret := C.SineBufferPlayer_play(&p.player, p.playingBuffer); ret != NoError {
		log.Println("buffer-player resume err:", errString(ret))
	}
}

func (p *BufferPlayer) Volume() float32 {
	v := C.SineBufferPlayer_getVolume(&p.player)
	return float32(v)
}

func (p *BufferPlayer) SetVolume(v float32) {
	C.SineBufferPlayer_setVolume(&p.player, C.ALfloat(v))
}

func (p *BufferPlayer) SetLoop(loop int) {

}

func (p *BufferPlayer) State() uint32 {
	st := C.SineBufferPlayer_state(&p.player)
	return uint32(st)
}

// StreamPlayer can play audio loaded as StreamData.
type StreamPlayer struct {
	player C.SineStreamPlayer
	format uint32
	sampleRate int32
	decoder      Decoder
}

func (p *StreamPlayer) initialize(engine *Engine) {
	var (
		numBuffers = C.ALuint(4)
	)
	if ret := C.SineStreamPlayer_init(&p.player, numBuffers); ret != NoError {
		log.Println("stream-player init err:", errString(ret))
	}
}

func (p *StreamPlayer) Play(stream *StreamData) {
	d := stream.decoder
	p.decoder = d
	p.format = formatCodes[stream.format()]
	p.sampleRate = d.SampleRate()
	p.fill()
	// play
	if ret := C.SineStreamPlayer_play(&p.player); ret != NoError {
		log.Println("stream-player play err:", errString(ret))
	}
}

func (p *StreamPlayer) Stop() {
	if ret := C.SineStreamPlayer_stop(&p.player); ret != NoError {
		log.Println("stream-player stop err:", errString(ret))
	}
	if d := p.decoder; d != nil {
		d.Rewind()
	}
	p.decoder = nil
}

func (p *StreamPlayer) Pause() {
	if ret := C.SineStreamPlayer_pause(&p.player); ret != NoError {
		log.Println("stream-player pause err:", errString(ret))
	}
}

func (p *StreamPlayer) Resume() {
	if ret := C.SineStreamPlayer_play(&p.player); ret != NoError {
		log.Println("stream-player resume err:", errString(ret))
	}
}

func (p *StreamPlayer) State() uint32 {
	st := C.SineStreamPlayer_state(&p.player)
	return uint32(st)
}

func (p *StreamPlayer) Volume() float32 {
	v := C.SineStreamPlayer_getVolume(&p.player)
	return float32(v)
}

func (p *StreamPlayer) SetVolume(v float32) {
	C.SineStreamPlayer_setVolume(&p.player, C.ALfloat(v))
}

func (p *StreamPlayer) Tick() {
	p.fill()
}

func (p *StreamPlayer) fill() {
	d := p.decoder
	if d == nil || d.ReachEnd() {
		return // return if no more data
	}

	free := int(C.SineStreamPlayer_freeBuffer(&p.player))
	if free == 0 {
		return // return if no more free buffer
	}

	// feed the free buffer one by one.
	for ;free > 0 ; free-- {
		if n := d.Decode(); n == 0 {
			break
		}
		var (
			buffer = d.Buffer()
			size   = C.ALsizei(len(buffer))
			format = C.ALenum(p.format)
			freq   = C.ALsizei(p.sampleRate)
		)
		C.SineStreamPlayer_feed(&p.player, unsafe.Pointer(&buffer[0]), size, format, freq)
	}
}

func errString(code C.ALenum) string {
	switch code {
	case InvalidName:
		return "invalid name"
	case InvalidEnum:
		return "invalid enum"
	case InvalidValue:
		return "invalid value"
	case InvalidOperation:
		return "invalid operation"
	case OutOfMemory:
		return "out of memory"
	}
	return "unknown"
}

const (
	FormatMono8    = 0x1100
	FormatMono16   = 0x1101
	FormatStereo8  = 0x1102
	FormatStereo16 = 0x1103
)

// AL error
const (
	NoError          = C.AL_NO_ERROR
	InvalidName      = C.AL_INVALID_NAME
	InvalidEnum      = C.AL_INVALID_ENUM
	InvalidValue     = C.AL_INVALID_VALUE
	InvalidOperation = C.AL_INVALID_OPERATION
	OutOfMemory      = C.AL_OUT_OF_MEMORY
)

// AL state
const (
	Initial = 0x1011
	Playing = 0x1012
	Paused  = 0x1013
	Stopped = 0x1014
)