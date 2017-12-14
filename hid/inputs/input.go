package inputs

import (
	"sync"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// 记录一帧之内的按键，一帧时间做多支持同时按6个按键
type SparseMap struct {
	keys [6]Key
	stat [6]bool
	used int
}

func (m *SparseMap) Put(k Key, st bool) {
	m.keys[m.used] = k
	m.stat[m.used] = st
	m.used ++
}

func (m *SparseMap) Clear() {
	m.used = 0
}

func (m *SparseMap) Get(k Key) (st, ok bool) {
	for i := 0; i < m.used; i++ {
		if m.keys[i] == k {
			st, ok = m.stat[i], true
			break
		}
	}
	return
}

type InputSystem struct {
	buttons map[string]*button
	axes	map[string]*VAxis

	// 记录每帧的按键状态
	// 无论是用数组还是哈希，这里的实现总之要达到快速
	// 查询一个按键的状态的效果
	dirty SparseMap
	mutex   sync.RWMutex

	// 按照button排序，这样同一个Button的绑定按键是
	// 是连续的。
	binds []KeyBind
}

func NewInputSystem() *InputSystem {
	in := &InputSystem{
		buttons:make(map[string]*button),
		axes:make(map[string]*VAxis),
	}
	Input = in
	return in
}

/// 查询虚拟按键的状态
func (in *InputSystem) Button(name string) *button {
	return in.buttons[name]
}

/// 将物理按键映射到虚拟按键
func (in *InputSystem) RegisterButton(name string, keys ...Key) {
	btn := NewButton()
	in.buttons[name] = btn
	for _, k := range keys {
		in.binds = append(in.binds, KeyBind{k, btn})
	}
	// sort binds!!
}

// 更新 Button 状态....
// TODO 此处的输入状态，更新有bug！！
func (in *InputSystem) Frame() {
	if n, dirty := len(in.binds), in.dirty.used; n > 0 && dirty > 0 {
		var st, ok bool
		var pr *button

		for _, bd := range in.binds {
			if s, o := in.dirty.Get(bd.key); o {
				st = st || s
				ok = ok || o
			}

			if pr != bd.btn {
				if ok {
					bd.btn.Update(st)
				}
				st, ok = false, false
			}

			pr = bd.btn
		}
	}
}

func (in *InputSystem) Reset() {
	// clear dirty map!!
	in.mutex.Lock()
	in.dirty.Clear()
	in.mutex.Unlock()
	// reset button state
	for _, v := range in.buttons {
		v.Reset()
	}
}

// 更新 key 的状态
func (in *InputSystem) SetKeyEvent(key int, pressed bool) {
	in.mutex.Lock()
	in.dirty.Put(Key(key), pressed)
	in.mutex.Unlock()
}

type Key int

type KeyBind struct {
	key Key
	btn *button
}

/// short API
func Button(name string) *button {
	return Input.Button(name)
}

func RegisterButton(name string, keys...Key) {
	Input.RegisterButton(name, keys...)
}

var Input *InputSystem

const (
	Grave = Key(glfw.KeyGraveAccent)
	Dash = Key(glfw.KeyMinus)
	Apostrophe = Key(glfw.KeyApostrophe)
	Semicolon = Key(glfw.KeySemicolon)
	Equals = Key(glfw.KeyEqual)
	Comma = Key(glfw.KeyComma)
	Period = Key(glfw.KeyPeriod)
	Slash = Key(glfw.KeySlash)
	Backslash = Key(glfw.KeyBackslash)
	Backspace = Key(glfw.KeyBackspace)
	Tab = Key(glfw.KeyTab)
	CapsLock = Key(glfw.KeyCapsLock)
	Space = Key(glfw.KeySpace)
	Enter = Key(glfw.KeyEnter)
	Escape = Key(glfw.KeyEscape)
	Insert = Key(glfw.KeyInsert)
	PrintScreen = Key(glfw.KeyPrintScreen)
	Delete = Key(glfw.KeyDelete)
	PageUp = Key(glfw.KeyPageUp)
	PageDown = Key(glfw.KeyPageDown)
	Home = Key(glfw.KeyHome)
	End = Key(glfw.KeyEnd)
	Pause = Key(glfw.KeyPause)
	ScrollLock = Key(glfw.KeyScrollLock)
	ArrowLeft = Key(glfw.KeyLeft)
	ArrowRight = Key(glfw.KeyRight)
	ArrowDown = Key(glfw.KeyDown)
	ArrowUp = Key(glfw.KeyUp)
	LeftBracket = Key(glfw.KeyLeftBracket)
	LeftShift = Key(glfw.KeyLeftShift)
	LeftControl = Key(glfw.KeyLeftControl)
	LeftSuper = Key(glfw.KeyLeftSuper)
	LeftAlt = Key(glfw.KeyLeftAlt)
	RightBracket = Key(glfw.KeyRightBracket)
	RightShift = Key(glfw.KeyRightShift)
	RightControl = Key(glfw.KeyRightControl)
	RightSuper = Key(glfw.KeyRightSuper)
	RightAlt = Key(glfw.KeyRightAlt)
	Zero = Key(glfw.Key0)
	One = Key(glfw.Key1)
	Two = Key(glfw.Key2)
	Three = Key(glfw.Key3)
	Four = Key(glfw.Key4)
	Five = Key(glfw.Key5)
	Six = Key(glfw.Key6)
	Seven = Key(glfw.Key7)
	Eight = Key(glfw.Key8)
	Nine = Key(glfw.Key9)
	F1 = Key(glfw.KeyF1)
	F2 = Key(glfw.KeyF2)
	F3 = Key(glfw.KeyF3)
	F4 = Key(glfw.KeyF4)
	F5 = Key(glfw.KeyF5)
	F6 = Key(glfw.KeyF6)
	F7 = Key(glfw.KeyF7)
	F8 = Key(glfw.KeyF8)
	F9 = Key(glfw.KeyF9)
	F10 = Key(glfw.KeyF10)
	F11 = Key(glfw.KeyF11)
	F12 = Key(glfw.KeyF12)
	A = Key(glfw.KeyA)
	B = Key(glfw.KeyB)
	C = Key(glfw.KeyC)
	D = Key(glfw.KeyD)
	E = Key(glfw.KeyE)
	F = Key(glfw.KeyF)
	G = Key(glfw.KeyG)
	H = Key(glfw.KeyH)
	I = Key(glfw.KeyI)
	J = Key(glfw.KeyJ)
	K = Key(glfw.KeyK)
	L = Key(glfw.KeyL)
	M = Key(glfw.KeyM)
	N = Key(glfw.KeyN)
	O = Key(glfw.KeyO)
	P = Key(glfw.KeyP)
	Q = Key(glfw.KeyQ)
	R = Key(glfw.KeyR)
	S = Key(glfw.KeyS)
	T = Key(glfw.KeyT)
	U = Key(glfw.KeyU)
	V = Key(glfw.KeyV)
	W = Key(glfw.KeyW)
	X = Key(glfw.KeyX)
	Y = Key(glfw.KeyY)
	Z = Key(glfw.KeyZ)
	NumLock = Key(glfw.KeyNumLock)
	NumMultiply = Key(glfw.KeyKPMultiply)
	NumDivide = Key(glfw.KeyKPDivide)
	NumAdd = Key(glfw.KeyKPAdd)
	NumSubtract = Key(glfw.KeyKPSubtract)
	NumZero = Key(glfw.KeyKP0)
	NumOne = Key(glfw.KeyKP1)
	NumTwo = Key(glfw.KeyKP2)
	NumThree = Key(glfw.KeyKP3)
	NumFour = Key(glfw.KeyKP4)
	NumFive = Key(glfw.KeyKP5)
	NumSix = Key(glfw.KeyKP6)
	NumSeven = Key(glfw.KeyKP7)
	NumEight = Key(glfw.KeyKP8)
	NumNine = Key(glfw.KeyKP9)
	NumDecimal = Key(glfw.KeyKPDecimal)
	NumEnter = Key(glfw.KeyKPEnter)
)

