package audio

import "golang.org/x/mobile/exp/audio/al"

/**
音频接口设计....音频和其它组件的

*/

type SourceState int

const (
	READY SourceState = iota
	PLAYING
	STOP
)

// 组件设计
type SourceComp struct {
	// 音频资源
	Id string

	// 静音
	Mute bool

	// 循环
	Loop bool

	// 音量 [0, 1]
	Volume float32

	// 声源状态
	state SourceState
}

func (s *SourceComp) Play(id string) {
	s.Id = id
	s.state = READY
}

// 系统设计
type AudioSystem struct {
	bank *Bank
	player *Player
	comps []SourceComp
}

func (sys *AudioSystem) Update(dt float32) {
	p := sys.player

	for i := range sys.comps {
		comp := &sys.comps[i]
		if comp.state == READY {
			p.Play(sys.id2buffer(comp.Id))
			comp.state = PLAYING
		}
	}
}

func (sys *AudioSystem) id2buffer(id string) al.Buffer {
	return sys.bank.GetBuffer(id)
}








