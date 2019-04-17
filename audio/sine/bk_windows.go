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

//+build windows

package sine

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	winmm = windows.NewLazySystemDLL("winmm")
)

var (
	procWaveOutOpen          = winmm.NewProc("waveOutOpen")
	procWaveOutClose         = winmm.NewProc("waveOutClose")
	procWaveOutPrepareHeader = winmm.NewProc("waveOutPrepareHeader")
	procWaveOutWrite         = winmm.NewProc("waveOutWrite")
)

type wavehdr struct {
	lpData          uintptr
	dwBufferLength  uint32
	dwBytesRecorded uint32
	dwUser          uintptr
	dwFlags         uint32
	dwLoops         uint32
	lpNext          uintptr
	reserved        uintptr
}

type waveformatex struct {
	wFormatTag      uint16
	nChannels       uint16
	nSamplesPerSec  uint32
	nAvgBytesPerSec uint32
	nBlockAlign     uint16
	wBitsPerSample  uint16
	cbSize          uint16
}

const (
	waveFormatPCM = 1
	whdrInqueue   = 16
)

type mmresult uint

const (
	mmsyserrNoerror       mmresult = 0
	mmsyserrError         mmresult = 1
	mmsyserrBaddeviceid   mmresult = 2
	mmsyserrAllocated     mmresult = 4
	mmsyserrInvalidhandle mmresult = 5
	mmsyserrNodriver      mmresult = 6
	mmsyserrNomem         mmresult = 7
	waveerrBadformat      mmresult = 32
	waveerrStillplaying   mmresult = 33
	waveerrUnprepared     mmresult = 34
	waveerrSync           mmresult = 35
)

func (m mmresult) String() string {
	switch m {
	case mmsyserrNoerror:
		return "MMSYSERR_NOERROR"
	case mmsyserrError:
		return "MMSYSERR_ERROR"
	case mmsyserrBaddeviceid:
		return "MMSYSERR_BADDEVICEID"
	case mmsyserrAllocated:
		return "MMSYSERR_ALLOCATED"
	case mmsyserrInvalidhandle:
		return "MMSYSERR_INVALIDHANDLE"
	case mmsyserrNodriver:
		return "MMSYSERR_NODRIVER"
	case mmsyserrNomem:
		return "MMSYSERR_NOMEM"
	case waveerrBadformat:
		return "WAVEERR_BADFORMAT"
	case waveerrStillplaying:
		return "WAVEERR_STILLPLAYING"
	case waveerrUnprepared:
		return "WAVEERR_UNPREPARED"
	case waveerrSync:
		return "WAVEERR_SYNC"
	}
	return fmt.Sprintf("MMRESULT (%d)", m)
}

type winmmError struct {
	fname    string
	errno    windows.Errno
	mmresult mmresult
}

func (e *winmmError) Error() string {
	if e.errno != 0 {
		return fmt.Sprintf("winmm error at %s: Errno: %d", e.fname, e.errno)
	}
	if e.mmresult != mmsyserrNoerror {
		return fmt.Sprintf("winmm error at %s: %s", e.fname, e.mmresult)
	}
	return fmt.Sprintf("winmm error at %s", e.fname)
}

func waveOutOpen(f *waveformatex) (uintptr, error) {
	const (
		waveMapper   = 0xffffffff
		callbackNull = 0
	)
	var w uintptr
	r, _, e := procWaveOutOpen.Call(uintptr(unsafe.Pointer(&w)), waveMapper, uintptr(unsafe.Pointer(f)),
		0, 0, callbackNull)
	runtime.KeepAlive(f)
	if e.(windows.Errno) != 0 {
		return 0, &winmmError{
			fname: "waveOutOpen",
			errno: e.(windows.Errno),
		}
	}
	if mmresult(r) != mmsyserrNoerror {
		return 0, &winmmError{
			fname:    "waveOutOpen",
			mmresult: mmresult(r),
		}
	}
	return w, nil
}

func waveOutClose(hwo uintptr) error {
	r, _, e := procWaveOutClose.Call(hwo)
	if e.(windows.Errno) != 0 {
		return &winmmError{
			fname: "waveOutClose",
			errno: e.(windows.Errno),
		}
	}
	// WAVERR_STILLPLAYING is ignored.
	if mmresult(r) != mmsyserrNoerror && mmresult(r) != waveerrStillplaying {
		return &winmmError{
			fname:    "waveOutClose",
			mmresult: mmresult(r),
		}
	}
	return nil
}

