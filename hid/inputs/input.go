package inputs

var buttons map[string]*VButton
var axes	map[string]*VAxis

func init() {
	buttons = make(map[string]*VButton)
	axes    = make(map[string]*VAxis)

	// default
	RegisterButton("Action", 1, 2)

}

/// 查询虚拟按键的状态
func Button(name string) *VButton {
	return buttons[name]
}

/// 将物理按键映射到虚拟按键
func RegisterButton(name string, keys ...Key) {
	buttons[name] = &VButton{
		Name:name,
		Keys:keys,
	}
}

/// 查询虚拟摇杆的状态
func Axis(name string) *VAxis {
	return axes[name]
}

/// TODO 如何抽象虚拟遥感？？
func RegisterAxis(name string) {

}




type Key int

type KeyState int

