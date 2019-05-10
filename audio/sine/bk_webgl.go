// Copyright 2017 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// vanishs 修改于github.com/hajimehoshi/oto

//+build js

package sine

import (
	"log"
	"time"

	"syscall/js"

	"korok.io/korok/hid"
)

type driver struct {
	sampleRate      int
	channelNum      int
	bitDepthInBytes int
	nextPos         float64
	tmp             []byte
	bufferSize      int
	context         js.Value
	lastTime        float64
	lastAudioTime   float64
	ready           bool
	bytesPerSecond  int
}

const audioBufferSamples = 3200

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func newDriver(sampleRate, channelNum, bitDepthInBytes, bufferSize int) (*driver, error) {

	p := &driver{
		sampleRate:      sampleRate,
		channelNum:      channelNum,
		bitDepthInBytes: bitDepthInBytes,
		context:         hid.AudioCtx,
		bufferSize:      max(bufferSize, audioBufferSamples*channelNum*bitDepthInBytes),
		bytesPerSecond:  sampleRate * channelNum * bitDepthInBytes,
	}

	// setCallback := func(event string) {
	// 	var f js.Func
	// 	f = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 		if !p.ready {
	// 			p.context.Call("resume")
	// 			p.ready = true
	// 		}
	// 		js.Global().Get("document").Call("removeEventListener", event, f)
	// 		return nil
	// 	})
	// 	js.Global().Get("document").Call("addEventListener", event, f)

	// }

	// // Browsers require user interaction to start the audio.
	// // https://developers.google.com/web/updates/2017/09/autoplay-policy-changes#webaudio
	// setCallback("touchend")
	// setCallback("keyup")
	// setCallback("mouseup")
	return p, nil
}

func toLR(data []byte) ([]float32, []float32) {
	const max = 1 << 15

	l := make([]float32, len(data)/4)
	r := make([]float32, len(data)/4)
	for i := 0; i < len(data)/4; i++ {
		l[i] = float32(int16(data[4*i])|int16(data[4*i+1])<<8) / max
		r[i] = float32(int16(data[4*i+2])|int16(data[4*i+3])<<8) / max
	}
	return l, r
}

func nowInSeconds() float64 {
	return js.Global().Get("performance").Call("now").Float() / 1000.0
}

func (p *driver) TryWrite(data []byte) (int, error) {
	// if !p.ready {
	// 	return 0, nil
	// }

	n := min(len(data), max(0, p.bufferSize-len(p.tmp)))
	p.tmp = append(p.tmp, data[:n]...)

	c := p.context.Get("currentTime").Float()
	now := nowInSeconds()

	if p.lastTime != 0 && p.lastAudioTime != 0 && p.lastAudioTime >= c && p.lastTime != now {
		// Unfortunately, currentTime might not be precise enough on some devices
		// (e.g. Android Chrome). Adjust the audio time with OS clock.
		c = p.lastAudioTime + now - p.lastTime
	}

	p.lastAudioTime = c
	p.lastTime = now

	if p.nextPos < c {
		p.nextPos = c
	}

	// It's too early to enqueue a buffer.
	// Highly likely, there are two playing buffers now.
	if c+float64(p.bufferSize/p.bitDepthInBytes/p.channelNum)/float64(p.sampleRate) < p.nextPos {
		return n, nil
	}

	le := audioBufferSamples * p.bitDepthInBytes * p.channelNum
	if len(p.tmp) < le {
		return n, nil
	}

	buf := p.context.Call("createBuffer", p.channelNum, audioBufferSamples, p.sampleRate)
	l, r := toLR(p.tmp[:le])
	tl := js.TypedArrayOf(l)
	tr := js.TypedArrayOf(r)
	if buf.Get("copyToChannel") != js.Undefined() {
		buf.Call("copyToChannel", tl, 0, 0)
		buf.Call("copyToChannel", tr, 1, 0)
	} else {
		// copyToChannel is not defined on Safari 11
		buf.Call("getChannelData", 0).Call("set", tl)
		buf.Call("getChannelData", 1).Call("set", tr)
	}
	tl.Release()
	tr.Release()

	s := p.context.Call("createBufferSource")
	s.Set("buffer", buf)
	s.Call("connect", p.context.Get("destination"))
	s.Call("start", p.nextPos)
	p.nextPos += buf.Get("duration").Float()

	p.tmp = p.tmp[le:]
	return n, nil
}

func (p *driver) Close() error {
	return nil
}

///////////////////////////////

// StaticData is small audio sampler, which will be load into memory directly.
type StaticData struct {
	bits []byte
	fmt  uint32
	freq int32
}

