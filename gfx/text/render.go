package text

import (
	"github.com/go-gl/mathgl/mgl32"
	"fmt"
)

// Text 本质上是一个RenderComp， 它只是对构建复杂RenderComp过程的封装！！
// 所有的渲染都可以变成对mesh的封装的过程，
// 概念：Text 只用来构建出一个 RenderComp，不占用 XXXComp的命名空间，
// XXXComp 是最底层的基础概念，应该把所有 XXXComp的概念封装成接口，这样可以
// 避免暴露出 id 索引的概念，况且 id 是非常难以使用的.
// 所以在上层使用的时候都是以面向对象的方式来进行的，到底层之后才会出现 id, xxxComp, xxxSystem

// 但是 unity 中 XXComp 是暴露给用户的!! （错，unity中也没有 textComp 这个概念，而是 textmesh），另外如果不保存 textComp 的引用，便会丢点这个数据
// unity 中给一个物体添加 text-mesh-comp 之后就可以绘制了
// comp 必须对应到一个 System 不然 comp 就找不到了！！
// 暂时再RenderComp中添加一个interface指向所有的外部上层的生成器对象，比如 text, shape, particle, skeleton!
// 否则会造成奇怪的现象，Text 不见了，但是 text 还会被绘制！！


///// 以上构建出一个合理的 Text 定义！
///// 还需要给 RenderComp 想一个合适的名字 Sprite ？
///// 在整个外部空间，使用 gameobject ?

type LabelComp struct {
	Font *Font

	id uint32

	// Color
	Color mgl32.Vec3

	// text data
	TextData

	RuneCount int

	String string
	CharSpacing []float32
}

// 重点在于生成 gfx.renderComp 的 mesh ！
func NewText(font *Font) (t *LabelComp){
	t = new(LabelComp)
	t.Font = font
	t.TextData.tex = uint16(font.Texture)
	return t
}


func (t *LabelComp) SetScale(s float32) {
//	t.Scale = s
}

func (t *LabelComp) SetString(fs string, argv ...interface{}) {
	t.String = fmt.Sprintf(fs, argv...)

	// init ebo, vbo
	t.RuneCount     = len(t.String)

	// fill data
	t.fillData()
}

func (t *LabelComp) RenderData() *TextData {
	return &t.TextData
}

// fill vbo/ebo with the string
//
//		+----------+
//		| . 	   |
//      |   .	   |
//		|     .    |
// 		|		.  |
// 		+----------+
// 1 * 1 quad for each char
func (t *LabelComp) fillData() {
	var xOffset float32
	var yOffset float32

	chars := make([]TextQuad, len(t.String))
	t.TextData.Chars = chars

	for i, r := range t.String {
		if glyph := t.Font.config.Find(r); glyph != nil {
			advance := float32(glyph.Advance)
			vw := glyph.Width
			vh := glyph.Height

			min, max := glyph.GetTexturePosition(t.Font)
			char := &chars[i]

			char.xOffset = xOffset
			char.yOffset = yOffset
			char.w, char.h = float32(vw), float32(vh)
			char.region.X1, char.region.Y1 = min.X, min.Y
			char.region.X2, char.region.Y2 = max.X, max.Y

			// left to right shit
			xOffset += advance
			yOffset += 0
		}
	}
}
