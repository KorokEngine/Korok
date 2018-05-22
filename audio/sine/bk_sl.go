// +build android

package sine

/*
 #cgo LDFLAGS: -landroid -lOpenSLES

 #include <SLES/OpenSLES.h>
 #include <SLES/OpenSLES_Android.h>

//////////////////// AudioEngine /////////////////////

typedef struct SineEngine {
	SLObjectItf object;
    SLEngineItf interface;
    SLObjectItf outputMixer;
} SineEngine;

SLresult Sine_init(SineEngine *engine) {
	// create engine & realize
	SLresult ret = slCreateEngine(&engine->object, 0, (void*)(0), 0, (void*)(0), (void*)(0));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	ret = (*engine->object)->Realize(engine->object, SL_BOOLEAN_FALSE);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	// get the engine interface, which is needed in order to create other objects
	ret = (*engine->object)->GetInterface(engine->object, SL_IID_ENGINE, &engine->interface);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	// create output mixer & realize
	ret = (*engine->interface)->CreateOutputMix(engine->interface, &engine->outputMixer, 0, 0, 0);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	ret = (*engine->outputMixer)->Realize(engine->outputMixer, SL_BOOLEAN_FALSE);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	return SL_RESULT_SUCCESS;
}

void Sine_close(SineEngine *engine) {

}

//////////////////// BufferPlayer /////////////////////

typedef struct SineBufferPlayer {
	// data
	SLDataSource audioSrc;
	SLDataSink   audioSink;

	// control
	SLVolumeItf volume;
	SLObjectItf playerObj;
	SLPlayItf   player;

	SLBufferQueueItf bufferQueue;
	SLboolean playing;
} SineBufferPlayer;

void SineBufferPlayer_callback(SLBufferQueueItf aSoundQueue, void* aContext) {
	((SineBufferPlayer*) aContext)->playing = SL_BOOLEAN_FALSE;
}

SLresult SineBufferPlayer_init(SineBufferPlayer *p, SineEngine *engine, SLuint32 numBuffers, SLuint32 numChannels) {
	SLDataLocator_AndroidSimpleBufferQueue locatorBufferQueue = {
		SL_DATALOCATOR_ANDROIDSIMPLEBUFFERQUEUE, numBuffers
	};
	SLDataFormat_PCM format = {
		SL_DATAFORMAT_PCM,
		numChannels,
		SL_SAMPLINGRATE_44_1,
		SL_PCMSAMPLEFORMAT_FIXED_16,
		SL_PCMSAMPLEFORMAT_FIXED_16,
		SL_SPEAKER_FRONT_LEFT|SL_SPEAKER_FRONT_RIGHT,
		SL_BYTEORDER_LITTLEENDIAN
	};

	p->audioSrc.pLocator = &locatorBufferQueue;
	p->audioSrc.pFormat  = &format;

	SLDataLocator_OutputMix outmix = {SL_DATALOCATOR_OUTPUTMIX, engine->outputMixer};
	p->audioSink.pLocator = &outmix;
	p->audioSink.pFormat  = (void*)(0);

	SLInterfaceID ids[2] = {SL_IID_ANDROIDSIMPLEBUFFERQUEUE, SL_IID_VOLUME};
	SLboolean req[2] = {SL_BOOLEAN_TRUE,SL_BOOLEAN_TRUE};

	SLresult ret;
	// create player & realize
	ret = (*engine->interface)->CreateAudioPlayer(engine->interface, &(p->playerObj), &(p->audioSrc), &(p->audioSink), 2, ids, req);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	ret = (*p->playerObj)->Realize(p->playerObj, SL_BOOLEAN_FALSE);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	// play interface
	ret = (*p->playerObj)->GetInterface(p->playerObj, SL_IID_PLAY, &(p->player));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	// create volume & realize
	ret = (*p->playerObj)->GetInterface(p->playerObj, SL_IID_VOLUME, &(p->volume));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	ret = (*p->playerObj)->GetInterface(p->playerObj, SL_IID_ANDROIDSIMPLEBUFFERQUEUE, &(p->bufferQueue));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	return SL_RESULT_SUCCESS;
}

void SineBufferPlayer_close(SineBufferPlayer *p) {
	(*p->playerObj)->Destroy(p->playerObj);
}

SLresult SineBufferPlayer_play(SineBufferPlayer *p, void* buffer, SLuint32 size) {
	// enqueue data
	SLresult ret;
	ret = (*p->bufferQueue)->Enqueue(p->bufferQueue, buffer, size);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	// play
	ret = (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_PLAYING);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	p->playing = SL_BOOLEAN_TRUE;
	return SL_RESULT_SUCCESS;
}

SLresult SineBufferPlayer_stop(SineBufferPlayer *p) {
	return (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_STOPPED);
}

SLresult SineBufferPlayer_pause(SineBufferPlayer *p) {
	return (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_PAUSED);
}

SLresult SineBufferPlayer_resume(SineBufferPlayer *p) {
	return (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_PLAYING);
}

SLresult SineBufferPlayer_queueCount(SineBufferPlayer *p, SLuint32 *count) {
	SLresult ret;
	SLBufferQueueState qState;
	ret = (*p->bufferQueue)->GetState(p->bufferQueue, &qState);
	*count = qState.count;
	return ret;
}

SLresult SineBufferPlayer_state(SineBufferPlayer *p, SLuint32 *state) {
	return (*p->player)->GetPlayState(p->player, state);
}

// unit: db
SLresult SineBufferPlayer_setVolume(SineBufferPlayer *p, SLmillibel v) {
	return (*p->volume)->SetVolumeLevel(p->volume, v);
}

SLresult SineBufferPlayer_getVolume(SineBufferPlayer *p, SLmillibel *v) {
	return (*p->volume)->GetVolumeLevel(p->volume, v);
}

//////////////////// StreamPlayer /////////////////////

typedef struct SineBuffer {
	void *data;
	SLuint32 size;
} SineBuffer;

typedef struct SineStreamPlayer {
	// data
	SLDataSource audioSrc;
	SLDataSink   audioSink;

	// control
	SLVolumeItf volume;
	SLObjectItf playerObj;
	SLPlayItf   player;

	// queue and buffer
	SLBufferQueueItf bufferQueue;
	SineBuffer buffers[8]; //max free buffers
	SLuint32 numBuffers;   //buffer in use
	SLuint32 freed;        //freed buffers
} SineStreamPlayer;

// just increase free buffer
void SineStreamPlayer_callback(SLBufferQueueItf aSoundQueue, void* aContext) {
	((SineStreamPlayer*) aContext)->freed += 1;
}

SLresult SineStreamPlayer_init(SineStreamPlayer *p, SineEngine *engine, SLuint32 numBuffers, SLuint32 numChannels) {
	SLDataLocator_AndroidSimpleBufferQueue locatorBufferQueue = {
		SL_DATALOCATOR_ANDROIDSIMPLEBUFFERQUEUE, numBuffers
	};
	SLDataFormat_PCM format = {
		SL_DATAFORMAT_PCM,
		numChannels,
		SL_SAMPLINGRATE_44_1,
		SL_PCMSAMPLEFORMAT_FIXED_16,
		SL_PCMSAMPLEFORMAT_FIXED_16,
		SL_SPEAKER_FRONT_LEFT|SL_SPEAKER_FRONT_RIGHT,
		SL_BYTEORDER_LITTLEENDIAN
	};

	p->audioSrc.pLocator = &locatorBufferQueue;
	p->audioSrc.pFormat  = &format;

	SLDataLocator_OutputMix outmix = {SL_DATALOCATOR_OUTPUTMIX, engine->outputMixer};
	p->audioSink.pLocator = &outmix;
	p->audioSink.pFormat  = (void*)(0);

	SLInterfaceID ids[2] = {SL_IID_ANDROIDSIMPLEBUFFERQUEUE, SL_IID_VOLUME};
	SLboolean req[2] = {SL_BOOLEAN_TRUE,SL_BOOLEAN_TRUE};

	SLresult ret;
	// create player & realize
	ret = (*engine->interface)->CreateAudioPlayer(engine->interface, &(p->playerObj), &(p->audioSrc), &(p->audioSink), 2, ids, req);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	ret = (*p->playerObj)->Realize(p->playerObj, SL_BOOLEAN_FALSE);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	// play interface
	ret = (*p->playerObj)->GetInterface(p->playerObj, SL_IID_PLAY, &(p->player));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	// create volume & realize
	ret = (*p->playerObj)->GetInterface(p->playerObj, SL_IID_VOLUME, &(p->volume));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	// setup buffer queue
	ret = (*p->playerObj)->GetInterface(p->playerObj, SL_IID_ANDROIDSIMPLEBUFFERQUEUE, &(p->bufferQueue));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	p->numBuffers = numBuffers;
	p->freed      = numBuffers;

	// register callback
	ret = (*p->bufferQueue)->RegisterCallback(p->bufferQueue, SineStreamPlayer_callback, p);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	return SL_RESULT_SUCCESS;
}

void SineStreamPlayer_feed(SineStreamPlayer *p, void* data, SLuint32 size) {
	SLuint32 freed = p->freed;
	if (freed <= 0) {
		return;
	}
	(*p->bufferQueue)->Enqueue(p->bufferQueue, data, size);
	p->freed = freed - 1;
}

void SineStreamPlayer_close(SineStreamPlayer *p) {
	(*p->playerObj)->Destroy(p->playerObj);
}

SLresult SineStreamPlayer_play(SineStreamPlayer *p) {
	return (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_PLAYING);
}

SLresult SineStreamPlayer_stop(SineStreamPlayer *p) {
	return (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_STOPPED);
}

SLresult SineStreamPlayer_pause(SineStreamPlayer *p) {
	return  (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_PAUSED);
}

SLresult SineStreamPlayer_queueCount(SineStreamPlayer *p, SLuint32 *count) {
	SLresult ret;
	SLBufferQueueState qState;
	ret = (*p->bufferQueue)->GetState(p->bufferQueue, &qState);
	*count = qState.count;
	return ret;
}

SLresult SineStreamPlayer_state(SineStreamPlayer *p, SLuint32 *state) {
	return (*p->player)->GetPlayState(p->player, state);
}

// unit: db
SLresult SineStreamPlayer_setVolume(SineStreamPlayer *p, SLmillibel v) {
	return (*p->volume)->SetVolumeLevel(p->volume, v);
}

SLresult SineStreamPlayer_getVolume(SineStreamPlayer *p, SLmillibel *v) {
	return (*p->volume)->GetVolumeLevel(p->volume, v);
}

 */
