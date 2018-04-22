package asset

type refCount struct {
	ref interface{} // 资源引用
	cnt int32       // 引用计数
}

type idCount struct {
	rid uint16 	// 资源ID
	cnt uint16  // 引用计数
}

/// 资源在各自相关的系统内管理，此处只提供工具方法
/// 和引用计数相关的操作



func Clear()  {
	// TODO
	//for _, v := range shaders {
	//	gl.DeleteProgram(v.Program)
	//}
	//
	//for _, v := range textures {
	//	gl.DeleteTextures(1, &v.Id)
	//}
}

var Shader *ShaderManager
var Texture *TextureManager
var Font *FontManager
var PSConfig *ParticleConfigManager
var Audio *AudioManager

func init() {
	Shader = &ShaderManager{}
	Audio = NewAudioManager()
	Texture = NewTextureManager()
	Font = NewFontManager()
	PSConfig = NewParticleConfigManager()
}