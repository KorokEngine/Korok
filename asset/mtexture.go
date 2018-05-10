package asset

import (
	"golang.org/x/mobile/asset"
	"korok.io/korok/gfx/bk"
	"korok.io/korok/gfx"

	"log"
	"fmt"
	"image"
	"errors"
	"io/ioutil"
	"encoding/json"
)

type TextureManager struct {
	repo map[string]idCount
	names map[string]uint32
}

func NewTextureManager() *TextureManager {
	return &TextureManager{
		make(map[string]idCount),
		make(map[string]uint32),
		}
}

// Load loads a single Texture file.
func (tm *TextureManager) Load(file string) {
	var rid, cnt uint16
	if v, ok := tm.repo[file]; ok {
		cnt = v.cnt
		rid = v.rid
	} else {
		// create bk.Texture2D
		id, err := tm.loadTexture(file)
		if err != nil {
			log.Println(err)
		}
		rid = id
	}
	tm.repo[file] = idCount{rid, cnt+1}
}

// Unload delete raw Texture and any related SubTextures.
func (tm *TextureManager) Unload(file string) {
	if v, ok := tm.repo[file]; ok {
		if v.cnt > 1 {
			tm.repo[file] = idCount{v.rid, v.cnt -1}
		} else {
			delete(tm.repo, file)
			bk.R.Free(v.rid)
			// maybe it's a atlas, try to delete
			gfx.R.Delete(file)

			log.Println("refCont == 0, delete resoruce!!")
		}
	}
}

// LoadAtlas loads the atlas with a description file.
// The SubTexture can be found by SubTexture's name.
func (tm *TextureManager) LoadAtlas(file, desc string) {
	var rid, cnt uint16
	if v, ok := tm.repo[file]; ok {
		cnt = v.cnt
		rid = v.rid
	} else {
		id, data, err := tm.loadAtlas(file, desc)
		if err != nil {
			log.Println(err)
			return
		}
		size := len(data.Frames)

		// new atlas
		at := gfx.R.NewAtlas(id, size, file)

		// fill
		for _, f := range data.Frames {
			at.AddItem(float32(f.Frame.X), float32(f.Frame.Y), float32(f.Frame.W), float32(f.Frame.H), f.Filename, f.Rotated)
		}
		rid = id
	}
	tm.repo[file] = idCount{rid, cnt+1}
}


// LoadAtlasIndexed loads the atlas with specified with/height/num.
func (tm *TextureManager) LoadAtlasIndexed(file string, width, height float32, row, col int) {
	var rid, cnt uint16
	if v, ok := tm.repo[file]; ok {
		cnt = v.cnt
		rid = v.rid
	} else {
		id, err := tm.loadTexture(file)
		if err != nil {
			log.Println(err)
		}
		size := row * col

		// new atlas
		at := gfx.R.NewAtlas(id, size, file)

		// fill
		for i := 0; i < row; i ++ {
			for j := 0; j < col; j ++ {
				at.AddItem(float32(j)*width, float32(i)*height, width, height, "", false)
			}
		}
	}
	tm.repo[file] = idCount{rid, cnt+1}
}

// Get returns the low-level Texture.
func (tm *TextureManager) Get(file string) gfx.Tex2D {
	rid := tm.repo[file]
	return gfx.NewTex(rid.rid)
}

// Get returns the low-level Texture.
func (tm *TextureManager) GetRaw(file string) (uint16, *bk.Texture2D)  {
	if v, ok := tm.repo[file]; ok {
		if ok, tex := bk.R.Texture(v.rid); ok {
			return v.rid, tex
		}
	}
	return bk.InvalidId, nil
}

// Atlas returns the Atlas.
func (tm *TextureManager) Atlas(file string) (at *gfx.Atlas, ok bool) {
	if _, ok := tm.repo[file]; ok {
		at = gfx.R.Atlas(file)
		ok = at != nil
	}
	return
}


func (tm *TextureManager) loadTexture(file string)(uint16, error)  {
	log.Println("load file:" + file)
	// 1. load file
	imgFile, err := asset.Open(file)
	if err != nil {
		return bk.InvalidId, fmt.Errorf("texture %q not found: %v", file, err)
	}
	// 2. decode image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return bk.InvalidId, err
	}
	// 3. create raw texture
	if id, _ := bk.R.AllocTexture(img); id != bk.InvalidId {
		return id, nil
	}
	return bk.InvalidId, errors.New("fail to load texture")
}

// 加载纹理图集
func (tm *TextureManager) loadAtlas(img, desc string)(id uint16, at *atlas, e error) {
	id, err := tm.loadTexture(img)
	if err != nil {
		e = err
		return
	}
	file, err := asset.Open(desc)
	defer file.Close()

	if err != nil {
		e = err
		return
	}
	d, err := ioutil.ReadAll(file)
	if err != nil {
		e = err
		return
	}
	at = &atlas{}
	e = json.Unmarshal(d, at)
	return
}

// Field int `json:"myName"`
// The file format is TexturePacker's generic json-array format.
// TexturePacker: https://www.codeandweb.com/texturepacker
type atlas struct {
	Meta struct {
		App string `json:"app"`
		Version string `json:"version"`
		Image string   `json:"image"`
		Format string  `json:"format"`
		Size struct{
			W, H int
		} `json:"size"`
		Scale float32 `json:"scale,string"`
	} `json:"meta"`

	Frames []struct{
		Filename string `json:"filename"`
		Frame struct{ X, Y, W, H int} `json:"frame"`
		Rotated bool `json:"rotated"`
		Trimmed bool `json:trimmed`
		Pivot struct{X, Y float32} `json:"pivot"`
	} `json:"frames"`
}