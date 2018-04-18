package bk

import (
	"korok.io/korok/hid/gl"
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
	// pixel-ratio = windows-size/frame-buffer-size
	pixelRatio float32

	// clips rect, index-0 is a default zero-rect.
	clips []Rect

	backBufferFbo uint32
}

func NewRenderContext(r *ResManager, ub *UniformBuffer) *RenderContext {
	return &RenderContext{
		R:  r,
		ub: ub,
		clips:make([]Rect, 1),
	}
}

func (ctx *RenderContext) Init() {
	ctx.vaoSupport = gl.NeedVao()
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

// Reset OpenGL state, then each frame has same starting state.
func (ctx *RenderContext) Reset() {
	ctx.clips = ctx.clips[:1]
	gl.Disable(gl.SCISSOR_TEST)
}

func (ctx *RenderContext) AddClipRect(x, y, w, h uint16) uint16 {
	index := uint16(len(ctx.clips))
	ratio := ctx.pixelRatio
	ctx.clips = append(ctx.clips, Rect{
		uint16(float32(x)*ratio),
		uint16(float32(y)*ratio),
		uint16(float32(w)*ratio),
		uint16(float32(h)*ratio)})
	return index
}

func (ctx *RenderContext) Draw(sortKeys []uint64, sortValues []uint16, drawList []RenderDraw) {
	// if vao support
	if defaultVao := ctx.vao; 0 != defaultVao {
		gl.BindVertexArray(defaultVao)
	}

	// init per-draw state
	var (
		currentState = RenderDraw{}
		shaderId     = InvalidId
		key          = SortKey{}
		primIndex    = uint8(uint64(0) >> ST.PT_SHIFT)
		prim         = g_PrimInfo[primIndex]
		stateBits    = ST.DEPTH_WRITE| ST.DEPTH_TEST_MASK| ST.RGB_WRITE| ST.ALPHA_WRITE| ST.BLEND_MASK| ST.PT_MASK
	)

	// Let's Render!!
	for item := range sortKeys {
		encodedKey := sortKeys[item]
		itemId := sortValues[item]
		key.Decode(encodedKey)

		draw := drawList[itemId]

		// 1. 求取变化的状态位
		newFlags := draw.state
		changedFlags := currentState.state ^ draw.state
		currentState.state = newFlags

		newStencil := draw.stencil
		changedStencil := currentState.stencil ^ draw.stencil
		currentState.stencil = newStencil

		// 2. Scissor
		if scissor := draw.scissor; currentState.scissor != scissor {
			currentState.scissor = scissor
			clip := ctx.clips[scissor]

			if clip.isZero() {
				gl.Disable(gl.SCISSOR_TEST)
			} else {
				gl.Enable(gl.SCISSOR_TEST)
				gl.Scissor(int32(clip.x), int32(clip.y), int32(clip.w), int32(clip.h))
			}
		}

		// 3. stencil
		if changedStencil != 0 {
			if newStencil != 0 {
				if (gDebug & DebugQueue) != 0 {
					log.Println("Renderc disable stencil")
				}
				gl.Enable(gl.STENCIL_TEST)
			} else {
				gl.Disable(gl.STENCIL_TEST)
				if (gDebug & DebugQueue) != 0 {
					log.Println("Renderc enable stencil")
				}
			}
		}

		// 4. state binding
		if (stateBits&changedFlags) != 0 {
			ctx.bindState(changedFlags, newFlags)
			pt := newFlags & ST.PT_MASK
			primIndex = uint8(pt >> ST.PT_SHIFT)
			prim = g_PrimInfo[primIndex]
		}

		var programChanged bool
		//var bindAttribs bool
		//var constantsChanged bool

		// 5. Update program
		if key.Shader != shaderId {
			shaderId = key.Shader
			var id = ctx.R.shaders[shaderId].GLShader.Program
			gl.UseProgram(id)
			programChanged = true
			//constantsChanged = true
			//bindAttribs = true
			//log.Println("bind program")
		}

		// 6. uniform binding
		if draw.uniformBegin < draw.uniformEnd {
			ctx.bindUniform(uint32(draw.uniformBegin), uint32(draw.uniformEnd))
		}

		// 7. texture binding 如果纹理的采样类型变化，也要重新绑定！！
		for stage := 0; stage < 2; stage++ {
			bind := draw.textures[stage]
			current := currentState.textures[stage]
			if InvalidId != bind {
				if current != bind || programChanged {
					texture := ctx.R.textures[bind]
					texture.Bind(int32(stage))
				}
			}
			currentState.textures[stage] = bind
		}

		// 8. index & vertex binding TODO 优化 attribute 绑定
		shader := ctx.R.shaders[shaderId]
		shader.BindAttributes(ctx.R, draw.vertexBuffers[:])

		if ib := draw.indexBuffer; ib != InvalidId && ib != currentState.indexBuffer {
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ctx.R.indexBuffers[ib].Id)
			currentState.indexBuffer = ib
		}

		// 9. draw
		if draw.indexBuffer != InvalidId {
			offset := int(draw.firstIndex) * 2 // 2 = sizeOf(unsigned_short)
			gl.DrawElements(prim, int32(draw.num), gl.UNSIGNED_SHORT, offset)
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
		case UniformInt1:
			gl.Uniform1iv(int32(loc), int32(num), (*int32)(data))
		case UniformVec1:
			gl.Uniform1fv(int32(loc), int32(num), (*float32)(data))
		case UniformVec4:
			gl.Uniform4fv(int32(loc), int32(num), (*float32)(data))
		case UniformMat3:
			gl.UniformMatrix3fv(int32(loc), int32(num), false, (*float32)(data))
		case UniformMat4:
			gl.UniformMatrix4fv(int32(loc), int32(num), false, (*float32)(data))
		case UniformSampler:
			gl.Uniform1i(int32(loc), *(*int32)(data))
		}
	}
}

func (ctx *RenderContext) bindAttributes() {

}

func (ctx *RenderContext) bindState(changedFlags, newFlags uint64) {
	if changedFlags&ST.DEPTH_WRITE != 0 {
		gl.DepthMask(newFlags&ST.DEPTH_WRITE != 0)
		log.Printf("depth mask state: %v", ST.DEPTH_WRITE&newFlags != 0)
	}

	if changedFlags&(ST.ALPHA_WRITE|ST.RGB_WRITE) != 0 {
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
	if changedFlags&ST.BLEND_MASK != 0 {
		blend := uint16(newFlags&ST.BLEND_MASK) >> ST.BLEND_SHIFT
		if blend != 0 {
			gl.Enable(gl.BLEND)
			gl.BlendFunc(g_Blend[blend].Src, g_Blend[blend].Dst)
		} else {
			gl.Disable(gl.BLEND)
		}
	}
}
