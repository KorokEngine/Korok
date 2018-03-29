package assets

import (
	"os"
	"fmt"

	"korok.io/korok/gfx/font"
)

type FontManager struct {
	repo map[string]refCount
}

func NewFontManager() *FontManager {
	return &FontManager{repo: make(map[string]refCount)}
}

func (fm *FontManager) LoadBitmap(name string, img, fc string) {
	var cnt int32 = 0
	var fnt interface{}

	if v, ok := fm.repo[name]; ok {
		cnt = v.cnt
		fnt = v.ref
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

	fm.repo[name] = refCount{fnt, cnt + 1}
	fmt.Println("load bitmap font sucess...", name)
}

func (fm *FontManager) LoadTrueType(name string, fc string) {
	var cnt int32 = 0
	var fnt interface{}

	if v, ok := fm.repo[name]; ok {
		cnt = v.cnt
		fnt = v.ref
	} else {
		fcr, err := os.Open(fc)
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := font.LoadTrueType(fcr, 24,  '0', 'z', 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		fnt = f
	}

	fm.repo[name] = refCount{fnt, cnt + 1}
	fmt.Println("load true-type font sucess...", name)
}

func (fm *FontManager) Unload(name string) {
	if v, ok := fm.repo[name]; ok {
		if v.cnt > 1 {
			fm.repo[name] = refCount{v.ref, v.cnt-1}
		} else {
			delete(fm.repo, name)
			// todo release font resource
			// v.ref.().Release()
		}
	}
}

func (fm *FontManager) GetFont(name string) (fnt font.Font) {
	if v, ok := fm.repo[name]; ok {
		fnt = v.ref.(font.Font)
	}
	return
}

