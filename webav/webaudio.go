package webav

import (
	"errors"

	"syscall/js"
)

type AudioContextAttributes struct {
	// LatencyHint: 这个参数表示了重放的类型, 参数是播放效果和资源消耗的一种权衡。
	// 可接受的值有 "balanced", "interactive" 和"playback"，
	// 默认值为 "interactive"。意思是 "平衡音频输出延迟和资源消耗", "提供最小的音频输出延迟最好没有干扰"和 "对比音频输出延迟，优先重放不被中断"。
	// 我们也可以用一个双精度的值来定义一个秒级的延迟数值做到更精确的控制。
	LatencyHint string
}

func ms2mi(x map[string]string) map[string]interface{} {
	y := make(map[string]interface{})
	for k, v := range x {
		y[k] = v
	}
	return y
}

// AudioDefaultAttributes 返回默认属性
func AudioDefaultAttributes() *AudioContextAttributes {
	return &AudioContextAttributes{"interactive"}
}

type AudioContext struct {
	js.Value
}

// AudioNewContext 获取音频对象
func AudioNewContext(ca *AudioContextAttributes) (*AudioContext, error) {

	if ca == nil {
		ca = AudioDefaultAttributes()
	}
	attrs := map[string]string{
		"latencyHint": ca.LatencyHint,
	}

	audioContext := js.Global().Get("AudioContext")
	if audioContext == js.Undefined() || audioContext == js.Null() {
		audioContext = js.Global().Get("webkitAudioContext")
	}
	if audioContext == js.Undefined() || audioContext == js.Null() {
		return nil, errors.New("Your browser doesn't appear to support Web Audio API.")
	}

	ac := audioContext.New(ms2mi(attrs))
	if ac == js.Null() {
		return nil, errors.New("Creating a Web Audio API context has failed.")
	}
	ctx := new(AudioContext)
	ctx.Value = ac
	return ctx, nil
}

func (c *AudioContext) Close() {
	c.Call("close")
}

func (c *AudioContext) CreateBuffer(numOfChannels, length, sampleRate int) js.Value {
	return c.Call("createBuffer", numOfChannels, length, sampleRate)
}

type AudioBufferSource struct {
	js.Value
}

func (c *AudioContext) CreateBufferSource() *AudioBufferSource {
	jsv := c.Call("createBufferSource")
	abs := new(AudioBufferSource)
	abs.Value = jsv
	return abs
}

func (c *AudioBufferSource) Start() js.Value {
	return c.Call("start")
}

func (c *AudioBufferSource) Stop() js.Value {
	return c.Call("stop")
}
