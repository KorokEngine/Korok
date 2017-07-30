package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

type Drawable struct {
	vao, vbo, tex uint32
}

type MeshRender struct {
	shader *Shader
}

func (mr *MeshRender) Draw(d *Drawable) {
	mr.shader.Use()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindTexture(gl.TEXTURE_2D, d.tex)

	gl.BindVertexArray(d.vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)

	gl.Disable(gl.BLEND)
}