func (d *StaticData) Create(fmt uint32, bits []byte, freq int32) {
	d.fmt = fmt
	d.bits = bits
	d.freq = freq
}

// StreamData will decode pcm-data at runtime. It's used to play big audio files(like .ogg).
type StreamData struct {
	decoder Decoder
}

func (d *StreamData) Create(file string, ft FileType) {
	decoder, err := factory.NewDecoder(file, ft)
	if err != nil {
		log.Println(err)
		return
	}
	d.decoder = decoder
}

type Engine struct {
	d *driver
}

func (eng *Engine) Initialize() {
	var err error
	eng.d, err = newDriver(44100, 2, 2, 4096)
	if err != nil {
		panic(err)
	}
}

func (eng *Engine) Destroy() {
	err := eng.d.Close()
	if err != nil {
		panic(err)
	}
}

// BufferPlayer can play audio loaded as StaticData.
type BufferPlayer struct {
	engine *Engine
	status uint32
}

func (p *BufferPlayer) initialize(engine *Engine) {
	p.engine = engine
	p.status = Stopped
}

func (p *BufferPlayer) Play(data *StaticData) {
	go func() {
		written := 0
		buf := data.bits
		for len(buf) > 0 {
			// if d.driver == nil {
			// 	return written, errClosed
			// }
			n, _ := engine.d.TryWrite(buf)
			written += n
			// if err != nil {
			// 	return written, err
			// }
			buf = buf[n:]
			// When not all buf is written, the underlying buffer is full.
			// Mitigate the busy loop by sleeping (#10).
			if len(buf) > 0 {
				t := time.Second * time.Duration(engine.d.bufferSize) / time.Duration(engine.d.bytesPerSecond) / 8
				time.Sleep(t)
			}
		}
	}()

	// engine.d.TryWrite(data.bits)
}

func (p *BufferPlayer) Stop() {

}

func (p *BufferPlayer) Pause() {

}

func (p *BufferPlayer) Resume() {

}

func (p *BufferPlayer) Volume() float32 {
	return 0
}

func (p *BufferPlayer) SetVolume(v float32) {
}

func (p *BufferPlayer) SetLoop(loop int) {

}

func (p *BufferPlayer) State() uint32 {
	return p.status
}

// StreamPlayer can play audio loaded as StreamData.
type StreamPlayer struct {
	engine  *Engine
	status  uint32
	feed    chan []byte
	decoder Decoder
}

func (p *StreamPlayer) initialize(engine *Engine) {
	p.engine = engine
	p.status = Stopped
	p.feed = make(chan []byte, 128)
}

func (p *StreamPlayer) Play(stream *StreamData) {
	p.decoder = stream.decoder
	p.fill()
	go func() {

		for {

			written := 0
			buf := <-p.feed
			for len(buf) > 0 {
				// if d.driver == nil {
				// 	return written, errClosed
				// }
				n, _ := engine.d.TryWrite(buf)
				written += n
				// if err != nil {
				// 	return written, err
				// }
				buf = buf[n:]
				// When not all buf is written, the underlying buffer is full.
				// Mitigate the busy loop by sleeping (#10).
				if len(buf) > 0 {
					t := time.Second * time.Duration(engine.d.bufferSize) / time.Duration(engine.d.bytesPerSecond) / 8
					time.Sleep(t)
				}
			}
		}

	}()
}

func (p *StreamPlayer) fill() {
	d := p.decoder
	if d == nil || d.ReachEnd() {
		return // return if no more data
	}

	// feed the free buffer one by one.
	for {
		if n := d.Decode(); n == 0 {
			break
		}
		// var (
		// 	buffer = d.Buffer()
		// 	size   = C.SLuint32(len(buffer))
		// )
		// C.SineStreamPlayer_feed(&p.player, unsafe.Pointer(&buffer[0]), size)

		p.feed <- d.Buffer()
	}
}

func (p *StreamPlayer) Stop() {

}

func (p *StreamPlayer) Pause() {
}

func (p *StreamPlayer) Resume() {

}

func (p *StreamPlayer) State() uint32 {
	return p.status
}

func (p *StreamPlayer) Volume() float32 {
	return 0
}

func (p *StreamPlayer) SetVolume(v float32) {
}

func (p *StreamPlayer) Tick() {
	p.fill()
}

const (
	FormatMono8    = 0x1100
	FormatMono16   = 0x1101
	FormatStereo8  = 0x1102
	FormatStereo16 = 0x1103
)

// AL error
const ()

// AL state
const (
	Initial = 0x1011
	Playing = 0x1012
	Paused  = 0x1013
	Stopped = 0x1014
)
