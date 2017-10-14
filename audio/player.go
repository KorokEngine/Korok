package audio

// player 管理 source 资源，执行音频播放
// 1. SamplerPlayer 播放音效
// 2. MusicPlayer 播放背景音乐

type SamplerPlayer struct {
}

func (p *SamplerPlayer) Play(id uint16) {

}

type MusicPlayer struct {
	id uint16
}

func NewMusicPlayer() (*MusicPlayer, error) {
	return nil, nil
}


