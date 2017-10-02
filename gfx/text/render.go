package text

import (
	"korok/gfx"

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

	// text mesh
	mesh gfx.Mesh

	RuneCount int

	String string
	CharSpacing []float32
}

// 重点在于生成 gfx.renderComp 的 mesh ！
func NewText(font *Font) (t *LabelComp){
	t = new(LabelComp)
	t.Font = font
	t.mesh.SetRawTexture(font.Texture)
	return t
}

func (t *LabelComp) Release() {
	t.mesh.Delete()
}

func (t *LabelComp) SetScale(s float32) {
//	t.Scale = s
}

func (t *LabelComp) SetString(fs string, argv ...interface{}) {
	t.String = fmt.Sprintf(fs, argv...)

	// init ebo, vbo
	t.RuneCount     = len(t.String)

	vboData, eboData := t.fillData()
	t.mesh.SetVertex(vboData)
	t.mesh.SetIndex(eboData)
	t.mesh.Setup()
}

func (t *LabelComp) Mesh() *gfx.Mesh {
	return &t.mesh
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
func (t *LabelComp) fillData() ([]float32, []int32) {
	var xOffset float32
	var yOffset float32

	vboIndexCount := len(t.String) * 4 * 4  // len * 4 * <x,y,u,v>
	eboIndexCount := len(t.String) * 6 	 // len * 6

	vboData := make([]float32, vboIndexCount)
	eboData := make([]int32, eboIndexCount)

	for i, r := range t.String {
		if glyph := t.Font.config.Find(r); glyph != nil {
			advance := float32(glyph.Advance)
			vw := glyph.Width
			vh := glyph.Height

			min, max := glyph.GetTexturePosition(t.Font)

			/// step = 16(4*4)
			vi := i * 16

			// 0- 3 互换, 1 - 2 互换
			// index (0, 0) <x,y,u,v>
			vboData[vi+0] = 0 + xOffset
			vboData[vi+1] = 0 + yOffset
			vboData[vi+2] = min.X
			vboData[vi+3] = max.Y


			// index (1,0) <x,y,u,v>
			vi += 4
			vboData[vi+0] = float32(vw) + xOffset
			vboData[vi+1] = 0 + yOffset
			vboData[vi+2] = max.X
			vboData[vi+3] = max.Y

			// index(1,1) <x,y,u,v>
			vi += 4
			vboData[vi+0] = float32(vw) + xOffset
			vboData[vi+1] = float32(vh) + yOffset
			vboData[vi+2] = max.X
			vboData[vi+3] = min.Y

			// index(0, 1) <x,y,u,v>
			vi += 4
			vboData[vi+0] = 0 + xOffset
			vboData[vi+1] = float32(vh) + yOffset
			vboData[vi+2] = min.X
			vboData[vi+3] = min.Y

			/// step = 6
			ei := i * 6
			bi := int32(i * 4)

			eboData[ei+0] = bi + 1
			eboData[ei+1] = bi + 2
			eboData[ei+2] = bi + 3

			eboData[ei+3] = bi + 0
			eboData[ei+4] = bi + 1
			eboData[ei+5] = bi + 3

			// left to right shit
			xOffset += advance
			yOffset += 0
		}
	}

	return vboData, eboData
}
