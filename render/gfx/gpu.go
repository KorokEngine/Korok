package gfx

import (
	"korok/render/bx"
	"log"
	"unsafe"
)

// GPU Resources
type GpuRes struct {
	indexBuffers  [CONFIG_MAX_INDEX_BUFFERS]IndexBufferGL
	vertexBuffers [CONFIG_MAX_VERTEX_BUFFERS]VertexBufferGL
	shaders       [CONFIG_MAX_SHADERS]ShaderGL
	programs      [CONFIG_MAX_PROGRAMS]ProgramGL
	textures      [CONFIG_MAX_TEXTURES]TextureGL
	vertexLayouts [CONFIG_MAX_VERTEX_LAYOUT]VertexLayout
	frameBuffers  [CONFIG_MAX_FRAME_BUFFERS]FrameBufferGL
	uniforms      [CONFIG_MAX_UNIFORMS]interface{}
}


//////// bgfx 中的把 GPU 资源管理和渲染逻辑混合在一起，这让渲染部分的逻辑变得复杂
//////// 把资源管理部分分离出来。单独处理
//////// 提供类似于 Content.LoadTexture() handle 的 方法来管理 资源
//////// 在渲染层把 Content 共享给 RenderContext，这样渲染依然可以使用所以直接查找原来的数据
//////// 同时消灭了可恶的 CommandBuffer ！


type VertexBuffer struct {
	stride 	uint16
}

type DynamicIndexBuffer struct {
	handle IndexBufferHandle

	offset 	uint32
	size 	uint32
	startIndex uint32
	flags 	uint16
}

type DynamicVertexBuffer struct {
	handle VertexBufferHandle
	layout VertexLayoutHandle

	offset 	uint32
	size 	uint32
	startVertex 	uint32
	numVertices		uint32
	stride 			uint16
	flags 			uint16
}


type VertexLayoutRef struct {
	vertexLayoutMap bx.HandleHashMap

	vertexLayoutRef        [CONFIG_MAX_VERTEX_LAYOUT]uint16
	vertexBufferRef        [CONFIG_MAX_VERTEX_BUFFERS]VertexLayoutHandle
	dynamicVertexBufferRef [CONFIG_MAX_DYNAMIC_VERTEX_BUFFERS]VertexLayoutHandle
}

func (vd *VertexLayoutRef)init() {
	// zero memory
}

func (vd *VertexLayoutRef) shutdown(handleAlloc bx.HandleAlloc) {
	N := handleAlloc.GetNumHandles()
	for ii := uint16(0); ii < N; ii++ {
		handle := handleAlloc.GetHandleAt(ii)
		vd.vertexLayoutRef[handle] = 0
		vd.vertexLayoutMap.RemoveByHandle(handle)
		handleAlloc.Free(handle)
	}
}

func (vd *VertexLayoutRef) find(hash uint32) VertexLayoutHandle {
	handle := vd.vertexLayoutMap.Find(hash)
	return VertexLayoutHandle(handle)
}

func (vd *VertexLayoutRef) add(declHandle VertexLayoutHandle, hash uint32) {
	vd.vertexLayoutRef[declHandle.idx] ++
	vd.vertexLayoutMap.Insert(hash, uint16(declHandle))
}

func (vd *VertexLayoutRef) addVertexBuffer(handle VertexBufferHandle, declHandle VertexLayoutHandle, hash uint32) {
	vd.vertexLayoutRef[handle.idx] = uint16(declHandle)
	vd.vertexLayoutRef[declHandle.idx] ++
	vd.vertexLayoutMap.Insert(hash, uint16(declHandle))
}

func (vd *VertexLayoutRef) addDynamicBuffer(handle DynamicVertexBufferHandle, declHandle VertexLayoutHandle, hash uint32) {
	vd.dynamicVertexBufferRef[handle.idx] = declHandle
	vd.vertexLayoutRef[declHandle.idx]++
	vd.vertexLayoutMap.Insert(hash, uint16(declHandle))
}

func (vd *VertexLayoutRef) release(decHandle VertexLayoutHandle) VertexLayoutHandle {
	if isValid(idx(decHandle)) {
		vd.vertexBufferRef[decHandle.idx].idx --
		if 0 == vd.vertexLayoutRef[decHandle.idx] {
			vd.vertexLayoutMap.RemoveByHandle(uint16(decHandle))
			return decHandle
		}
	}

	invalid := INVALID_HANDLE
	return invalid
}

