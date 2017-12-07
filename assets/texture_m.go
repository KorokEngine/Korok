package assets

import (
	"log"
	"os"
	"fmt"
	"image"
	"errors"

	"korok/gfx/bk"
	"korok/gfx"
)

type TextureManager struct {
	repo map[string]RefCount
}

func NewTextureManager() *TextureManager {
	return &TextureManager{make(map[string]RefCount)}
}

func (tm *TextureManager) Load(file string) {
	var rid, cnt uint16
	if v, ok := tm.repo[file]; ok {
		cnt = v.cnt
	} else {
		id, err := tm.loadTexture(file)
		if err != nil {
			log.Println(err)
		}
		rid = id
	}
	tm.repo[file] = RefCount{rid, cnt}
}

func (tm *TextureManager) GetTexture(file string) (uint16, *bk.Texture2D)  {
	if v, ok := tm.repo[file]; ok {
		if ok, tex := bk.R.Texture(v.rid); ok {
			return v.rid, tex
		}
	}
	return bk.InvalidId, nil
}

func (tm *TextureManager) Unload(file string) {
	if v, ok := tm.repo[file]; ok {
		if v.cnt > 1 {
			tm.repo[file] = RefCount{v.rid, v.cnt -1}
		} else {
			delete(tm.repo, file)
			bk.R.Free(v.rid)
		}
	}
}

func (tm *TextureManager) loadTexture(file string)(uint16, error)  {
	log.Println("load file:" + file)
	// 1. load file
	imgFile, err := os.Open(file)
	if err != nil {
		return bk.InvalidId, fmt.Errorf("texture %q not found: %v", file, err)
	}
	// 2. decode image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return bk.InvalidId, err
	}
	// 3. create
	if id, _ := bk.R.AllocTexture(img); id != bk.InvalidId {
		return id, nil
	}
	return bk.InvalidId, errors.New("fail to load texture")
}

// helper method
func AsSubTexture(id uint16) *gfx.SubTex{
	return &gfx.SubTex{
		TexId: id,
		Region: gfx.Region {
			0, 0,
			1, 1,
		},
	}
}