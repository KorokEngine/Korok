package sine

import (
	"log"
	"sort"
)

const (
	MaxChannelSize = 8
)

type SoundPollCallback func()

type SoundChannel struct {
	BufferPlayer
	playing bool
	priority int
	channelId int
}

// manager priority and audio play
type SoundPool struct {
	R *AudioManger
	cb SoundPollCallback

	// state
	pause bool
	mute  bool

	nextChanId int
	channels []*SoundChannel  // quick ref
	_channels  []SoundChannel // real data
}

func (sp *SoundPool) initialize (am *AudioManger, engine *Engine, maxChannel int) *SoundPool {
	sp.R = am
	size := maxChannel
	if size == 0 {
		size = MaxChannelSize
	}
	sp.channels = make([]*SoundChannel, size)
	sp._channels = make([]SoundChannel, size)

	for i := range sp.channels {
		sp._channels[i].initialize(engine)
		sp.channels[i] = &sp._channels[i]
	}
	return sp
}

func (sp *SoundPool) Destroy() {
	//al.CloseDevice()
}

func (sp *SoundPool) Mute(mute bool) {
	sp.mute = mute
}

func (sp *SoundPool) Pause(pause bool) {
	sp.pause = pause
}


func (sp *SoundPool) Tick() {
	for _, ch := range sp.channels {
		if ch.playing && ch.State() == Stopped {
			ch.playing = false
			ch.priority = 0
		}
	}
}

func (sp *SoundPool) Play(id uint16, priority int) (chanId int){
	sound, ok := sp.R.Sound(id)
	if !ok {
		log.Println("invalid sound id:", id)
		return
	}
	static, ok := sound.Data.(*StaticData)
	if !ok {
		return
	}

	// allocate a channel
	channel, ok := sp.allocChannel(priority)
	if !ok {
		log.Print("no channel available")
		return
	}

	// play sample with the channel
	chanId = sp.nextChanId; sp.nextChanId ++
	channel.channelId = chanId
	channel.Play(static)
	return chanId
}

// channels ascend order
func (sp *SoundPool) allocChannel(p int) (sc *SoundChannel, ok bool) {
	// allocate a channel
	if c := sp.channels[0]; !c.playing || c.priority < p {
		sc = c; ok = true
		c.priority = p
		c.playing = true
	}

	// update priority
	if sc != nil {
		sort.SliceStable(sp.channels, func(i, j int) bool {
			return sp.channels[i].priority < sp.channels[j].priority
		})
	}
	return
}

func (sp *SoundPool) StopChan(chanId int) {
	if ch, ok := sp.findChannel(chanId); ok {
		ch.Stop()
	}
}

func (sp *SoundPool) PauseChan(chanId int) {
	if ch, ok := sp.findChannel(chanId); ok {
		ch.Pause()
	}
}

func (sp *SoundPool) ResumeChan(chanId int) {
	if ch, ok := sp.findChannel(chanId); ok {
		ch.Resume()
	}
}

func (sp *SoundPool) SetVolume(chanId int, leftVolume, rightVolume float32) {
	if ch, ok := sp.findChannel(chanId); ok {
		ch.SetVolume(leftVolume, rightVolume)
	}
}

func (sp *SoundPool) SetLoop(chanId int, loop int) {
	if ch, ok := sp.findChannel(chanId); ok {
		ch.SetLoop(loop)
	}
}

func (sp *SoundPool) findChannel(chanId int) (sc *SoundChannel, ok bool) {
	for _, ch := range sp.channels {
		if ch.channelId == chanId {
			sc, ok = ch, true; break
		}
	}
	return
}

func (sp *SoundPool) SetCallback(cb SoundPollCallback) {
	sp.cb = cb
}
