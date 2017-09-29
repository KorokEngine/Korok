package bk

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"unsafe"
)

type RenderContext struct {
	// data
	R  *ResManager
	ub *UniformBuffer

	// draw list
	drawList []RenderDraw

	// draw state
	vao uint32
	vaoSupport bool

	// window rect
	wRect Rect

	backBufferFbo uint32
}

func (ctx *RenderContext) Init() {

}

func (ctx *RenderContext) Shutdown() {

}

func (ctx *RenderContext) Draw(sortKeys []uint64, sortValues []uint16, drawList []RenderDraw) {
	// 1. 绑定 VAO 和 FrameBuffer
	if defaultVao := ctx.vao; 0 != defaultVao {
		gl.BindVertexArray(defaultVao)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, ctx.backBufferFbo)

	// 2. 更新分辨率
	// ctx.updateResolution(&render.resolution)

	// 3. 初始化 per-draw state
	currentState := RenderDraw{}

	shaderId := InvalidId
	key := SortKey{}
	view := UINT16_MAX

	primIndex := uint8(uint64(0) >> ST.PT_SHIFT)
	prim := R_PrimInfo[primIndex]

	var viewHasScissor bool
	viewScissorRect := Rect{}
	viewScissorRect.clear()


	// 8. Render!!
	for item := range sortKeys{
		encodedKey := sortKeys[item]
		itemId 	   := sortValues[item]
		key.Decode(encodedKey)

		draw := drawList[itemId]

		// 1. 求取变化的状态位
		newFlags := draw.state
		changedFlags := currentState.state ^ draw.state // TODO golang 异或
		currentState.state = newFlags

		newStencil := draw.stencil
		changedStencil := currentState.stencil ^ draw.stencil
		currentState.stencil = newStencil

		// 2. Scissor?
		scissor := draw.scissor
		if currentState.scissor != scissor {
			currentState.scissor = scissor

			if scissor.isZero() {
				gl.Disable(gl.SCISSOR_TEST)
			} else {

				gl.Enable(gl.SCISSOR_TEST)
				gl.Scissor(int32(scissor.x),
					int32(ctx.wRect.h - scissor.h - scissor.y),
					int32(scissor.w),
					int32(scissor.h))
			}
		}

		// 3. stencil?
		if 0 != changedStencil {
			if 0 != newStencil {
				gl.Enable(gl.STENCIL_TEST)
				//// stencil not supported!!!
			}
		}

		// 4. state binding
		if 0 != (0 |
			ST.DEPTH_WRITE |
			ST.DEPTH_TEST_MASK |
			ST.RGB_WRITE |
			ST.ALPHA_WRITE |
			ST.BLEND_MASK |
			ST.BLEND_EQUATION_MASK |
			ST.PT_MASK ) & changedFlags {

			ctx.bindState(changedFlags, newFlags)

			pt := newFlags & ST.PT_MASK
			primIndex = uint8(pt >> ST.PT_SHIFT)
			prim = R_PrimInfo[primIndex]
		} /// End state change

		var programChanged bool
		var bindAttribs bool
		var constantsChanged bool

		/// 5. update program
		if key.Shader != shaderId {
			shaderId = key.Shader
			var id uint32 = ctx.R.shaders[shaderId].GLShader.Program
			gl.UseProgram(id)
			programChanged = true
			constantsChanged = true
			bindAttribs = true
		}

		/// 6. uniform binding
		if draw.uniformBegin < draw.uniformEnd {
			ctx.bindUniform(uint32(draw.uniformBegin), uint32(draw.uniformEnd))
		}

		/// 7. texture binding 如果纹理的采样类型变化，也要重新绑定！！
		for stage := 0; stage < 2; stage ++ {
			bind := draw.textures[stage]
			current := currentState.textures[stage]

			if  current != bind || programChanged {
				texture := ctx.R.textures[bind]
				texture.Bind(int32(stage))
			}

			current = bind
		}

		/// 8. vertex binding
		for stream := 0; stream < 2; stream ++ {
			vbStream := draw.vertexBuffers[stream]

			vSlot   := ctx.R.shaders[shaderId].stream[stream]
			vBuffer := ctx.R.vertexBuffers[vbStream.vertexBuffer]
			vLayout := ctx.R.vertexLayouts[vbStream.vertexLayout]

			var num uint8
			var _type AttribType
			var normalized, asInt bool

			gl.BindBuffer(gl.ARRAY_BUFFER, vBuffer.Id)
			for _, attr := range vLayout.attributes {
				Vertex_decode(attr, &num, &_type, &normalized, &asInt)

				gl.VertexAttribPointer( uint32(vSlot),			// Slot in Program
										int32(num), 			// component size
										R_AttribType[_type], 	// vertex type
										normalized, 			// normalized
										int32(vLayout.stride), 	// vertex stride
										unsafe.Pointer(0), 		// vertex pointer ! TODO
										)
			}
		}

		/// 9. draw
		if draw.indexBuffer != InvalidId {
			gl.DrawElements(prim, int32(draw.num), gl.UNSIGNED_SHORT, nil)
		} else {
			gl.DrawArrays(prim, int32(draw.firstIndex), int32(draw.num))
		}
	}
}

