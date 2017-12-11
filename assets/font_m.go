package assets

import (
	"os"
	"fmt"
	"korok/gfx/font"
)

type FRefCount struct {
	cnt int32
	fnt *font.Font
}

type FontManager struct {
	repo map[string]FRefCount
}

func NewFontManager() *FontManager {
	return &FontManager{repo: make(map[string]FRefCount)}
}

func (fm *FontManager) LoadBitmap(name string, img, fc string) {
	var cnt int32 = 0
	var fnt *font.Font

	if v, ok := fm.repo[name]; ok {
		cnt = v.cnt
		fnt = v.fnt
	} else {
		ir, err := os.Open(img)
		if err != nil {
			fmt.Println(err)
			return
		}
		fcr, err := os.Open(fc)
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := font.LoadBitmap(ir, fcr, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		fnt = f
	}

	fm.repo[name] = FRefCount{cnt + 1, fnt}
	fmt.Println("load bitmap font sucess...", name)
}

func (fm *FontManager) LoadTrueType(name string, fc string) {
	var cnt int32 = 0
	var fnt *font.Font

	if v, ok := fm.repo[name]; ok {
		cnt = v.cnt
		fnt = v.fnt
	} else {
		fcr, err := os.Open(fc)
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := font.LoadTrueType(fcr, 1,  '0', 'z', 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		fnt = f
	}

	fm.repo[name] = FRefCount{cnt + 1, fnt}
	fmt.Println("load true-type font sucess...", name)
}

func (fm *FontManager) Unload(name string) {
	if v, ok := fm.repo[name]; ok {
		if v.cnt > 1 {
			fm.repo[name] = FRefCount{v.cnt-1, v.fnt}
		} else {
			delete(fm.repo, name)
			v.fnt.Release()
		}
	}
}

func (fm *FontManager) GetFont(name string) (fnt *font.Font) {
	if v, ok := fm.repo[name]; ok {
		fnt = v.fnt
	}
	return
}