func waveOutPrepareHeader(hwo uintptr, pwh *wavehdr) error {
	r, _, e := procWaveOutPrepareHeader.Call(hwo, uintptr(unsafe.Pointer(pwh)), unsafe.Sizeof(wavehdr{}))
	runtime.KeepAlive(pwh)
	if e.(windows.Errno) != 0 {
		return &winmmError{
			fname: "waveOutPrepareHeader",
			errno: e.(windows.Errno),
		}
	}
	if mmresult(r) != mmsyserrNoerror {
		return &winmmError{
			fname:    "waveOutPrepareHeader",
			mmresult: mmresult(r),
		}
	}
	return nil
}

func waveOutWrite(hwo uintptr, pwh *wavehdr) error {
	r, _, e := procWaveOutWrite.Call(hwo, uintptr(unsafe.Pointer(pwh)), unsafe.Sizeof(wavehdr{}))
	runtime.KeepAlive(pwh)
	if e.(windows.Errno) != 0 {
		return &winmmError{
			fname: "waveOutWrite",
			errno: e.(windows.Errno),
		}
	}
	if mmresult(r) != mmsyserrNoerror {
		return &winmmError{
			fname:    "waveOutWrite",
			mmresult: mmresult(r),
		}
	}
	return nil
}

type header struct {
	buffer  []byte
	waveHdr *wavehdr
}

func newHeader(waveOut uintptr, bufferSize int) (*header, error) {
	h := &header{
		buffer: make([]byte, bufferSize),
	}
	h.waveHdr = &wavehdr{
		lpData:         uintptr(unsafe.Pointer(&h.buffer[0])),
		dwBufferLength: uint32(bufferSize),
	}
	if err := waveOutPrepareHeader(waveOut, h.waveHdr); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *header) Write(waveOut uintptr, data []byte) error {
	if len(data) != len(h.buffer) {
		return errors.New("oto: len(data) must equal to len(h.buffer)")
	}
	copy(h.buffer, data)
	if err := waveOutWrite(waveOut, h.waveHdr); err != nil {
		return err
	}
	return nil
}

type driver struct {
	out            uintptr
	headers        []*header
	tmp            []byte
	bufferSize     int
	bytesPerSecond int
}

func newDriver(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes int) (*driver, error) {
	numBlockAlign := channelNum * bitDepthInBytes
	f := &waveformatex{
		wFormatTag:      waveFormatPCM,
		nChannels:       uint16(channelNum),
		nSamplesPerSec:  uint32(sampleRate),
		nAvgBytesPerSec: uint32(sampleRate * numBlockAlign),
		wBitsPerSample:  uint16(bitDepthInBytes * 8),
		nBlockAlign:     uint16(numBlockAlign),
	}
	w, err := waveOutOpen(f)
	if err != nil {
		return nil, err
	}

	const numBufs = 2
	p := &driver{
		out:            w,
		headers:        make([]*header, numBufs),
		bufferSize:     bufferSizeInBytes,
		bytesPerSecond: sampleRate * channelNum * bitDepthInBytes,
	}
	runtime.SetFinalizer(p, (*driver).Close)
	for i := range p.headers {
		var err error
		p.headers[i], err = newHeader(w, p.bufferSize)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *driver) TryWrite(data []byte) (int, error) {
	n := min(len(data), max(0, p.bufferSize-len(p.tmp)))
	p.tmp = append(p.tmp, data[:n]...)
	if len(p.tmp) < p.bufferSize {
		return n, nil
	}

	var headerToWrite *header
	for _, h := range p.headers {
		// TODO: Need to check WHDR_DONE?
		if h.waveHdr.dwFlags&whdrInqueue == 0 {
			headerToWrite = h
			break
		}
	}
	if headerToWrite == nil {
		return n, nil
	}

	if err := headerToWrite.Write(p.out, p.tmp); err != nil {
		// This error can happen when e.g. a new HDMI connection is detected (#51).
		const errorNotFound = 1168
		werr := err.(*winmmError)
		if werr.fname == "waveOutWrite" && werr.errno == errorNotFound {
			return 0, nil
		}
		return 0, err
	}

	p.tmp = nil
	return n, nil
}

func (p *driver) Close() error {
	runtime.SetFinalizer(p, nil)
	// TODO: Call waveOutUnprepareHeader here
	if err := waveOutClose(p.out); err != nil {
		return err
	}
	return nil
}

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
