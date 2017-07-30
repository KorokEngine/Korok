package anim

import (
	"korok/gfx"
	"korok/assets"


	"github.com/go-gl/gl/v3.2-core/gl"
	"korok/anim/spine"
)

type SkeletonRender struct {
	shader *gfx.Shader
	vao, vbo, ebo uint32
	tex uint32
	buffer []float32
}

func LoadSpine() {
}

func NewSkeletonRender(shader *gfx.Shader) *SkeletonRender {
	sr := new(SkeletonRender)
	sr.shader = shader
	sr.tex = assets.GetTexture("assets/spine/alien.png").Id
	sr.buffer = make([]float32, 16)

	gl.GenVertexArrays(1, &sr.vao)
	gl.GenBuffers(1, &sr.vbo)
	gl.GenBuffers(1, &sr.ebo)

	shader.Use()

	gl.BindVertexArray(sr.vao)

	//// VAO 缺省的VBO和EBO

	var vertices = []float32{
		// Pos
		0.0, 0.0, 0, 0,
		200.0, 0.0, 1, 0,
		200.0, 200.0,  1, 1,
		0.0, 200.0, 0, 1,
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, sr.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) * 4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)

	var indices = []int32 {
		0, 1, 2,
		0, 2, 3,
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, sr.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.DYNAMIC_DRAW)


	// VAO 捕获顶点属性配置
	vertAttrib := uint32(gl.GetAttribLocation(sr.shader.Program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(sr.shader.Program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	gl.BindVertexArray(0)

	return sr
}

// 每次更新... 不需要VAO吧 - NO!!! 每次都必须包围VAO才能绘制
func (sr*SkeletonRender) draw(vertex []float32, uv []float32) {
	sr.shader.Use()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)


	// update vertex
	for i := 0;  i < 4; i++ {
		sr.buffer[i*4 + 0] = vertex[i*2 + 0]
		sr.buffer[i*4 + 1] = vertex[i*2 + 1]
		sr.buffer[i*4 + 2] = uv[i*2 + 0]
		sr.buffer[i*4 + 3] = uv[i*2 + 1]
	}

	// draw ..

	gl.BindVertexArray(sr.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, sr.vbo)

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(sr.buffer) * 4, gl.Ptr(sr.buffer))

	gl.BindTexture(gl.TEXTURE_2D, sr.tex)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)


	gl.Disable(gl.BLEND)
}

func (sr*SkeletonRender)Draw(skeleton *spine.Skeleton) {
	for _, slot := range skeleton.Slots {
		//
		if attachment, ok := slot.Attachment.(*spine.RegionAttachment); ok {
			// 计算得到插件坐标
			vert := attachment.Update(slot)
			//fmt.Println("index:", i, " name:", attachment.Name(), " update:", vert)
			sr.draw(vert[0:], attachment.Uvs[0:])
		} else {
			// fmt.Println("index:",i , " is null")
		}
	}
}

////
//var vertices = [8]float32{
//	// Pos
//	200.0, 200.0,
//	0.0, 200.0,
//	0.0, 0.0,
//	200.0, 0.0,
//}

//var uvs = [8]float32 {
//	0, 0,
//	1, 0,
//	1, 1,
//	0, 1,
//}
//var uvs = [8]float32 {
//	0.997, 0.998,
//	0.991, 0.998,
//	0.991, 0.990,
//	0.997, 0.990,
//}

//attach := g.skeleton.Slots[5].Attachment.(*spine.RegionAttachment)
//
//fmt.Println("slot name:", attach.Name())
//fmt.Println("uvs:", attach.Uvs)
//fmt.Println("vert:", attach.Update(g.skeleton.Slots[3]))
//
//vertices := attach.Update(g.skeleton.Slots[5])
//uvs      := attach.Uvs
//g.skRender.draw(vertices[0:], uvs[0:])
//g.skRender.Draw(g.skeleton)
//_, bone := g.skeleton.FindBone("back-thigh")
//if bone != nil {
//	//fmt.Println("bone m00:", bone.M00, "m01:", bone.M01, " m10:", bone.M10, " m11:", bone.M11)
//	fmt.Println("bone trans:", bone.X, bone.Y, " scale:", bone.ScaleX, bone.ScaleY, " rotation:", bone.Rotation)
//}