import "C"
import (
	//"unsafe"
	//"log"
	"log"
	//"unsafe"
	//"unsafe"
	"unsafe"
	"math"
)

// StaticData is small audio sampler, which will be load into memory directly.
type StaticData struct {
	bits []byte
	fmt uint32
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
	engine C.SineEngine
}

func (eng *Engine) Initialize() {
	if ret := C.Sine_init(&eng.engine); ret != OK {
		log.Println("fail to initialize opensl")
	}
}

func (eng *Engine) Destroy() {
	C.Sine_close(&eng.engine)
}

// BufferPlayer can play audio loaded as StaticData.
type BufferPlayer struct {
	player C.SineBufferPlayer
	volume float32
	state int
}

func (p *BufferPlayer) initialize(engine *Engine) {
	var (
		numBuffers = C.SLuint32(2)
		numChannels = C.SLuint32(2)
	)
	if ret := C.SineBufferPlayer_init(&p.player, &engine.engine, numBuffers, numChannels); ret != OK {
		log.Println("buffer-player init err:", ret)
	}
}

func (p *BufferPlayer) destroy() {
	C.SineBufferPlayer_close(&p.player)
}

func (p *BufferPlayer) Play(d *StaticData) {
	p.state = Playing
	if ret := C.SineBufferPlayer_play(&p.player, unsafe.Pointer(&d.bits[0]), C.SLuint32(len(d.bits))); ret != OK {
		log.Println("buffer-player play err:", ret)
	}
}

