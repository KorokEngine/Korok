package bk

// 应该尽可能的编码一些重要的状态切换信息

const (
	SK_LayerMask   uint32 = 0xF0000000
	SK_ShaderMask  uint32 = 0x0F000000
	SK_BlendMask   uint32 = 0x00F00000
	SK_TextureMask uint32 = 0x000FF000
)

// SortKey FORMAT
// 64bit:
// 0000 - 0000000000 -       00000 -  000 - 0000000000
//  ^        ^                  ^      ^       ^
//  |        |                  |      |       |
//  |      z-order            shader   |    texture
// Layer                             blend
//
type SortKey struct {
	Layer   uint16
	Order   uint16
	Shader  uint16
	Blend   uint16
	Texture uint16
}

func (sk *SortKey) Encode() (key uint64) {
	return 0 |
		uint64(sk.Layer  )<<28 |
		uint64(sk.Shader )<<23 |
		uint64(sk.Blend  )<<20 |
		uint64(sk.Texture)<<10
}

func (sk *SortKey) Decode(key uint64) {
	sk.Texture = uint16((key >>  0) & (2<<10 - 1))
	sk.Texture = uint16((key >> 10) & (2<<10 - 1))
	sk.Blend   = uint16((key >> 20) & (2<< 3 - 1))
	sk.Shader  = uint16((key >> 23) & (2<< 5 - 1))
	sk.Layer   = uint16((key >> 28) & (2<< 4 - 1))
}

func SkDecode(key uint64) (sk SortKey) {
	return
}

//
//func (sk *SortKey) SetTexture(t uint32) {
//	v := uint32(*sk)
//	v = (v & ^TextureMask) | (t << 12 & TextureMask)
//	*sk = SortKey(v)
//}
