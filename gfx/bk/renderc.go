package bk

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"log"
)

type RenderContext struct {
	// data
	R  *ResManager
	ub *UniformBuffer

	// draw state
	vao        uint32
	vaoSupport bool

	// window rect
	wRect Rect

	backBufferFbo uint32
}

func NewRenderContext(r *ResManager, ub *UniformBuffer) *RenderContext {
	return &RenderContext{
		R:  r,
		ub: ub,
	}
}

func (ctx *RenderContext) Init() {
	ctx.vaoSupport = true
	if ctx.vaoSupport {
		gl.GenVertexArrays(1, &ctx.vao)
	}
}

func (ctx *RenderContext) Shutdown() {
	if ctx.vao != 0 {
		gl.BindVertexArray(0)
		gl.DeleteVertexArrays(1, &ctx.vao)
	}
}

func (ctx *RenderContext) Draw(sortKeys []uint64, sortValues []uint16, drawList []RenderDraw) {
	// 1. 绑定 VAO 和 FrameBuffer
	if defaultVao := ctx.vao; 0 != defaultVao {
		gl.BindVertexArray(defaultVao)
	}
	// gl.BindFramebuffer(gl.FRAMEBUFFER, ctx.backBufferFbo)

	// 2. 更新分辨率
	// ctx.updateResolution(&render.resolution)

	// 3. 初始化 per-draw state
	currentState := RenderDraw{}
	//
	shaderId := InvalidId
	key := SortKey{}
	//
	primIndex := uint8(uint64(0) >> ST.PT_SHIFT)
	prim := g_PrimInfo[primIndex]

	// 8. Render!!

	sortKeys = sortKeys[1:]
	sortValues = sortValues[1:]

	for item := range sortKeys {
		encodedKey := sortKeys[item]
		itemId := sortValues[item]
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
				if (g_debug & DEBUG_Q) != 0 {
					log.Println("Renderc disable scissor")
				}
				gl.Disable(gl.SCISSOR_TEST)
			} else {
				if (g_debug & DEBUG_Q) != 0 {
					log.Printf("Renderc enable scissor: %q", scissor)
				}
				gl.Enable(gl.SCISSOR_TEST)
				gl.Scissor(int32(scissor.x),
					int32(ctx.wRect.h-scissor.h-scissor.y),
					int32(scissor.w),
					int32(scissor.h))
			}
		}

		// 3. stencil?
		if 0 != changedStencil {
			if 0 != newStencil {
				if (g_debug & DEBUG_Q) != 0 {
					log.Println("Renderc disable stencil")
				}
				gl.Enable(gl.STENCIL_TEST)
				//// stencil not supported!!!
			} else {
				gl.Disable(gl.STENCIL_TEST)
				if (g_debug & DEBUG_Q) != 0 {
					log.Println("Renderc enable stencil")
				}
			}
		}

		// 4. state binding
		if 0 != (0|
			ST.DEPTH_WRITE|
			ST.DEPTH_TEST_MASK|
			ST.RGB_WRITE|
			ST.ALPHA_WRITE|
			ST.BLEND_MASK|
			ST.PT_MASK)&changedFlags {

			ctx.bindState(changedFlags, newFlags)

			pt := newFlags & ST.PT_MASK
			primIndex = uint8(pt >> ST.PT_SHIFT)
			prim = g_PrimInfo[primIndex]
		} /// End state change

		var programChanged bool
		//var bindAttribs bool
		//var constantsChanged bool

		/// 5. Update program
		if key.Shader != shaderId {
			shaderId = key.Shader
			var id uint32 = ctx.R.shaders[shaderId].GLShader.Program
			gl.UseProgram(id)
			programChanged = true
			//constantsChanged = true
			//bindAttribs = true
		}

		/// 6. uniform binding
		if draw.uniformBegin < draw.uniformEnd {
			ctx.bindUniform(uint32(draw.uniformBegin), uint32(draw.uniformEnd))
		}

		/// 7. texture binding 如果纹理的采样类型变化，也要重新绑定！！
		for stage := 0; stage < 2; stage++ {
			bind := draw.textures[stage]
			current := currentState.textures[stage]

			if current != bind || programChanged {
				texture := ctx.R.textures[bind]
				texture.Bind(int32(stage))
			}

			current = bind
		}

		/// 8. vertex binding
		shader := ctx.R.shaders[shaderId]
		shader.BindAttributes(ctx.R, draw.vertexBuffers[:])

		/// 9. draw
		if draw.indexBuffer != InvalidId {
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ctx.R.indexBuffers[draw.indexBuffer&0x0FFF].Id)
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
		data := ub.ReadPointer(uint32(size) * uint32(num))

		switch UniformType(uType) {
		case UniformIntN:
			gl.Uniform1iv(int32(loc), int32(num), (*int32)(data))
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
	if changedFlags&ST.DEPTH_WRITE != 0 {
		gl.DepthMask(ST.DEPTH_WRITE&newFlags != 0)
		log.Printf("depth mask state: %v", ST.DEPTH_WRITE&newFlags != 0)
	}

	if (ST.ALPHA_WRITE|ST.RGB_WRITE)&changedFlags != 0 {
		alpha := (newFlags & ST.ALPHA_WRITE) != 0
		rgb := (newFlags & ST.RGB_WRITE) != 0
		gl.ColorMask(rgb, rgb, rgb, alpha)

		log.Printf("color mask state: (%d, %d)", rgb, alpha)
	}

	if changedFlags&ST.DEPTH_TEST_MASK != 0 {
		_func := (newFlags & ST.DEPTH_TEST_MASK) >> ST.DEPTH_TEST_SHIFT

		if _func != 0 {
			gl.Enable(gl.DEPTH_TEST)
			gl.DepthFunc(g_CmpFunc[_func])

			log.Printf("set depth-test func: %d", _func)
		} else {
			if newFlags&ST.DEPTH_WRITE != 0 {
				gl.Enable(gl.DEPTH_TEST)
				gl.DepthFunc(gl.ALWAYS)

				log.Println("set depth-test always")
			} else {
				gl.Disable(gl.DEPTH_TEST)

				log.Println("disable depth-test")
			}
		}
	}

	/// 所谓 blend independent 可以实现顺序无关的 alpha 混合
	/// http://www.openglsuperbible.com/2013/08/20/is-order-independent-transparency-really-necessary/
	if ((ST.BLEND_MASK) & newFlags) != 0 {
		enabled := (ST.BLEND_MASK & newFlags) != 0

		blend := uint16(newFlags&ST.BLEND_MASK) >> ST.BLEND_SHIFT
		if enabled {
			gl.Enable(gl.BLEND)
			gl.BlendFunc(g_Blend[blend].Src, g_Blend[blend].Dst)

			log.Printf("set blend-func: %d", blend)
		} else {
			gl.Disable(gl.BLEND)
			log.Println("disable blend-func")
		}
	}
}
