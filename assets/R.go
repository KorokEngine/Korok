package assets

import (
	"log"
	"os"
	"fmt"
	"image"


	"github.com/go-gl/gl/v3.2-core/gl"
	"korok/gfx"
	"korok/gfx/text"
)

var textures = make(map[string]*gfx.Texture2D)
var shaders  = make(map[string]*gfx.GLShader)
var fonts    = make(map[string]*text.Font)

func LoadShader() {
	shader, err := loadShader(vertex, color)
	if err != nil {
		log.Println(err)
	}
	shaders["dft"] = shader

	shader, err = loadShader(bVertex, bColor)
	if err != nil {
		log.Println(err)
	}
	shaders["batch"] = shader

	shader, err = loadShader(pVertex, pColor)
	if err != nil {
		log.Println(err)
	}
	shaders["particle"] = shader

	shader, err = loadShader(tVertex, tColor)
	if err != nil {
		fmt.Println(err)
	}
	shaders["text"] = shader
}

func GetShader(key string) *gfx.GLShader {
	return shaders[key]
}

func LoadTexture(file string) {
	texture, err := loadTexture(file)
	if err != nil {
		log.Println(err)
	}
	textures[file] = texture
}

func GetTexture(file string) *gfx.Texture2D  {
	return textures[file]
}

func LoadFont(img, fc string) {
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
	f, err := text.LoadBitmap(ir, fcr, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("load font sucess...")
	fonts[fc] = f
}

func GetFont(fc string) *text.Font{
	return fonts[fc]
}

func Clear()  {
	for _, v := range shaders {
		gl.DeleteProgram(v.Program)
	}

	for _, v := range textures {
		gl.DeleteTextures(1, &v.Id)
	}
}

func loadShader(vertex, fragment string) (*gfx.GLShader, error){
	shader := &gfx.GLShader{}
	program, err := gfx.Compile(vertex, fragment)
	if err != nil {
		return nil, err
	}
	shader.Program = program
	return shader, nil
}

func loadTexture(file string)(*gfx.Texture2D, error)  {
	log.Println("load file:" + file)
	// 1. load file
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found: %v", file, err)
	}
	// 2. decode image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	// 3. create
	texture, err :=  gfx.Generate(float32(img.Bounds().Dx()), float32(img.Bounds().Dy()), img)
	if err != nil {
		return nil, err
	}
	return texture, nil
}


