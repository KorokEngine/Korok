package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"korok/gfx/bk"
)

/// A Sprite Batch TypeRender
type BatchRender struct {
	pipeline PipelineState
	BatchContext
	RenderContext
}

func NewBatchRender(shader bk.GLShader) *BatchRender {
	br := new(BatchRender)

	br.pipeline.BlendFunc = BF_Add
	br.pipeline.GLShader = shader
	shader.Use()

	shader.SetInteger("tex\x00", 0)
	gl.BindFragDataLocation(shader.Program, 0, gl.Str("outputColor\x00"))

	// vertex layout
	pos := VertexAttr {
		Size: 2,
		Type: gl.FLOAT,
		Normalized: false,
		Stride: 20,
		Offset: 0,
	}
	uv := VertexAttr {
		Size: 2,
		Type: gl.FLOAT,
		Normalized: false,
		Stride: 20,
		Offset: 8,
	}
	color := VertexAttr {
		Size: 4,
		Type: gl.UNSIGNED_BYTE,
		Normalized: false,
		Stride: 20,
		Offset: 16,
	}
	br.pipeline.VertexLayout = append(br.pipeline.VertexLayout, pos, uv, color)

	// uniform layout
	p := bk.Uniform{
		Data: 0, 		// index of uniform data
		Slot: shader.GetUniformLocation("projection\x00"), 		// slot in shader
		Type: bk.UniformMat4,
		Count: 1,
	}
	br.pipeline.UniformLayout = append(br.pipeline.UniformLayout, p)

	return br
}

/**
if batch.ready && batch.compatible {

}

对于 Batch Render 来说，是无法知道外面的 RenderComp 的排序状况的，
那么也无法知道目前传入的 RenderData 是应该同上一批 Batch 合并，还是
应该提交batch还是应该建立一个新的batch

 */
func (br *BatchRender) Draw(d RenderData, pos, scale mgl32.Vec2, rot float32) {
	quad := d.(Quad)
	// 计算顶点, scale, rot TODO
	vertex := quad.buf
	for i := range vertex {
		vertex[i].XY[0] += pos[0]
		vertex[i].XY[1] += pos[1]
	}

	if br.BatchContext.Ready() && br.BatchContext.Compatible() {
		// br.BatchContext.
	}

}

type BatchContext struct {
	B Batch

}

func (*BatchContext) Ready() bool {
	return false
}

func (*BatchContext) Compatible() bool {
	return false
}

func (*BatchContext) Begin() {
}

func (*BatchContext) Draw() {

}

func (*BatchContext) End() {

}

