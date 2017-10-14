package ap

/// Initialize the AudioPlayer
func Init() error{
	return g_ctx.Init()
}

/// Destroy AudioPlayer
func Destroy() {
	g_ctx.Destroy()
}

// Mute the AudioPlayer
func Mute(mute bool) {
	g_ctx.Mute(mute)
}

// Pause the AudioPlayer
func Pause(pause bool) {
	g_ctx.Pause(pause)
}

/// Play a sound (by id)
///
/// default:priority=0
func Play(id uint16, priority uint16) {
	g_ctx.Play(id, priority)
}

func Stop(id uint16) {
	g_ctx.Stop(id)
}

// advance to next frame
func NextFrame() {
	g_ctx.NextFrame()
}

func SetDecoderFactory(factory DecoderFactory) {
	g_df = factory
}

////////// static & global filed

var R *AudioManger
var g_ctx *PlayContext
var g_df DecoderFactory

func init() {
	R = NewAudioManager()
	g_ctx = NewPlayContext(R)
}


