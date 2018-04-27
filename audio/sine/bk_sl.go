// +build android

package sine

/*
 #cgo LDFLAGS: -landroid -lOpenSLES

 #include <SLES/OpenSLES.h>
 #include <SLES/OpenSLES_Android.h>


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

int lala(SineEngine *engine) {
	return 3;
}

void Sine_close(SineEngine *engine) {

}

//SLresult Sine_createAudioPlayer(SLObjectItf *objectItf,
//                            SLDataSource *audioSource,
//                            SLDataSink *audioSink) {
//	const SLInterfaceID ids[] = {SL_IID_BUFFERQUEUE};
//    const SLboolean reqs[] = {SL_BOOLEAN_TRUE};
//
//    return (*mEngineInterface)->CreateAudioPlayer(mEngineInterface, objectItf, audioSource,
//                                                  audioSink,
//                                                  sizeof(ids) / sizeof(ids[0]), ids, reqs);
//}

typedef struct SineBufferPlayer {
	// data
	SLDataSource audioSrc;
	SLDataSink   audioSink;

	// control
	SLVolumeItf fdPlayerVolume;
	SLObjectItf playerObj;
	SLPlayItf   player;

	SLBufferQueueItf bufferQueue;
} SineBufferPlayer;

SLresult SineBufferPlayer_init(SineBufferPlayer *p, SineEngine *engine) {
	SLDataLocator_AndroidSimpleBufferQueue locatorBufferQueue;
	locatorBufferQueue.locatorType = SL_DATALOCATOR_ANDROIDSIMPLEBUFFERQUEUE;
	locatorBufferQueue.numBuffers = 16;

	SLDataFormat_PCM format;
	format.formatType = SL_DATAFORMAT_PCM;
	format.numChannels = 2;
	format.samplesPerSec = SL_SAMPLINGRATE_44_1;
	format.bitsPerSample = SL_PCMSAMPLEFORMAT_FIXED_16;
	format.containerSize = SL_PCMSAMPLEFORMAT_FIXED_16;
	format.channelMask   = SL_SPEAKER_FRONT_LEFT|SL_SPEAKER_FRONT_RIGHT;
	format.endianness    = SL_BYTEORDER_LITTLEENDIAN;

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
	ret = (*p->playerObj)->GetInterface(p->playerObj, SL_IID_VOLUME, &(p->fdPlayerVolume));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	ret = (*p->playerObj)->GetInterface(p->playerObj, SL_IID_ANDROIDSIMPLEBUFFERQUEUE, &(p->bufferQueue));
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	return SL_RESULT_SUCCESS;
}

//void SineBufferPlayer_close(SineBufferPlayer *p) {
//	stream.mObjectInterface.Destroy(stream.mObjectInterface);
//}

SLresult SineBufferPlayer_play(SineBufferPlayer *p, void* buffer, SLuint32 size) {
	// enqueue data
	SLresult ret;
	ret = (*p->bufferQueue)->Clear(p->bufferQueue);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	(*p->bufferQueue)->Enqueue(p->bufferQueue, buffer, size);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	// play
	ret = (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_PLAYING);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	return SL_RESULT_SUCCESS;
}

int Simple1(SineBufferPlayer *p) {
	return 1;
}

SLresult SineBufferPlayer_stop(SineBufferPlayer *p) {
	return SL_RESULT_SUCCESS;
}

SLresult SineBufferPlayer_pause(SineBufferPlayer *p) {
	SLresult ret;
	ret = (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_PAUSED);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	return SL_RESULT_SUCCESS;
}

SLresult SineBufferPlayer_resume(SineBufferPlayer *p) {
	SLresult ret;
	ret = (*p->player)->SetPlayState(p->player, SL_PLAYSTATE_PLAYING);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}
	return SL_RESULT_SUCCESS;
}

SLresult SineBufferPlayer_state(SineBufferPlayer *p, SLuint32 *state) {
	SLresult ret;

	SLBufferQueueState qState;
	ret = (*p->bufferQueue)->GetState(p->bufferQueue, &qState);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	SLuint32 pState;
	ret = (*p->player)->GetPlayState(p->player, &pState);
	if (ret != SL_RESULT_SUCCESS) {
		return ret;
	}

	if ((pState == SL_PLAYSTATE_PLAYING) && (qState.count == 0)) {
		*state = SL_PLAYSTATE_PLAYING; // TODO 状态有待优化..
	}
	return SL_RESULT_SUCCESS;
}
//
//void SineStream_setVolume(stream *SineStream, v float) {
//	//tood gain_to_attenuation... 不知道何意
//	stream.fdPlayerVolume->SetVolumeLevel(stream.fdPlayerVolume, )
//}
//
//float SineStream getVolume(stream *SineStream) {
//	SLmillibel millibel;
//	fbPlayerVolume->GetVolumeLevel(stream.fdPlayerVolume, &millibel);
//	// todo
//	return
//}

typedef struct SineStreamPlayer {

} SineStreamPlayer;
//

 */
import "C"
import (
	//"unsafe"
	//"log"
	"log"
	//"unsafe"
	//"unsafe"
	"unsafe"
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
	file string
	ft FileType
}

func (d *StreamData) Create(file string, ft FileType) {
	d.file = file
	d.ft = ft
}


type Engine struct {
	engine C.SineEngine
}

func (eng *Engine) Initialize() {
	if ret := C.Sine_init(&eng.engine); ret != OK {
		log.Println("fail to initialize opensl")
	}


	r3 := C.lala(&eng.engine)
	log.Println("lala:", r3)
}

func (eng *Engine) Destroy() {

}

// BufferPlayer can play audio loaded as StaticData.
type BufferPlayer struct {
	player C.SineBufferPlayer
}

func (p *BufferPlayer) initialize(engine *Engine) {
	if ret := C.SineBufferPlayer_init(&p.player, &engine.engine); ret != OK {
		log.Println("buffer-player init err:", ret)
	}
}

func (p *BufferPlayer) Play(d *StaticData) {
	//if p == nil {
	//	log.Println("p is nil")
	//}
	//if d == nil {
	//	log.Println("d is nil")
	//}
	//if d.bits == nil {
	//	log.Println("bits is nil")
	//}
	if ret := C.SineBufferPlayer_play(&p.player, unsafe.Pointer(&d.bits[0]), C.SLuint32(len(d.bits))); ret != OK {
		log.Println("buffer-player play err:", ret)
	}
}

func (p *BufferPlayer) Stop() {
	if ret := C.SineBufferPlayer_stop(&p.player); ret != OK {
		log.Println("buffer-player stop err:", ret)
	}
}

func (p *BufferPlayer) Pause() {
	if ret := C.SineBufferPlayer_pause(&p.player); ret != OK {
		log.Println("buffer-player pause err:", ret)
	}
}

func (p *BufferPlayer) Resume() {
	if ret := C.SineBufferPlayer_resume(&p.player); ret != OK {
		log.Println("buffer-player resume err:", ret)
	}
}

// StreamPlayer can play audio loaded as StreamData.
type StreamPlayer struct {
	player *C.SineStreamPlayer
}

func (p *StreamPlayer) initialize(engine *Engine) {

}

func (p *StreamPlayer) Play(d *StreamData) {

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



