package inputs


// JustPressed/JustReleased 表示按键事件
// IsDown 表示按键状态
//
// 0x0001 记录当前状态
// 0x0002 记录pressed状态
// 0x0004 记录release状态
type button struct {
	id uint16
	state uint16
}
var btnId uint16 = 0
var buttonCache = make([]button, 64)
func NewButton() (btn *button) {
	btn = &buttonCache[btnId]
	btnId ++
	return btn
}

func (btn button) JustPressed() bool {
	return (btn.state & 0x02) == 2
}

func (btn button) JustReleased() bool {
	return (btn.state & 0x04) == 4
}

func (btn button) Down() bool {
	return (btn.state & 0x01) == 1
}

func (btn *button) Update(down bool){
	d := (btn.state & 0x01) == 1
	btn.state = 0
	if down {
		btn.state |= 0x01
	}
	if d && !down {
		btn.state |= 0x04
	} else if !d && down {
		btn.state |= 0x02
	}
}

func (btn *button) Reset() {
	btn.state = btn.state & 0x01
}