func (p *BufferPlayer) Stop() {
	p.state = Stopped
	if ret := C.SineBufferPlayer_stop(&p.player); ret != OK {
		log.Println("buffer-player stop err:", ret)
	}
}

func (p *BufferPlayer) Pause() {
	p.state = Paused
	if ret := C.SineBufferPlayer_pause(&p.player); ret != OK {
		log.Println("buffer-player pause err:", ret)
	}
}

func (p *BufferPlayer) Resume() {
	p.state = Playing
	if ret := C.SineBufferPlayer_resume(&p.player); ret != OK {
		log.Println("buffer-player resume err:", ret)
	}
}

// OpenSL's state is different from OpenAL. In OpenAL, if buffer
// queue exhausted, OpenAL will issue a 'Stop' state. But in SL,
// buffer-queue has nothing to do with player state. It' still
// playing even though queue exhausted.


func (p *BufferPlayer) State() int {
	if st := p.state; st != Playing {
		return st
	}
	ct := C.SLuint32(0)
	if ret := C.SineBufferPlayer_queueCount(&p.player, &ct); ret != OK {
		log.Print("buffer state err:", ret)
		return Initial
	}
	if ct == 0 {
		return Stopped
	}
	return Playing
}

func (p *BufferPlayer) Volume() float32 {
	return p.volume
}

