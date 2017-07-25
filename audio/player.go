package audio

import "golang.org/x/mobile/exp/audio/al"

type Player struct {
	sources []al.Source
}

func NewPlayer(size int) (*Player, error){
	err := al.OpenDevice()
	if err != nil {
		return nil, err
	}

	p := new(Player)
	p.sources = al.GenSources(size)

	return p, nil
}

func (p*Player) Play(b al.Buffer) {
	s := p.available()
	s.QueueBuffers(b)
}

// 选取一个可用的源
func (p*Player) available() al.Source{
	return p.sources[0]
}
