package ap

/// Initialize the AudioPlayer
func Init() error{
	return ctx.Init()
}

/// Destroy AudioPlayer
func Destroy() {
	ctx.Destroy()
}

// Mute the AudioPlayer
func Mute(mute bool) {
	ctx.Mute(mute)
}

// Pause the AudioPlayer
func Pause(pause bool) {
	ctx.Pause(pause)
}

/// Play a sound (by id)
///
/// default:priority=0
func Play(id uint16, priority uint16) {
	ctx.Play(id, priority)
}

func Stop(id uint16) {
	ctx.Stop(id)
}

// advance to next frame
func NextFrame() {
	ctx.NextFrame()
}

func SetDecoderFactory(factory DecoderFactory) {
	factory = factory
}

////////// static & global filed

var R *AudioManger
var ctx *PlayContext
var factory DecoderFactory

func init() {
	R = NewAudioManager()
	ctx = NewPlayContext(R)
}


