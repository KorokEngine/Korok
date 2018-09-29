package input

import (
	"sync"
	"korok.io/korok/math/f32"
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

	// 触摸/鼠标, 最多支持10个手指头同时触摸
	// 通常情况下，active < 1
	active int
	pointerButton [10]button
	pointers [10]PointerInput
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

func (in *InputSystem) AnyKeyChanged() bool {
	if in.dirty.used > 0 {
		return true
	}
	return false
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
func (in *InputSystem) AdvanceFrame() {
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
	for i := 0; i <= in.active; i++ {
		in.pointerButton[i].Reset()
	}
}

// 更新 key 的状态
func (in *InputSystem) SetKeyEvent(key int, pressed bool) {
	in.mutex.Lock()
	in.dirty.Put(Key(key), pressed)
	in.mutex.Unlock()
}

// 更新 Mouse/Touch 状态
func (in *InputSystem) SetPointerEvent(key int, pressed bool, x, y float32) {
	if key != -1000 {
		in.mutex.Lock()
		in.pointerButton[key].Update(pressed)
		in.pointers[key].MousePos = f32.Vec2{x, y}
		if key > in.active {
			in.active = key
		}
		in.mutex.Unlock()
	} else {
		// 如果是鼠标总是记录在 0 的位置
		// 如果是手指... 这就尴尬了..需要特殊处理
		in.pointers[0].MousePos = f32.Vec2{x, y}
	}
}

type Key int

type KeyBind struct {
	key Key
	btn *button
}

// short API
func Button(name string) button {
	return *Input.Button(name)
}

func AnyKeyChanged() bool {
	return Input.AnyKeyChanged()
}

func RegisterButton(name string, keys...Key) {
	Input.RegisterButton(name, keys...)
}

func PointerButton(pb KeyPoint) button {
	return Input.pointerButton[pb]
}

func PointerPosition(pb KeyPoint) PointerInput {
	return Input.pointers[pb]
}

// Touch event
func Touch(fi FingerId) (btn button, pos, delta f32.Vec2) {
	btn = Input.pointerButton[fi]
	p := Input.pointers[fi]
	pos, delta = p.MousePos, p.MouseDelta
	return
}

// Mouse event
func Mouse(key int) (btn *button, pos, delta f32.Vec2) {
	btn = &Input.pointerButton[key]
	p := Input.pointers[key]
	pos, delta = p.MousePos, p.MouseDelta
	return
}

var Input *InputSystem

// mouse or finger button
type KeyPoint int
const (
	KeyPoint1 KeyPoint = iota
	KeyPoint2
	KeyPoint3
	KeyPoint4
	KeyPoint5
	KeyPoint6
	KeyPoint7
	KeyPoint8
	KeyPoint9
	KeyPointX
)

