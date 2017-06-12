package gfx



import (
	"fmt"
	"image"
	_ "image/png"
	_ "image/jpeg"
	"image/draw"

	"github.com/go-gl/gl/v3.2-core/gl"
	"log"
)

/**
	处理纹理相关问题
 */

type Texture2D struct{
	Width, Height float32
	Id uint32
}

func Generate(width, height float32, image image.Image)  (*Texture2D, error){
	texture := &Texture2D{Width:width, Height:height}
	id, err:= newTexture(image)
	if err != nil {
		return nil, err
	}
	texture.Id = id
	return texture, nil
}

func (t *Texture2D)Bind()  {
	gl.BindTexture(gl.TEXTURE_2D, t.Id)
}


func newTexture(img image.Image) (uint32, error) {
	// 3. copy image
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	r, g, b, a := rgba.At(5, 5).RGBA()
	log.Printf("rgba file, R:%x G:%x B:%x A:%x", r, g, b, a)
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