func (ctx *RenderContext) updateResolution() {

}

func (ctx *RenderContext) bindUniform(begin, end uint32) {
	ctx.ub.Seek(begin)

	for ub := ctx.ub; ub.GetPos() < end; {
		opcode := ub.ReadUInt32()

		if opcode == uint32(UniformEnd) {
			break
		}

		var uType, loc, size, num uint8
		Uniform_decode(opcode, &uType, &loc, &size, &num)

		data := ub.ReadPointer(uint32(size * num))
		switch UniformType(uType) {
		case UniformIntN:
			gl.Uniform1iv(int32(loc), int32(num), (*uint32)(data))
		case UniformFloatN:
			gl.Uniform1fv(int32(loc), int32(num), (*float32)(data))
		case UniformMat4:
			gl.Uniform1fv(int32(loc), 16, (*float32)(data))
		}
	}
}

func (ctx *RenderContext) bindAttributes() {

}

func (ctx *RenderContext) bindState(changedFlags, newFlags uint64) {
	if changedFlags & ST.DEPTH_WRITE != 0 {
		gl.DepthMask(ST.DEPTH_WRITE & newFlags != 0)
	}

	if changedFlags & ST.DEPTH_TEST_MASK != 0 {
		_func := (newFlags & ST.DEPTH_TEST_MASK) >> ST.DEPTH_TEST_SHIFT

		if _func != 0 {
			gl.Enable(gl.DEPTH_TEST)
			gl.DepthFunc(R_CmpFunc[_func])
		} else {
			if newFlags & ST.DEPTH_WRITE != 0 {
				gl.Enable(gl.DEPTH_TEST)
				gl.DepthFunc(gl.ALWAYS)
			} else {
				gl.Disable(gl.DEPTH_TEST)
			}
		}
	}

	if (ST.ALPHA_WRITE|ST.RGB_WRITE) & changedFlags != 0 {
		alpha := (newFlags & ST.ALPHA_WRITE) != 0
		rgb   := (newFlags & ST.RGB_WRITE) != 0
		gl.ColorMask(rgb, rgb, rgb, alpha)
	}

	/// 所谓 blend independent 可以实现顺序无关的 alpha 混合
	/// http://www.openglsuperbible.com/2013/08/20/is-order-independent-transparency-really-necessary/

	if ((ST.BLEND_MASK | ST.BLEND_EQUATION_MASK) & newFlags) != 0 || (blendFactor != draw.rgba) {
		enabled := (ST.BLEND_MASK & newFlags) != 0

		blend := uint32(newFlags & ST.BLEND_MASK) >> ST.BLEND_SHIFT
		srcRGB := (blend    ) & 0xFF
		dstRGB := (blend>> 4) & 0xFF
		srcA   := (blend>> 8) & 0xFF
		dstA   := (blend>>12) & 0xFF

		equ    := uint32(newFlags & ST.BLEND_EQUATION_MASK) >> ST.BLEND_EQUATION_SHIFT
		equRGB := (equ    ) & 0x7
		equA   := (equ>> 3) & 0x7

		if enabled{
			gl.Enable(gl.BLEND)
			gl.BlendFuncSeparate(R_BlendFactor[srcRGB], R_BlendFactor[dstRGB], R_BlendFactor[srcA], R_BlendFactor[dstA])
			gl.BlendEquationSeparate(R_BlendEquation[equRGB], R_BlendEquation[equA])
		} else {
			gl.Disable(gl.BLEND)
		}
	}
}
