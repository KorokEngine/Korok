package asset

import (
	"korok.io/korok/gfx/font"

	"fmt"
	"log"

	"korok.io/korok/asset/res"
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
		ir, err := res.Open(img)
		if err != nil {
			fmt.Println(err)
			return
		}
		fcr, err := res.Open(fc)
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

func (fm *FontManager) LoadTrueType(name string, file string, lc font.TTFConfig) {
	var cnt int32 = 0
	var fnt interface{}

	if v, ok := fm.repo[name]; ok {
		cnt = v.cnt
		fnt = v.ref
	} else {
		fcr, err := res.Open(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := font.LoadTrueType(fcr, lc)
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
			fm.repo[name] = refCount{v.ref, v.cnt - 1}
		} else {
			ref := fm.repo[name].ref
			delete(fm.repo, name)
			fnt := ref.(font.Disposer)
			fnt.Dispose()

			log.Println("dispose font:", name)
		}
	}
}

func (fm *FontManager) Get(name string) (fnt font.Font, ok bool) {
	if v, ok := fm.repo[name]; ok {
		fnt, ok = v.ref.(font.Font)
	}
	return
}
