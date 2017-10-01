package bk

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

/**
处理纹理相关问题
*/

type Texture2D struct {
	Width, Height float32
	Id            uint32
}

func Generate(width, height float32, image image.Image) (*Texture2D, error) {
	texture := &Texture2D{Width: width, Height: height}
	id, err := newTexture(image)
	if err != nil {
		return nil, err
	}
	texture.Id = id
	return texture, nil
}

func (t *Texture2D) Bind(stage int32) {
	gl.BindTexture(gl.TEXTURE_2D, t.Id)
}

func (t *Texture2D) Sub(x, y float32, w, h float32) *SubTex {
	subTex := &SubTex{Texture2D: t}
	subTex.Min = mgl32.Vec2{x, y}
	subTex.Max = mgl32.Vec2{x + w, y + h}
	return subTex
}

func (t *Texture2D) Destroy() {
	// TODO impl
}

func newTexture(img image.Image) (uint32, error) {
	// 3. copy image
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	// 4. upload texture
	var texture uint32
	// 4.1 apply space
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// 4.2 params
	// 大小插值
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// 环绕方式
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// 4.3 upload
	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Dx()),
		int32(rgba.Rect.Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	return texture, nil
}

///// 还需要抽象 SubTexture 的概念出来
type SubTex struct {
	*Texture2D

	// location -
	Min, Max mgl32.Vec2
}
