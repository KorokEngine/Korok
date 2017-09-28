package gfx

import (
	"log"
	"korok/render/bx"
	"github.com/go-gl/gl/v3.2-core/gl"
)

/// OpenGL Model
/// Define: Texture/ Buffer/ Program/ Shader/
type IndexBufferGL struct {
	id 		uint32
	size 	uint32
	flags 	uint16
}

func (ib *IndexBufferGL) create(size uint32, data interface{}, flags uint16) {
	ib.size = size
	ib.flags = flags

	gl.GenBuffers(1, &ib.id)

	if 0 == ib.id {
		log.Println("Failed to generate buffer id.")
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.id)
	if data == nil {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(size), nil, gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(size), gl.Ptr(data), gl.STATIC_DRAW)
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

/// discard=false
func (ib *IndexBufferGL) update(offset uint32, size uint32, data interface{}, discard bool) {
	if 0 == ib.id {
		log.Println("Updating invalid index buffer.")
	}

	if discard {
		// orphan buffer
		ib.destroy()
		ib.create(ib.size, nil, ib.flags)
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.id)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(offset), int(size), gl.Ptr(data))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (ib *IndexBufferGL) destroy() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &ib.id)
}

type VertexBufferGL struct {
	id     uint32
	target uint32
	size   uint32
	layout VertexLayoutHandle
}

/// draw indirect >= es 3.0 or gl 4.0
func (vb *VertexBufferGL) create(size uint32, data interface{}, layoutHandle VertexLayoutHandle, flags uint16) {
	vb.size = size
	vb.layout = layoutHandle
	vb.target = gl.ARRAY_BUFFER

	gl.GenBuffers(1, &vb.id)
	if vb.id == 0 {
		log.Println("Failed to generate buffer id")
	}
	gl.BindBuffer(vb.target, vb.id)
	if data == nil {
		gl.BufferData(vb.target, int(size), gl.Ptr(data), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(vb.target, int(size), gl.Ptr(data), gl.STATIC_DRAW)
	}
	gl.BindBuffer(vb.target, 0)
}

/// discard = false
func (vb *VertexBufferGL) update(offset uint32, size uint32, data interface{}, discard bool) {
	if vb.id == 0 {
		log.Println("Updating invalid vertex buffer")
	}

	if discard {
		vb.destroy()
		vb.create(vb.size, nil, vb.layout, 0)
	}

	gl.BindBuffer(vb.target, vb.id)
	gl.BufferSubData(vb.target, int(offset), int(size), gl.Ptr(data))
	gl.BindBuffer(vb.target, 0)
}

func (vb *VertexBufferGL) destroy() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0 )
	gl.DeleteBuffers(1, &vb.id)
}

/// 纹理的处理还是很复杂的！！具体需要重新分析
type TextureGL struct {
	id 	uint32
	rbo uint32
	target 	uint32
	fmt 	uint32
	xType 	uint32

	flags 	uint32
	currentSamplerHash uint32
	width 	uint32
	height 	uint32
	depth 	uint32
	numMips	uint8
	requestedFormat uint8
	TextureFormat 	uint8
}

func (tex *TextureGL) init(target uint32, width, height, depth uint32, numMips uint8, flags uint32) {
	tex.target = target
	tex.numMips = numMips
	tex.flags = flags
	tex.width = width
	tex.height = height
	tex.depth = depth
	tex.currentSamplerHash = UINT32_MAX

	gl.GenTextures(1, &tex.id)

	if tex.id == 0 {
		log.Println("Failed to generate textue id.")
	}

	gl.BindTexture(target, tex.id)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
}

/// 解析纹理的 numMips/width/height/depth/format/
/// 从 flag 中取出 compute/srgb/msaa 等信息
/// 然后 glImage 创建 Texture
func (tex *TextureGL) create(mem *Memory, flags uint32, skip uint8) {

}

func (tex *TextureGL) destroy() {
	gl.BindTexture(tex.target, 0 )
	gl.DeleteTextures(1, &tex.id)

	if tex.rbo != 0 {
		gl.DeleteRenderbuffers(1, &tex.rbo)
		tex.rbo = 0
	}
}

func (tex *TextureGL) overrideInternal(ptr uintptr) {
	tex.destroy()
	tex.flags |= TEXTURE_INTERNAL_SHARED
	tex.id = uint32(ptr)
}

func (tex *TextureGL) update(side, mip uint8, rect Rect, z, depth, pitch uint16, mem *Memory) {

}

func (tex *TextureGL) setSampleState(flags uint32, rgba [4]float32) {
	if CONFIG_RENDERER_OPENGLES < 30 {

	}
}

func (tex *TextureGL) commit(stage uint32, flags uint32, palette [][4]float32) {

}

func (tex *TextureGL) resolve() {

}

func (tex *TextureGL) isCubeMap() bool{
	return gl.TEXTURE_CUBE_MAP == tex.target || gl.TEXTURE_CUBE_MAP_ARRAY == tex.target
}

type ShaderGL struct {
	id 		uint32
	xType 	uint32
	hash 	uint32
}

func (sh *ShaderGL) create(mem *Memory) {
	// TODO 这里的 Hash 计算
	sh.hash = bx.Murmur2A(uint32(mem.data), mem.size)

}

func (sh *ShaderGL) destroy() {
	if 0 != sh.id {
		gl.DeleteShader(sh.id)
		sh.id = 0
	}
}

