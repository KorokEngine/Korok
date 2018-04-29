package sine

func Init(df DecoderFactory) {
	factory = df
	engine = &Engine{}; engine.Initialize()
	bufferPlayer = &BufferPlayer{}; bufferPlayer.initialize(engine)
	streamPlayer = &StreamPlayer{}; streamPlayer.initialize(engine)
}

func Destroy() {
	engine.Destroy()
}

func Tick() {
	streamPlayer.Tick()
}

func Play(id uint16) {
	if sound, ok := R.Sound(id); ok {
		switch d :=  sound.Data.(type) {
		case *StaticData:
			PlayStatic(d)
		case *StreamData:
			PlayStream(d)
		}
	}
}

func Stop(id uint16) {

}

func Pause(id uint16) {

}

func Resume(id uint16) {

}

func PlayStream(d *StreamData) *StreamPlayer{
	streamPlayer.Play(d)
	return streamPlayer
}

func PlayStatic(d *StaticData) *BufferPlayer {
	bufferPlayer.Play(d)
	return bufferPlayer
}

func init() {
	R = NewAudioManager()
}

// public field
var R *AudioManger

// private field
var engine *Engine
var factory DecoderFactory

var bufferPlayer *BufferPlayer
var streamPlayer *StreamPlayer

func NewBufferPlayer() *BufferPlayer {
	bufferPlayer = &BufferPlayer{}
	bufferPlayer.initialize(engine)
	return bufferPlayer
}

func NewStreamPlayer() *StreamPlayer {
	streamPlayer = &StreamPlayer{}
	streamPlayer.initialize(engine)
	return streamPlayer;
}

