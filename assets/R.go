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
var shaders  = make(map[string]*gfx.Shader)
var fonts    = make(map[string]*text.Font)

func LoadShader() {
	shader, err := loadShader(vertex, color)
	if err != nil {
		log.Println(err)
	}
	shaders["dft"] = shader

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

func GetShader(key string) *gfx.Shader  {
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

func loadShader(vertex, fragment string) (*gfx.Shader, error){
	shader := &gfx.Shader{}
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

// 缺省的Shader写在这里

var vertex =
	`
	#version 330
	uniform mat4 projection;
	uniform mat4 model;

	layout (location = 0) in vec4 vert; // <vec2 pos, vec2 tex>

	out vec2 fragTexCoord;
	void main() {
	    fragTexCoord = vert.zw;
	    gl_Position = projection * model * vec4(vert.xy, 0, 1);
	}
	` + "\x00"

var color =
	`
	#version 330
	uniform sampler2D tex;
	in vec2 fragTexCoord;
	out vec4 outputColor;

	void main() {
	    outputColor = texture(tex, fragTexCoord);

	    //if (outputColor.a == 0.0 && outputColor.r == 0.0 && outputColor.g == 0.0 && outputColor.b == 0.0) {
	    //	outputColor = vec4(1, 0, 0, 1);
	    //}
	}
	` + "\x00"

// Shader for Particle-System
var pVertex = `
	#version 330 core
	layout (location = 0) in vec4 vertex; // <vec2 position, vec2 texCoords>

	out vec2 TexCoords;
	out vec4 ParticleColor;

	uniform mat4 projection;
	uniform vec2 offset;
	uniform vec4 color;

	void main()
	{
	    float scale = 10.0f;
	    TexCoords = vertex.zw;
	    ParticleColor = color;
	    gl_Position = projection * vec4((vertex.xy * scale) + offset, 0.0, 1.0);
	}
` + "\x00"

var pColor = `
	#version 330 core
	in vec2 TexCoords;
	in vec4 ParticleColor;
	out vec4 color;

	uniform sampler2D sprite;

	void main()
	{
	    color = (texture(sprite, TexCoords) * ParticleColor);
	}
` + "\x00"

// Shader for TextRender
var tVertex = `
	#version 330 core
	layout (location = 0) in vec4 vertex; // <vec2 pos, vec2 tex>
	out vec2 TexCoords;

	uniform mat4 projection;
	uniform vec3 model;					  // <x,y, scale>

	void main()
	{
	    gl_Position = projection * vec4(vertex.x + model.x, vertex.y + model.y, 0.0, 1.0);
	    TexCoords = vertex.zw;
	}
	` + "\x00"

var tColor = `
	#version 330 core
	in vec2 TexCoords;
	out vec4 color;

	uniform sampler2D text;
	uniform vec3 textColor;

	void main()
	{
	    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(text, TexCoords).r);
	    color = vec4(textColor, 1.0) * sampled;
	}
	` + "\x00"

