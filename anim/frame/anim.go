package frame

import (
	"korok.io/korok/gfx"
)

// implement sprite-animation system

// Frames data of Sprite Animation
type Animation struct {
	Name string
	Start, Len int
	Loop bool
}

// Sprite Animation System
type SpriteEngine struct {
	// raw frames
	frames []gfx.Tex2D
	// raw animation
	data []Animation
	// mapping from name to index
	names map[string]int

	// sprite and animate table
	st *gfx.SpriteTable
	at *FlipbookTable
}

func NewEngine() *SpriteEngine {
	return &SpriteEngine{ names:make(map[string]int) }
}

func (eng *SpriteEngine) RequireTable(tables []interface{}) {
	for _, t := range tables {
		switch table := t.(type) {
		case *gfx.SpriteTable:
			eng.st = table
		case *FlipbookTable:
			eng.at = table
		}
	}
}

// 创建新的动画数据
// 现在 subText 还是指针，稍后会全部用 id 来索引。
// 动画资源全部存储在一个大的buffer里面，在外部使用索引引用即可.
// 采用这种设计，删除动画将会变得麻烦..
// 或者说无法删除动画，只能全部删除或者完全重新加载...
// 如何动画以组的形式存在，那么便可以避免很多问题
//
func (eng *SpriteEngine) NewAnimation(name string, frames []gfx.Tex2D, loop bool) {
	// copy frames
	start, size := len(eng.frames), len(frames)
	eng.frames = append(eng.frames, frames...)
	// new animation
	eng.data = append(eng.data, Animation{name, start, size, loop})
	// keep mapping
	eng.names[name] = len(eng.data)-1
}

// 返回动画定义 - 好像并没有太大的意义
func (eng *SpriteEngine) Animation(name string) (anim *Animation, seq []gfx.Tex2D) {
	if ii, ok := eng.names[name]; ok {
		anim = &eng.data[ii]
		seq  = eng.frames[anim.Start:anim.Start+anim.Len]
	}
	return
}

func (eng *SpriteEngine) Update(dt float32) {
	var (
		at, st = eng.at, eng.st
		anims  = at.comps[:at.index]
	)

	// update animation state
	for i := range anims {
		if am := &anims[i]; am.running {
			var (
				id    = eng.names[am.define]
				data  = eng.data[id]
			)
			am.gfi = data.Start+int(am.frameIndex)
			if am.dt += dt; am.dt > am.rate {
				am.ii = am.ii + 1
				am.dt = 0
				frame := am.ii% data.Len

				// frame end
				if frame == 0 {
					if am.loop && am.typ == PingPong {
						am.reverse = !am.reverse
						am.ii += 1 // skip one frame, or it'll repeat last frame
						frame += 1
					}
					if !am.loop {
						am.running = false
						break
					}
				}

				if am.reverse {
					frame = data.Len-frame-1
				}
				// update frame index
				am.lastFrameIndex = am.frameIndex
				am.frameIndex = uint16(frame)
				am.gfi = data.Start+frame
			}

			// update sprite-component
			comp := st.Comp(am.Entity)
			frame := eng.frames[am.gfi]
			comp.SetSprite(frame)
		}
	}
}
