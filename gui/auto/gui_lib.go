package auto

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	g "korok.io/korok/gui"
)

type LayoutMan struct {
	*g.Context
	Options
	fallback layout

	layouts map[string]*layout
	depth int
	current *layout

	// sqNum should be same for  layout and drawing
	sqNum int
}

func (lm *LayoutMan) initialize() {
	lm.fallback.Initialize(gContext.Theme)
	lm.current = &lm.fallback
	lm.layouts = make(map[string]*layout)
}

func (lm *LayoutMan) Text(id g.ID, text string, style *g.TextStyle, opt *Options)  *Element {
	var (
		elem, ready = lm.BeginElement(id, opt)
		size f32.Vec2
		font = style.Font
	)
	if font == nil {
		font = lm.Context.Theme.Font
	}

	// draw text 最好返回最新的大小..
	if ready {
		size = lm.Context.DrawText(&elem.Rect, text, style)
	} else {
		size = lm.Context.CalcTextSize(text, 0, font, style.Size)
	}

	elem.Rect.W = size[0]
	elem.Rect.H = size[1]

	lm.EndElement(elem)
	return nil
}

// Widgets: InputEditor
func (lm *LayoutMan) InputText(hint string, lyt layout, style *g.InputStyle) {

}

// Widget: Image
func (lm *LayoutMan) Image(id g.ID, tex gfx.Tex2D, style *g.ImageStyle, opt *Options) {
	var (
		elem, ready = lm.BeginElement(id, opt)
	)

	if ready {
		lm.Context.DrawImage(&elem.Rect, tex, style)
	} else {
		if opt != nil {
			elem.W = opt.W
			elem.H = opt.H
		}
	}

	lm.EndElement(elem)
}

// Widget: Button
func (lm *LayoutMan) Button(id g.ID, text string, style *g.ButtonStyle, opt *Options) (event g.EventType) {
	var (
		elem, ready = lm.BeginElement(id, opt)
	)

	if ready {
		lm.Context.Button(id, &elem.Rect, text, style)
	} else {
		font := style.Font
		if font == nil {
			font = lm.Context.Theme.Font
		}
		size := style.TextStyle.Size
		textSize := lm.CalcTextSize(text, 0, font, size)
		extW := style.Padding.Left+style.Padding.Right
		extH := style.Padding.Top+style.Padding.Bottom
		elem.W, elem.H = textSize[0]+extW, textSize[1]+extH
	}
	lm.EndElement(elem)
	return
}


func (lm *LayoutMan) renderTextClipped(text string, bb *g.Rect, style *g.TextStyle) {
	x, y := g.Gui2Game(bb.X, bb.Y)
	font := lm.Theme.Text.Font
	if bb.W == 0 {
		lm.DrawList.AddText(f32.Vec2{x, y}, text, font, 12, 0xFF000000, 0)
	} else {
		lm.DrawList.AddText(f32.Vec2{x, y}, text, font, 12, 0xFF000000, bb.W)
	}
}

func (lm *LayoutMan) ImageButton(id g.ID, normal, pressed gfx.Tex2D, style *g.ImageButtonStyle, opt *Options) ( event g.EventType) {
	var (
		elem, ready = lm.BeginElement(id, opt)
		bb = &elem.Rect
	)
	if ready {
		event = lm.ClickEvent(id, bb)
		var tex gfx.Tex2D
		if event & g.EventDown != 0 {
			tex = pressed
		} else {
			tex = normal
		}
		lm.DrawImage(bb, tex, &style.ImageStyle)
	} else {
		if opt != nil {
			elem.W = opt.W
			elem.H = opt.H
		}
	}
	lm.EndElement(elem)
	return
}

// Slider 需要设定一些自定义的属性，目前没有想好如何实现，先把逻辑实现了
// 用两种颜色来绘制
func (lm *LayoutMan) Slider(id g.ID, value *float32, style *g.SliderStyle, opt *Options) (e g.EventType){
	var (
		elem, ready = lm.BeginElement(id, opt)
		bb = &elem.Rect
	)

	if ready {
		// 说明滑动了，那么应该使用最新的值，而不是传入的值
		if v, event := lm.Context.CheckSlider(id, bb); event & g.EventDragging != 0 {
			*value = v
			e = event
		}

		lm.DrawRect(bb, style.Bar, 5)
		lm.DrawCircle(bb.X+bb.W*(*value), bb.Y+bb.H/2, 10, style.Knob)
	} else {
		// 设置默认的宽高
		if elem.W == 0 {
			elem.W = 120
		}
		if elem.H == 0 {
			elem.H = 10
		}
	}

	lm.EndElement(elem)
	return
}

func (lm *LayoutMan) DefineLayout(name string, xt ViewType) {
	if l, ok := lm.layouts[name]; ok {
		lm.current = l
	} else {
		l := &layout{}; l.Initialize(gContext.Theme)
		lm.layouts[name] = l
		lm.current = l
	}
}

func (lm *LayoutMan) Clear(names ...string) {
	if size := len(names); size > 0 {
		for i := 0; i < size; i++ {
			delete(lm.layouts, names[i])
		}
	} else {
		lm.layouts = make(map[string]*layout)
	}
}

// 计算单个UI元素
// 如果有大小则记录出偏移和Margin
// 否则只返回元素
func (lm *LayoutMan) BeginElement(id g.ID, opt *Options) (elem *Element, ok bool){
	return lm.current.BeginElement(id, opt)
}

// 结束绘制, 每绘制完一个元素都要偏移一下光标
func (lm *LayoutMan) EndElement(elem *Element) {
	lm.current.EndElement(elem)
}

func (lm *LayoutMan) BeginLayout(id g.ID, opt *Options, xtype LayoutType) {
	lm.depth ++
	if elem, ok := lm.current.BeginLayout(id, opt, xtype); ok {
		// debug-draw
		if g.DebugDraw {
			lm.DrawDebugBorder(elem.X, elem.Y, elem.W, elem.H, 0xFF00FF00)
		}
	}
}

func (lm *LayoutMan) EndLayout() {
	lm.depth --
	lm.current.EndLayout()
	if d := lm.depth; d == 0 {
		lm.current = &lm.fallback
	}
}
