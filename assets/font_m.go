package assets

import (
	"os"
	"fmt"
	"korok/gfx/font"
)

type FontManager struct {
	repo map[string]RefCount
}

func (fm *FontManager) Load(img, fc string) {
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
	fmt.Println("load font sucess...", f)
}

func (fm *FontManager) GetFont(fc string) {
	// return fm.repo[fc]
	return
}