func (vd *VertexLayoutRef) releaseVertexBuffer(handle VertexBufferHandle) VertexLayoutHandle {
	declHandle := vd.vertexBufferRef[handle.idx]
	declHandle = vd.release(declHandle)
	vd.vertexBufferRef[handle.idx].idx = kInvalidHandle
	return declHandle
}

func (vd *VertexLayoutRef) releaseDynamicBuffer(handle DynamicVertexBufferHandle) VertexLayoutHandle {
	declHandle := vd.dynamicVertexBufferRef[handle.idx]
	declHandle = vd.release(declHandle)
	vd.dynamicVertexBufferRef[handle.idx].idx = kInvalidHandle
	return declHandle
}



type GPUContentManager struct {
	vertexBuffers [CONFIG_MAX_VERTEX_BUFFERS]VertexBuffer

	dynIndexBufferAllocator NonLocalAllocator
	dynVertexBufferAllocator NonLocalAllocator
	dynamicIndexBuffers [CONFIG_MAX_DYNAMIC_INDEX_BUFFERS]DynamicIndexBuffer
	dynamicVertexBuffers[CONFIG_MAX_DYNAMIC_VERTEX_BUFFERS]DynamicVertexBuffer

	numFreeDynamicIndexBufferHandles uint16
	numFreeDynamicVertexBufferHandles uint16

	freeDynamicIndexBufferHandle [CONFIG_MAX_DYNAMIC_INDEX_BUFFERS]DynamicIndexBufferHandle
	freeDynamicVertexBufferHandle[CONFIG_MAX_DYNAMIC_VERTEX_BUFFERS]DynamicVertexBufferHandle

	dynamicIndexBufferHandle  bx.HandleAlloc
	dynamicVertexBufferHandle bx.HandleAlloc


	// 缓存少量 Uniform 信息
	uniformSet 		map[uint16]bool
	uniformHashMap 	bx.HandleHashMap
	uniformRef 		[CONFIG_MAX_UNIFORMS]struct{
		name     string
		uType    UniformType
		num      uint16
		refCount int16
	}

	shaderHashMap	bx.HandleHashMap
	shaderRef 		[CONFIG_MAX_SHADERS]struct{
		uniforms []UniformHandle
		name 	string
		hash 	uint32
		num 	uint16
		refCount int16
	}

	programHashMap 	bx.HandleHashMap
	programRef 		[CONFIG_MAX_PROGRAMS]struct{
		vsh 	ShaderHandle
		fsh 	ShaderHandle
		refCount int16
	}

	textureRef 		[CONFIG_MAX_TEXTURES]struct{
		name 		string
		refCount 	int16
		bbRatio 	BackBufferRatio
		format 		uint8
		numMips 	uint8
		owned 		bool
	}
	frameBufferRef  [CONFIG_MAX_FRAME_BUFFERS]struct{
		window bool
		th	   [CONFIG_MAX_FB_ATTACHMENTS]TextureHandle
	}
	layoutRef VertexLayoutRef

	indexBufferHandle  bx.HandleAlloc
	vertexBufferHandle bx.HandleAlloc
	vertexLayoutHandle bx.HandleAlloc

	shaderHandle 	  bx.HandleAlloc
	programHandle     bx.HandleAlloc
	textureHandle     bx.HandleAlloc
	frameBufferHandle bx.HandleAlloc
	uniformHandle     bx.HandleAlloc
}


func (ctx *GPUContentManager) createShader(mem *Memory) ShaderHandle {
	return INVALID_HANDLE
}

func (ctx *GPUContentManager) getShaderUniforms(handle ShaderHandle) []UniformHandle{
	if !handle.isValid() {
		log.Println("Passing invalid shader handle to getShaderUniforms")
	}
	sr := ctx.shaderRef[handle.idx]
	return sr.uniforms
}

/// TODO setName 做什么用的呢？
func (ctx *GPUContentManager) setName(handle ShaderHandle, name string) {
	cmdBuf := ctx.getCommandBuffer(Command_SetName)
	cmdBuf.writeUInt16(uint16(handle))
	len := uint16(len(name))
	cmdBuf.writeUInt16(len)
	cmdBuf.writeString(name)
}

func (ctx *GPUContentManager) destroyShader(handle ShaderHandle) {
	if !handle.isValid() {
		log.Println("Passing invalid shader handle to destroyShader")
		return
	}
	ctx.shaderDecRef(handle)
}

func (ctx *Context) shaderTakeOwnership(handle ShaderHandle) {
	ctx.shaderDecRef(handle)
}

func (ctx *GPUContentManager) shaderIncRef(handle ShaderHandle) {
	ctx.shaderRef[handle.idx].refCount ++
}

