package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"korok/gfx/bk"
)

/// A Sprite Batch TypeRender
type BatchRender struct {
	BatchContext

	stateFlags uint64
	rgba       uint32

	// shader program
	program uint16

	// uniform handle
	umh_PJ uint16 	// Projection
	umh_S0 uint16 	// Sampler0
}

func NewBatchRender(vsh, fsh string) *BatchRender {
	br := new(BatchRender)

	br.stateFlags |= bk.ST_BLEND.ADDITIVE

	// setup shader
	if id, sh := bk.R.AllocShader(vsh, fsh); id != bk.InvalidId {
		br.program = id

		// setup attribute
		sh.AddAttributeBinding("xyuv", 0, P4C4[0])
		sh.AddAttributeBinding("rgba", 0, P4C4[1])

		// setup uniform
		br.umh_PJ, _ = bk.R.AllocUniform(id, "proj\x00", bk.UniformMat4, 1)

		// TODO
		sh.SetInteger("tex\x00", 0)
		gl.BindFragDataLocation(sh.Program, 0, gl.Str("outputColor\x00"))

	}
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
	quad := d.(*Quad)
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

