package audio

import (
	"golang.org/x/mobile/exp/audio/al"
	"io"
	"korok/audio/codec"
)

/**
音频数据管理
 */
type Bank struct {
	buffers []al.Buffer
}

// 加载一段音频
// buffer - al的Buffer缓存
// id     - 音频 id
// err    - 错误
func (b *Bank) Load(reader io.Reader) (buffer al.Buffer, id string, err error) {
	// gen buffer
	bufs := al.GenBuffers(1)
	b.buffers = append(b.buffers, bufs[0])

	// decode file
	audio, err := codec.Decode("wav", reader)

	stream := audio.Stream()
	buf := make([]uint8, 0)
	io.ReadFull(stream, buf)

	// upload buffer data
	bufs[0].BufferData(0, buf, 44100)

	return bufs[0], "", nil
}

// TODO
func (b *Bank) GetBuffer(id string) al.Buffer {
	return 0
}


