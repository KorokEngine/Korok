package asset

import (
	"korok.io/korok/audio/ap"
	"strings"
	"log"
)

type AudioManager struct {
	repo map[string]idCount
}

func NewAudioManager() *AudioManager {
	return &AudioManager{make(map[string]idCount)}
}

// Load loads a single Texture file.
func (am *AudioManager) Load(file string, stream bool) {
	if typ := audioType(file); typ == ap.None {
		log.Println("not implemented audio type:", file)
		return
	}
	var rid, cnt uint16
	if v, ok := am.repo[file]; ok {
		cnt = v.cnt
		rid = v.rid
	} else {
		id, _ := ap.R.LoadSound(file, audioType(file), sourceType(stream))
		rid = id
	}
	am.repo[file] = idCount{rid, cnt+1}
}

// Unload delete raw Texture and any related SubTextures.
func (am *AudioManager) Unload(file string) {
	if v, ok := am.repo[file]; ok {
		if v.cnt > 1 {
			am.repo[file] = idCount{v.rid, v.cnt -1}
		} else {
			delete(am.repo, file)
			ap.R.UnloadSound(v.rid)
			log.Println("refCont == 0, delete resoruce!!")
		}
	}
}

func (am *AudioManager) Get(file string) (id uint16){
	if v, ok := am.repo[file]; ok {
		id = v.rid
	}
	return
}

func audioType(name string) ap.FileType {
	switch true {
	case strings.HasSuffix(name, ".wav"):
		return ap.WAV
	case strings.HasSuffix(name, ".ogg"):
		return ap.VORB
	default:
		return ap.None
	}
}

func sourceType(stream bool) ap.SourceType {
	if stream {
		return ap.Stream
	} else {
		return ap.Static
	}
}


