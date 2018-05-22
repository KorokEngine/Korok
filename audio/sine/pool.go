package sine

import (
	"log"
	"sort"
)

const (
	MaxChannelSize = 8
)

// TODO
type SoundPollCallback func(id uint16)

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
	volume float32

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
	sp.volume = 1 // default value

	for i := range sp.channels {
		sp._channels[i].initialize(engine)
		sp.channels[i] = &sp._channels[i]
		sp.channels[i].priority = -1
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

// TODO: This method has little Latency. Burst effects will
// use up all the play-channel quickly.
func (sp *SoundPool) Tick() {
	for _, ch := range sp.channels {
		if ch.playing && ch.State() == Stopped {
			ch.playing = false
			ch.priority = -1
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
	ch, ok := sp.allocChannel(priority)
	if !ok {
		log.Print("no channel available")
		return
	}

	// play sample with the channel
	sp.nextChanId ++; chanId = sp.nextChanId
	ch.channelId = chanId
	if ch.playing {
		ch.Stop()
	}
	ch.Play(static)
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

// SetChanVolume sets volume for the specified channel. It
// may fail if can't find the channel.
func (sp *SoundPool) SetChanVolume(chanId int, v float32) {
	if ch, ok := sp.findChannel(chanId); ok {
		ch.SetVolume(v)
	}
}

// GetChanVolume gets volume from the specified channel. Return
// false if not found.
func (sp *SoundPool) GetChanVolume(chanId int) (float32, bool) {
	if ch, ok := sp.findChannel(chanId); ok {
		return ch.Volume(), true
	}
	return 0, false
}

// SetVolume set volume for all the channels in the pool.
func (sp *SoundPool) SetVolume(v float32) {
	for _, ch := range sp.channels {
		ch.SetVolume(v)
	}
	sp.volume = v
}

func (sp *SoundPool) Volume() float32 {
	return sp.volume
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
