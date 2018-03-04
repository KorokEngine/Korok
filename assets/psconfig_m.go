package assets

// 粒子系统配置文件管理
type ParticleConfigManager struct {

}

func NewParticleConfigManager() *ParticleConfigManager {
	return &ParticleConfigManager{}
}

func (pcm *ParticleConfigManager) Load(file string) {

}

func (pcm *ParticleConfigManager) Unload(file string) {

}

func (pcm *ParticleConfigManager) Get(file string) interface{} {
	return nil
}