func (ctx *GPUContentManager) shaderDecRef(handle ShaderHandle) {
	sr := &ctx.shaderRef[handle.idx]
	sr.refCount --

	if refs := sr.refCount; refs == 0 {
		ok := ctx.submit.FreeShader(handle)
		if !ok {
			log.Printf("Shader handle %d is already destroyed!", handle.idx)
		}

		cmdBuf := ctx.getCommandBuffer(Command_DestroyShader)
		cmdBuf.writeUInt16(uint16(handle))

		if sr.num != 0 {
			for i, num := uint16(0), sr.num; i < num; i++ {
				ctx.destroyUniform(sr.uniforms[i])
			}
			sr.uniforms = nil
			sr.num = 0
		}

		ctx.shaderHashMap.RemoveByHandle(uint16(handle.idx))
	}
}

func (ctx *GPUContentManager) createProgram(vsh ShaderHandle, fsh ShaderHandle, destroyShaders bool) ProgramHandle {
	if !vsh.isValid() || !fsh.isValid() {
		log.Printf("Vertex/fragment shader is invalid (vsh %d, fsh %d).", vsh.idx, fsh.idx)
		return INVALID_HANDLE
	}
	id := idx(ctx.programHashMap.Find(uint32(fsh.idx) << 16 | uint32(vsh.idx)))
	if id != kInvalidHandle {
		handle := ProgramHandle{id}
		pr := &ctx.programRef[handle.idx]
		pr.refCount ++
		ctx.shaderIncRef(pr.vsh)
		ctx.shaderIncRef(pr.fsh)
		return handle
	}

	vsr := ctx.shaderRef[vsh.idx]
	fsr := ctx.shaderRef[fsh.idx]

	if vsr.hash != fsr.hash {
		log.Println("Vertex shader output doesn't match fragment shader input")
		return INVALID_HANDLE
	}

	handle := ProgramHandle{idx(ctx.programHandle.Alloc())}

	if !handle.isValid() {
		log.Println("Failed to allocated program handle")
	}

	if handle.isValid() {
		ctx.shaderIncRef(vsh)
		ctx.shaderIncRef(fsh)
		pr := &ctx.programRef[handle.idx]
		pr.vsh = vsh
		pr.fsh = fsh
		pr.refCount = 1

		key := uint32(fsh.idx) << 16 | uint32(vsh.idx)
		ok := ctx.programHashMap.Insert(key, uint16(handle.idx))

		if !ok {
			log.Printf("Program already exists (key: %d, handle: %3d)!", key, handle.idx)
		}

		cmdBuf := ctx.getCommandBuffer(Command_CreateProgram)
		cmdBuf.writeUInt16(uint16(handle))
		cmdBuf.writeUInt16(uint16(vsh))
		cmdBuf.writeUInt16(uint16(fsh))
	}

	if destroyShaders {
		ctx.shaderTakeOwnership(vsh)
		ctx.shaderTakeOwnership(fsh)
	}

	return handle
}

func (ctx *GPUContentManager) destroyProgram(handle ProgramHandle) {
	pr := &ctx.programRef[handle.idx]
	ctx.shaderDecRef(pr.vsh)

	if pr.fsh.isValid() {
		ctx.shaderDecRef(pr.fsh)
	}

	pr.refCount --
	if refs := pr.refCount; refs == 0 {
		ok := ctx.submit.FreeProgram(handle)
		if !ok {
			log.Printf("Program handle %d is already destroyed!", handle.idx)
		}
		cmdBuf := ctx.getCommandBuffer(Command_DestroyProgram)
		cmdBuf.writeUInt16(uint16(handle))

		ctx.programHashMap.RemoveByHandle(uint16(handle.idx))
	}
}

func (ctx *GPUContentManager) createTexture(mem *Memory, flags uint32, skip uint8, info *TextureInfo, ratio BackBufferRatio) TextureHandle{
	if info == nil {
		info = &TextureInfo{}
	}

	/// Calculate Size todo
	{
		info.format = TEXTURE_FORMAT_COUNT
		info.storageSize = 0
		info.width = 0
		info.height = 0
		info.depth = 0
		info.numMips = 0
		info.bitsPerPixel = 0
		info.cubeMap = false
	}

	handle := TextureHandle{idx(ctx.textureHandle.Alloc())}

	if !handle.isValid() {
		log.Println("Failed to allocate texture handle")
	}

	if handle.isValid() {
		ref := ctx.textureRef[handle.idx]
		ref.refCount = 1
		ref.bbRatio = ratio
		ref.format = uint8(info.format)
		ref.numMips = info.numMips	// todo origin: numMips = imageContainer.numMips
		ref.owned = false

		cmdBuf := ctx.getCommandBuffer(Command_CreateTexture)
		cmdBuf.writeUInt16(uint16(handle))
		cmdBuf.writeMemory(mem)
		cmdBuf.writeUInt32(flags)
		cmdBuf.writeUInt8(skip)
	}

	return handle
}