func (p *BufferPlayer) SetVolume(v float32) {
	p.volume = v
	f := math.Log10(float64(v)) * 2000
	if ret := C.SineBufferPlayer_setVolume(&p.player, C.SLmillibel(f)); ret != OK {
		log.Print("stream-player set volume err:", ret)
	}
}

func (p *BufferPlayer) SetLoop(n int) {

}

// StreamPlayer can play audio loaded as StreamData.
type StreamPlayer struct {
	player C.SineStreamPlayer
	decoder      Decoder
	volume float32
	state int
}

func (p *StreamPlayer) initialize(engine *Engine) {
	var (
		numBuffers = C.SLuint32(2)
		numChannels = C.SLuint32(2)
	)
	if ret := C.SineStreamPlayer_init(&p.player, &engine.engine, numBuffers, numChannels); ret != OK {
		log.Println("stream-player init err:", ret)
	}
	p.volume = 1
}

func (p *StreamPlayer) Play(d *StreamData) {
	p.decoder = d.decoder
	p.fill()
	p.state = Playing
	// play
	if ret := C.SineStreamPlayer_play(&p.player); ret != OK {
		log.Println("stream-player play err:", ret)
	}
}

func (p *StreamPlayer) Stop() {
	p.state = Stopped
	if ret := C.SineStreamPlayer_stop(&p.player); ret != OK {
		log.Println("stream-player stop err:", ret)
	}
	if d := p.decoder; d != nil {
		d.Rewind()
	}
	p.decoder = nil
}

func (p *StreamPlayer) Pause() {
	p.state = Paused
	if ret := C.SineStreamPlayer_pause(&p.player); ret != OK {
		log.Println("stream-player pause err:", ret)
	}
}

func (p *StreamPlayer) Resume() {
	p.state = Playing
	if ret := C.SineStreamPlayer_play(&p.player); ret != OK {
		log.Println("stream-player resume err:", ret)
	}
}

func (p *StreamPlayer) State() int {
	if st := p.state; st != Playing {
		return st
	}
	ct := C.SLuint32(0)
	if ret := C.SineBufferPlayer_queueCount(&p.player, &ct); ret != OK {
		log.Print("buffer state err:", ret)
		return Initial
	}
	if ct == 0 {
		return Stopped
	}
	return Playing
}

func (p *StreamPlayer) Volume() float32 {
	return p.volume
}

func (p *StreamPlayer) SetVolume(v float32) {
	var db C.SLmillibel
	if v > 0 {
		p.volume = v
		db = C.SLmillibel(math.Log10(float64(v))*2000)
	} else {
		p.volume = 0
		db = C.SL_MILLIBEL_MIN
	}
	if ret := C.SineStreamPlayer_setVolume(&p.player, db); ret != OK {
		log.Print("stream-player set volume err:", ret)
	}
}

func (p *StreamPlayer) Tick() {
	p.fill()
}

func (p *StreamPlayer) fill() {
	d := p.decoder
	if d == nil || d.ReachEnd() {
		return // return if no more data
	}

	free := int(p.player.freed)
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
			size   = C.SLuint32(len(buffer))
		)
		C.SineStreamPlayer_feed(&p.player, unsafe.Pointer(&buffer[0]), size)
	}
}

const (
	FormatMono8    = 0x1100
	FormatMono16   = 0x1101
	FormatStereo8  = 0x1102
	FormatStereo16 = 0x1103
)

const (
	OK = C.SL_RESULT_SUCCESS
)

const (
	Initial = 0x0
	Playing = C.SL_PLAYSTATE_PLAYING
	Paused  = C.SL_PLAYSTATE_PAUSED
	Stopped = C.SL_PLAYSTATE_STOPPED
)


