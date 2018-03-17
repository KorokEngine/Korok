package gfx

import (
	"korok.io/korok/gfx/bk"
	"unsafe"
)


// Anchor type
type Anchor uint8
const(
	ANCHOR_CENTER Anchor = 0x00
	ANCHOR_LEFT          = 0x01
	ANCHOR_RIGHT         = 0x02
	ANCHOR_UP  			 = 0x04
	ANCHOR_DOWN          = 0x08
)

type Size struct {
	Width, Height float32
}

type Region struct {
	X1, Y1 float32
	X2, Y2 float32
}


type bkTex struct {
	id, padding uint16
}

func (tex bkTex) Tex() uint16 {
	return uint16(tex.id)
}

func (tex bkTex) Region() Region {
	return Region{0, 0, 1, 1}
}

func (tex bkTex) Size() (sz Size) {
	if ok, t := bk.R.Texture(uint16(tex.id)); ok {
		sz = Size{t.Width, t.Height}
	}
	return
}

func NewTex(id uint16) Sprite {
	return bkTex{id:id}
}

// SubTexture = (bk-Texture << 16) + SubTextureId
type SubTex struct {
	id uint32
}

func (tex SubTex) Tex() uint16 {
	return R.texId(tex.id)
}

func (tex SubTex) Region() Region {
	return R.region(tex.id)
}

func (tex SubTex) Size() Size {
	return R.size(tex.id)
}

// A small struct to describe a group of sub-texture
type Atlas struct {
	id, aid uint16

	w, h float32

	// compiled sub-texture coordinate
	regions []Region

	// size of sub-texture
	sizes []Size

	// name of this atlas
	names []string

	// start index and size
	index, size uint16
}

func (at *Atlas) initialize(size int) {
	buffer := make([]byte, size * int(sizeOfRegion+sizeOfSize+sizeOfString))

	at.regions = (*[1<<16]Region)(unsafe.Pointer(&buffer[0]))[:]
	at.sizes = (*[1<<16]Size)(unsafe.Pointer(&buffer[size*int(sizeOfRegion)]))[:]
	at.names = (*[1<<16]string)(unsafe.Pointer(&buffer[size*int(sizeOfRegion+sizeOfSize)]))[:]
	at.index = 0
	at.size = uint16(size)
}

func (at *Atlas) release() {
	at.regions = nil
	at.sizes = nil
	at.names = nil
}

func (at *Atlas) AddItem(x, y, w, h float32, name string) {
	ii := at.index; at.index++

	at.sizes[ii] = Size{w, h}
	at.regions[ii] = Region{
		X1: x/at.w, Y1: y/at.h,
		X2: (x+w)/at.w, Y2:(y+h)/at.h,
	}
	at.names[ii] = name
}

func (at *Atlas) GetByName(name string) (tex SubTex) {
	for i := range at.names {
		if at.names[i] == name {
			tex = SubTex{uint32(i)} // todo atlas index
		}
	}
	return
}

func (at *Atlas) GetByIndex(index int) (tex SubTex) {
	return SubTex{uint32(index)}
}

func (at *Atlas) Region(ii int) Region {
	return at.regions[ii]
}


func (at *Atlas) Size(ii int) Size {
	return at.sizes[ii]
}

// Texture Resource Manager
type TexManager struct {
	atlases []Atlas
	frees []uint16

	// name to id
	names map[string]int

	// index and capacity
	index, cap uint16
}

// 纹理图集的管理是以纹理为单位.
func (tm *TexManager) NewAtlas(id uint16, size int, name string) (at *Atlas){
	if n := len(tm.frees); n > 0 {
		at = &tm.atlases[tm.frees[size-1]]
		tm.frees = tm.frees[:size-1]
	} else {
		ii := len(tm.atlases)
		tm.atlases = append(tm.atlases, Atlas{aid:uint16(ii)})
		at = &tm.atlases[ii]
		tm.names[name] = ii
	}

	at.initialize(size)
	at.id = id
	_, tex := bk.R.Texture(id)
	at.w, at.h = tex.Width, tex.Height

	return
}

func (tm *TexManager) Delete(name string) {
	if ii, ok := tm.names[name]; ok {
		tm.atlases[ii].release()
		tm.frees = append(tm.frees, uint16(ii))
	}
}

func (tm *TexManager) Atlas(name string) (at *Atlas) {
	if ii, ok := tm.names[name]; ok {
		at = &tm.atlases[ii]
	}
	return
}

// Region returns sub-texture's Region by id.
func (tm *TexManager) region(id uint32) (rg Region) {
	var (
		ai = id >> 16
		ii = id & 0xFFFF
	)
	at := tm.atlases[ai]
	rg = at.regions[ii]
	return
}

// Size returns sub-texture's Size by id.
func (tm *TexManager) size(id uint32) (sz Size) {
	var (
		ai = id >> 16
		ii = id & 0xFFFF
	)
	at := tm.atlases[ai]
	sz = at.sizes[ii]
	return
}

func (tm *TexManager) texId(id uint32) uint16 {
	return tm.atlases[id>>16].id
}

// TextureManager as a global variable.
var R *TexManager

// init
func init() {
	R = &TexManager{names:make(map[string]int, 0)}

	sizeOfRegion = unsafe.Sizeof(Region{})
	sizeOfSize = unsafe.Sizeof(Size{})
	sizeOfString = unsafe.Sizeof("")
}

var sizeOfRegion uintptr
var sizeOfSize uintptr
var sizeOfString uintptr
