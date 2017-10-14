package ap

import (
	"golang.org/x/mobile/exp/audio/al"

	"log"
	"sort"
)

/// 音频系统设计：
/// 1. 资源管理在 AudioManager 里面，所有的资源通过 id 索引
/// 2. 播放系统(PlayerContext)管理8个硬件播放通道
/// 3. 音频资源有一个默认的优先级，播放的时候也可以设置一个临时的优先级
///
/// 调用play方法只会参与优先级计算，在 Frame() 方法中会根据优先级执行最
/// 最终播放.

const (
	MAX_CHANNEL_SIZE = 8
)

type playCall struct {
	id uint16
	p  uint16
}

type chanRef struct {
	p uint16
	ref *Channel
}

// manager priority and audio play
type PlayContext struct {
	R *AudioManger

	// state
	pause bool
	mute  bool

	// hardware channel
	playChan []chanRef

	// priority queue:4 3 2 1
	pQueue []playCall
	pIndex int

	// underlying data
	p_ref  [MAX_CHANNEL_SIZE]chanRef
	p_call [MAX_CHANNEL_SIZE]playCall
	p_chan [MAX_CHANNEL_SIZE]Channel
}

func NewPlayContext(am *AudioManger) *PlayContext {
	pc := &PlayContext {
		R: am,

		pause:false,
		mute:false,
	}
	pc.pQueue = pc.p_call[:]
	pc.pIndex = 0
	pc.playChan = pc.p_ref[:]

	for i := range pc.p_chan {
		pc.playChan[i].ref = &pc.p_chan[i]
	}
	return pc
}

// init device/source/listener
func (pc *PlayContext) Init() error{
	err := al.OpenDevice()
	if err != nil {
		return err
	}
	al.SetListenerPosition(al.Vector{0, 0, 0})

	if err := pc.p_chan[0].Create(Stream);err != nil {
		return err
	}
	for i := 1; i < MAX_CHANNEL_SIZE; i++ {
		if err := pc.p_chan[i].Create(Static); err != nil {
			return err
		}
	}
	return nil
}

func (pc *PlayContext) Destroy() {
	al.CloseDevice()
}

func (pc *PlayContext) Mute(mute bool) {

}

func (pc *PlayContext) Pause(pause bool) {

}

// 直接把Play命令放入优先级的问题，

// 音频资源默认包含一个优先级，如果此处设置的优先级不为0，那么使用此处的
// 设置，否则使用默认的优先级
func (pc *PlayContext) Play(id uint16, priority uint16) {
	// 得到优先级
	p := priority
	if p == 0 {
		if ok, sound := pc.R.Sound(id); ok {
			p = sound.Priority
		} else {
			log.Println("Invalid source id")
			return
		}
	}
	// 加入优先级队列
	play := playCall{id, priority}
	insert := pc.pIndex
	for i := 0; i < pc.pIndex; i++ {
		if pc.pQueue[i].p < p {
			insert = i
		}
	}

	log.Println("insert position:", insert, " res id:", id)

	if insert < MAX_CHANNEL_SIZE {
		// queue is full
		if insert == pc.pIndex {
			pc.pQueue[insert] = play
		} else {
			copy(pc.pQueue[insert+1:], pc.pQueue[insert:pc.pIndex-1])
			pc.pQueue[insert] = play
		}
		pc.pIndex ++
	}
}

func (pc *PlayContext) Stop(id uint16) {

}

// play
func (pc *PlayContext) NextFrame() {
	// update channel state
	for _, ch := range pc.playChan {
		ch.ref.UpdateState()
		if ch.ref.State == STOP {
			ch.p = 0
		}
	}

	// sort chanRef 4 3 2 1
	sort.Slice(pc.playChan, func(i, j int) bool {
		return pc.playChan[i].p < pc.playChan[j].p
	})

	// play priority queue
	for i,j := 0, 0; i < pc.pIndex && j < MAX_CHANNEL_SIZE; i,j = i+1, j+1 {
		if play, ref := pc.pQueue[i], pc.playChan[j];play.p > ref.p {
			channel := ref.ref
			if channel.State != STOP {
				channel.Halt()
			}
			channel.Play(&pc.R.soundPool[play.id])
		}
	}
	// reset queue
	pc.pIndex = 0
}
