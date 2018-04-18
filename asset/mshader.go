package asset

type ShaderManager struct {
}

func (sm *ShaderManager) GetShaderStr(key string) (string, string) {
	switch key {
	case "dft", "mesh":
		return vertex, color
	case "batch":
		return bVertex, bColor
	}
	return "", ""
}