func (ctx *GPUContentManager) setTextureName(handle TextureHandle, name string) {
	ref := ctx.textureRef[handle.idx]
	ref.name = name

	// setName(convert(handle), name)
}

func (ctx *Context) destroyTexture(handle TextureHandle) {
	if !handle.isValid() {
		log.Println("Passing invalid texture handle to destroyTexture")
		return
	}

	ctx.textureDecRef(handle)
}

func (ctx *Context) resizeTexture(handle TextureHandle, width, height uint16, numMips uint8) {

}

func (ctx *GPUContentManager) textureTakeOwnership(handle TextureHandle) {
	ref := ctx.textureRef[handle.idx]
	if !ref.owned {
		ref.owned = true
		ctx.textureDecRef(handle)
	}
}

func (ctx *GPUContentManager) textureIncRef(handle TextureHandle) {
	ctx.textureRef[handle.idx].refCount ++
}

func (ctx *GPUContentManager) textureDecRef(handle TextureHandle) {
	ref := &ctx.textureRef[handle.idx]
	ref.refCount --

	if refs := ref.refCount; 0 == refs {
		ref.name = ""

		ok := ctx.submit.FreeTexture(handle)
		if !ok {
			log.Printf("Texture handle %d is already destroyed!", handle.idx)
		}

		cmdBuf := ctx.getCommandBuffer(Command_DestroyTexture)
		cmdBuf.writeUInt16(uint16(handle))
	}
}

func (ctx *GPUContentManager) updateTexture(handle TextureHandle, side uint8, mip uint8, x, y, z, width, height, depth, pitch uint16, mem *Memory) {
	cmdBuf := ctx.getCommandBuffer(Command_UpdateTexture)
	cmdBuf.writeUInt16(uint16(handle))
	cmdBuf.writeUInt8(side)
	cmdBuf.writeUInt8(mip)

	rect := Rect {
		x: x,
		y: y,
		width: width,
		height: height,
	}

	cmdBuf.writeRect(&rect)
	cmdBuf.writeUInt16(z)
	cmdBuf.writeUInt16(depth)
	cmdBuf.writeUInt16(pitch)
	cmdBuf.writeMemory(mem)
}

func (ctx *GPUContentManager) checkFrameBuffer(num uint8, attachment []Attachment) bool {
	var color uint8
	var depth uint8

	for ii := uint8(0); ii < num; ii++ {
		texHandle := attachment[ii].handle
		if bimg.isDepth(ctx.textureRef[texHandle.idx].format) {
			depth ++
		} else {
			color ++
		}
	}

	return uint32(color) <= g_caps.limits.MaxFBAttachments && depth <= 1
}

func (ctx *GPUContentManager) createFrameBuffer(num uint8, attachment []Attachment, destroyTextures bool) FrameBufferHandle {
	if !ctx.checkFrameBuffer(num, attachment) {
		log.Printf("Too many frame buffer attachments (num attachments: %d, max color attachements %d)!",
			num,
			g_caps.limits.MaxFBAttachments)
	}
	handle := FrameBufferHandle{idx(ctx.frameBufferHandle.Alloc())}
	if !handle.isValid() {
		log.Println("Failed to allocate frame buffer handle")
	}

	if handle.isValid() {
		cmdBuf := ctx.getCommandBuffer(Command_CreateFrameBuffer)
		cmdBuf.writeUInt16(uint16(handle))
		cmdBuf.writeBool(false)
		cmdBuf.writeUInt8(num)

		ref := &ctx.frameBufferRef[handle.idx]
		ref.window = false
		bx.MemFill(unsafe.Pointer(ref.th), int(unsafe.Sizeof(TextureHandle{})))

		bbRatio := ctx.textureRef[attachment[0].handle.idx].bbRatio
		for i := uint8(0); i < num ; i++ {
			texHandle := attachment[i].handle
			if ctx.textureRef[texHandle.idx].bbRatio != bbRatio {
				log.Println("Mismatch in texture back-buffer ratio")
			}
			ref.th[i] = texHandle
			ctx.textureIncRef(texHandle)
		}
		cmdBuf.writeAttachmentArray(&attachment[0], int(num))
	}

	if destroyTextures {
		for i := 0; i < int(num); i++ {
			ctx.textureTakeOwnership(attachment[i].handle)
		}
	}
	return handle
}

