package audio

// TODO
// positional audio component/system

type SourceComp struct {
	// 音频资源
	Id uint16
	// 优先级
	P  uint16

	// 静音
	Mute bool

	// 循环
	Loop bool

	// 音量 [0, 1]
	Volume float32
}

