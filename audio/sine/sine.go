package sine

/**
Sine is a low-level audio library. You can use it to play pcm-streams, pcm-buffer.
Compressed audio format like ogg/mp3 is not supported directly. You can implement
the DecoderFactory/Decoder interface to play those audio format.
*/
var (
	factory DecoderFactory
	engine *Engine
)

func Init(df DecoderFactory) {
	factory = df
	engine.Initialize()
}

func Destroy() {
	engine.Destroy()
}

func init() {
	R = NewAudioManager()
	engine = &Engine{};
}

// public field
var R *AudioManger

func NewBufferPlayer() (bp *BufferPlayer) {
	bp = &BufferPlayer{}; bp.initialize(engine)
	return
}

func NewStreamPlayer() (sp *StreamPlayer) {
	sp = &StreamPlayer{}; sp.initialize(engine)
	return
}

func NewSoundPool() *SoundPool {
	soundPool := &SoundPool{}; soundPool.initialize(R, engine, 8)
	return soundPool
}