package asset

import (
	"korok.io/korok/audio/sine"
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
	if typ := audioType(file); typ == sine.None {
		log.Println("not implemented audio type:", file)
		return
	}
	var rid, cnt uint16
	if v, ok := am.repo[file]; ok {
		cnt = v.cnt
		rid = v.rid
	} else {
		id, _ := sine.R.LoadSound(file, audioType(file), sourceType(stream))
		rid = id
	}
	am.repo[file] = idCount{rid, cnt+1}
	log.Print("load file:", file)
}

// Unload delete raw Texture and any related SubTextures.
func (am *AudioManager) Unload(file string) {
	if v, ok := am.repo[file]; ok {
		if v.cnt > 1 {
			am.repo[file] = idCount{v.rid, v.cnt -1}
		} else {
			delete(am.repo, file)
			sine.R.UnloadSound(v.rid)
			log.Println("refCont == 0, delete resoruce!!")
		}
	}
}

func (am *AudioManager) Get(file string) (id uint16, ok bool){
	if v, ook := am.repo[file]; ook {
		id = v.rid
		ok = true
	}
	return
}

func audioType(name string) sine.FileType {
	switch true {
	case strings.HasSuffix(name, ".wav"):
		return sine.WAV
	case strings.HasSuffix(name, ".ogg"):
		return sine.VORB
	default:
		return sine.None
	}
}

func sourceType(stream bool) sine.SourceType {
	if stream {
		return sine.Stream
	} else {
		return sine.Static
	}
}