func (ctx *GPUContentManager) getTexture(handle FrameBufferHandle, attachment uint8) TextureHandle {
	ref := ctx.frameBufferRef[handle.idx]
	if !ref.window {
		if uint32(attachment) < CONFIG_MAX_FB_ATTACHMENTS {
			return ref.th[attachment]
		} else {
			return ref.th[CONFIG_MAX_FB_ATTACHMENTS-1] // TODO ? 是否要减一
		}
	}
	return INVALID_HANDLE
}

func (ctx *GPUContentManager) destroyFrameBuffer(handle FrameBufferHandle) {
	ok := ctx.submit.FreeFrameBuffer(handle)
	if !ok {
		log.Printf("Frame buffer handle %d is already destroyed!", handle.idx)
	}
	cmdBuf := ctx.getCommandBuffer(Command_DestroyFrameBuffer)
	cmdBuf.writeUInt16(uint16(handle))

	if ref := &ctx.frameBufferRef[handle.idx]; !ref.window {
		for i := range ref.th {
			if th := ref.th[i]; th.isValid() {
				ctx.textureDecRef(th)
			}
		}
	}
}

func (ctx *GPUContentManager) createUniform(name string, _type UniformType, num uint16) UniformHandle {
	num = math.UInt16_max(1, num)

	id := idx(ctx.uniformHashMap.Find(bx.MurmurStr(name)))
	if id != kInvalidHandle {
		handle := UniformHandle{id}
		uniform := &ctx.uniformRef[handle.idx]

		if uniform.uType != _type {
			log.Printf("Uniform type mismatch (type: %d, expected %d).",
				_type,
				uniform.uType)
		}

		oldSize := g_uniformTypeSize[uniform.uType]
		newSize := g_uniformTypeSize[_type]

		if oldSize < newSize || uniform.num < num {
			if oldSize < newSize {
				uniform.uType = _type
			}
			uniform.num = math.UInt16_max(uniform.num, num)

			cmdBuf := ctx.getCommandBuffer(Command_CreateUniform)
			cmdBuf.writeUInt16(uint16(handle))
			cmdBuf.writeUInt8(uint8(uniform.uType))
			cmdBuf.writeUInt16(uniform.num)
			len := uint8(len(name) + 1)
			cmdBuf.writeUInt8(len)
			cmdBuf.writeString(name)
		}

		uniform.refCount ++
		return handle
	}

	handle := UniformHandle{idx(ctx.uniformHandle.Alloc())}

	if !handle.isValid() {
		log.Println("Failed to allocate uniform handle")
	}

	if handle.isValid() {
		log.Printf("Creating uniform (handle %3d) %s", handle.idx, name)

		uniform := ctx.uniformRef[handle.idx]
		uniform.name = name
		uniform.refCount = 1
		uniform.uType = _type
		uniform.num = num

		ok := ctx.uniformHashMap.Insert(bx.MurmurStr(name), uint16(handle.idx))
		if !ok {
			log.Printf("Uniform already exists (name: %s)!", name)
		}

		cmdBuf := ctx.getCommandBuffer(Command_CreateUniform)
		cmdBuf.writeUInt16(uint16(handle))
		cmdBuf.writeUInt8(uint8(_type))
		cmdBuf.writeUInt16(num)

		len := uint8(len(name) + 1)
		cmdBuf.writeUInt8(len)
		cmdBuf.writeString(name)
	}

	return handle
}

func (ctx *GPUContentManager) getUniformInfo(handle UniformHandle) (info *UniformInfo) {
	uniform := ctx.uniformRef[handle.idx]
	info.name = uniform.name
	info.uType = uniform.uType
	info.num = uniform.num
	return
}

func (ctx *GPUContentManager) destroyUniform(handle UniformHandle) {
	uniform := &ctx.uniformRef[handle.idx]
	if uniform.refCount <= 0 {
		log.Printf("Destroying already destroyed uniform %d", handle.idx)
	}
	uniform.refCount --
	if refs := uniform.refCount; refs == 0 {
		ok := ctx.submit.FreeUniform(handle)

		if !ok {
			log.Printf("Uniform handle %d is already destroyed!", handle.idx)
		}
		uniform.name = ""
		ctx.uniformHashMap.RemoveByHandle(uint16(handle.idx))

		cmdBuf := ctx.getCommandBuffer(Command_DestroyUniform)
		cmdBuf.writeUInt16(uint16(handle))
	}
}
