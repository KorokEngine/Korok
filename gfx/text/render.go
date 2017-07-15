package text

import (
	"korok/gfx"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v3.2-core/gl"

	"github.com/golang/freetype/truetype"
	"github.com/golang/freetype"

	"io/ioutil"
	"log"
	"fmt"
)

type LabelComp struct {
	Font *Font

	// Position
	Position mgl32.Vec2

	// Scale
	Scale float32

	// Color
	Color mgl32.Vec3

	//
	vao uint32
	vbo uint32
	ebo uint32

	vboData []float32
	vboIndexCount int
	eboData []int32
	eboIndexCount int

	RuneCount int

	vertices []float32

	String string
	CharSpacing []float32
}

func NewText(font *Font) (t *LabelComp){
	t = new(LabelComp)

	t.Font = font

	gl.GenVertexArrays(1, &t.vao)
	gl.GenBuffers(1, &t.vbo)
	gl.GenBuffers(1, &t.ebo)

	// bind fixture and attribute
	gl.BindVertexArray(t.vao)

	gl.BindTexture(gl.TEXTURE_2D, t.Font.Texture)
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)
	return
}

func (t *LabelComp) Release() {
	gl.DeleteBuffers(1, &t.vbo)
	gl.DeleteBuffers(1, &t.ebo)
	gl.DeleteVertexArrays(1, &t.vao)
}

func (t *LabelComp) SetScale(s float32) {
	t.Scale = s
}

func (t *LabelComp) SetString(fs string, argv ...interface{}) {
	t.String = fmt.Sprintf(fs, argv...)

	// init ebo, vbo
	glfloat_size := int(4)

	t.vboIndexCount = len(t.String) * 4 * 4  // len * 4 * <x,y,u,v>
	t.eboIndexCount = len(t.String) * 6 	 // len * 6
	t.RuneCount     = len(t.String)

	t.vboData = make([]float32, t.vboIndexCount, t.vboIndexCount)
	t.eboData = make([]int32, t.eboIndexCount, t.eboIndexCount)

	// generate the basic vbo
	t.fillData()

	// bind vbo/vao
	gl.BindVertexArray(t.vao)

	gl.BufferData(gl.ARRAY_BUFFER, glfloat_size*t.vboIndexCount, gl.Ptr(t.vboData), gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, glfloat_size*t.eboIndexCount, gl.Ptr(t.eboData), gl.DYNAMIC_DRAW)

	gl.BindVertexArray(0)

	fmt.Println(t.vboData)
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

	for i, r := range t.String {
		if glyph := t.Font.config.Find(r); glyph != nil {
			advance := float32(glyph.Advance)
			vw := glyph.Width
			vh := glyph.Height

			min, max := glyph.GetTexturePosition(t.Font)

			/// step = 16(4*4)
			vi := i * 16

			// TODO 纹理是倒置，临时通过坐标映射解决
			// 0- 3 互换, 1 - 2 互换
			// index (0, 0) <x,y,u,v>
			t.vboData[vi+0] = 0 + xOffset
			t.vboData[vi+1] = 0 + yOffset
			t.vboData[vi+2] = min.X
			t.vboData[vi+3] = max.Y


			// index (1,0) <x,y,u,v>
			vi += 4
			t.vboData[vi+0] = float32(vw) + xOffset
			t.vboData[vi+1] = 0 + yOffset
			t.vboData[vi+2] = max.X
			t.vboData[vi+3] = max.Y

			// index(1,1) <x,y,u,v>
			vi += 4
			t.vboData[vi+0] = float32(vw) + xOffset
			t.vboData[vi+1] = float32(vh) + yOffset
			t.vboData[vi+2] = max.X
			t.vboData[vi+3] = min.Y

			// index(0, 1) <x,y,u,v>
			vi += 4
			t.vboData[vi+0] = 0 + xOffset
			t.vboData[vi+1] = float32(vh) + yOffset
			t.vboData[vi+2] = min.X
			t.vboData[vi+3] = min.Y

			/// step = 6
			ei := i * 6
			bi := int32(i * 4)

			t.eboData[ei+0] = bi + 1
			t.eboData[ei+1] = bi + 2
			t.eboData[ei+2] = bi + 3

			t.eboData[ei+3] = bi + 0
			t.eboData[ei+4] = bi + 1
			t.eboData[ei+5] = bi + 3

			// left to right shit
			xOffset += advance
			yOffset += 0
		}
	}
}

type Renderer struct {
	shader *gfx.Shader
	TTF *truetype.Font
}

func NewTextRenderer(shader *gfx.Shader) *Renderer {
	renderer := &Renderer{}
	renderer.shader = shader
	return renderer
}

func (renderer *Renderer) Load(file string, size uint32)  {
	ttfBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return
	}

	ttf, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		log.Println(err)
		return
	}
	renderer.TTF = ttf
}

// Render current text component!!
func (renderer *Renderer) RenderText(comp *LabelComp) {
	renderer.shader.Use()
	renderer.shader.SetVector3f("model\x00", comp.Position[0], comp.Position[1], 10)
	renderer.shader.SetVector3f("textColor\x00", comp.Color[0], comp.Color[1], comp.Color[2])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, comp.Font.Texture)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BindVertexArray(comp.vao)
	gl.DrawElements(gl.TRIANGLES, int32(comp.eboIndexCount), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
	gl.Disable(gl.BLEND)
}