type SwapChainGL struct {

}

type FrameBufferGL struct {
	swapChain *SwapChainGL
	fbo 	[2]uint32

	width 	uint32
	height 	uint32
	denseIdx uint16
	num 	uint8
	numTh 	uint8
	needPresent bool

	attachment [CONFIG_MAX_FB_ATTACHMENTS]Attachment
}

func (fb *FrameBufferGL) createWithAttachment(num uint8, attachment []Attachment) {
	gl.GenFramebuffers(1, &fb.fbo[0])

	fb.denseIdx = UINT16_MAX
	fb.numTh = num
	copy(fb.attachment[:], attachment)

	fb.needPresent = false
	fb.postReset()
}

func (fb *FrameBufferGL) create(denseIdx uint16, nwh interface{}, width, height uint32, depthFormat TextureFormat ) {
	fb.swapChain = nil // TODO fbo
	fb.width = width
	fb.height = height
	fb.numTh = 0
	fb.denseIdx = denseIdx
	fb.needPresent = false
}

func (fb *FrameBufferGL) postReset() {
	if 0 != fb.fbo[0] {
		gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo[0])

		var needResolve bool
		var buffers [CONFIG_MAX_FB_ATTACHMENTS]uint32
		var colorIdx uint32

		for i := 0; i < fb.numTh; i++ {
			if handle := fb.attachment[i].handle; handle.isValid() {
				texture := s_renderGL.textures[handle.idx]
				if 0 == colorIdx {
					fb.width = math.UInt32_max(texture.width >> fb.attachment[i].mip, 1)
					fb.height = math.UInt32_max(texture.height >> fb.attachment[i].mip, 1)
				}

				attachment := gl.COLOR_ATTACHMENT0 + colorIdx
				format := texture.TextureFormat

				if iimg.isDepth(format) {
					info := bimg.GetBlockInfo(format)
					if 0 < info.stencilBits {
						attachment = gl.DEPTH_STENCIL_ATTACHMENT
					} else if 0 == info.depthBits {
						attachment = gl.STENCIL_ATTACHMENT
					} else {
						attachment = gl.DEPTH_ATTACHMENT
					}
				} else {
					buffers[colorIdx] = attachment; colorIdx ++
				}

				if 0 != texture.rbo {
					gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, attachment, gl.RENDERBUFFER, texture.rbo)
				} else {
					target := texture.target
					gl.FramebufferTexture2D(gl.FRAMEBUFFER,
						attachment,
						target,
						texture.id,
						int32(fb.attachment[i].mip))
				}

				needResolve = needResolve || ((0 != texture.rbo) && (0 != texture.id))
			}
		}

		fb.num = uint8(colorIdx)

		/// validate frame buffer
		frameBufferValidate()

		if needResolve {
			gl.GenFramebuffers(1, &fb.fbo[1])
			gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo[1])

			colorIdx = 0
			for i := uint8(0); i < fb.numTh; i ++ {
				if handle := fb.attachment[i].handle; handle.isValid() {
					// texture
				}

			}

		}
	}
}

func (fb *FrameBufferGL) destroy() uint16{
	if fb.num != 0 {
		if fb.fbo[1] == 0 {
			gl.DeleteFramebuffers(1, &fb.fbo[0])
		} else {
			gl.DeleteFramebuffers(2, &fb.fbo[1])
		}
		fb.num = 0
	}

	if fb.swapChain != nil {
		// s_renderGL.glctx.
	}
	fb.fbo[0], fb.fbo[1] = 0, 0
	denseIdx := fb.denseIdx
	fb.needPresent = false
	fb.numTh = 0
	return denseIdx
}

func (fb *FrameBufferGL) resolve() {
	if fb.fbo[1] != 0 {
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fb.fbo[0])
		gl.ReadBuffer(gl.COLOR_ATTACHMENT0)
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, fb.fbo[1])
		gl.BlitFramebuffer(0,
			0,
			int32(fb.width),
			int32(fb.height),
			0,
			0,
			int32(fb.width),
			int32(fb.height),
			gl.COLOR_BUFFER_BIT,
			gl.LINEAR)
		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fb.fbo[0])
		gl.ReadBuffer(gl.NONE)
		gl.BindFramebuffer(gl.FRAMEBUFFER)
	}
}

func (fb *FrameBufferGL) discard(flags uint16) {

}

type ProgramGL struct {
	id uint32

	unboundUsedAttrib [ATTRIB_COUNT]uint8	// For tracking unbound used attributes between begin()/end()
	usedCount 	uint8
	used		[ATTRIB_COUNT]uint8
	attributes 	[ATTRIB_COUNT]uint32

	sampler 	[CONFIG_MAX_SAMPLERS]uint32
	numSamplers uint8

	constantBuffer *UniformBuffer
	predefinedUniform [UNIFORM_COUNT]PreDefUniformEnum
	numPredefined 	uint8
}

func (p *ProgramGL) create(vsh, fsh *ShaderGL) {

}

func (p *ProgramGL) destroy() {

}

func (p *ProgramGL) init() {

}

// baseVertex = 0
func (p *ProgramGL) bindInstanceData(stride uint32, baseVertex uint32) {

}

func (p *ProgramGL) bindAttributeBeign() {

}

func (p *ProgramGL) bindAttributeEnd() {

}

