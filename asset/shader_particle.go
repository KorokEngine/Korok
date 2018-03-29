package asset

// GLShader for Particle-System
// 粒子采用 GPU 计算的方式
// inVertex: <x,y,scale,rot>
// inColor: <r,g,b,a>
//
var pVertex = `
	#version 330
	uniform mat4 projection;
	uniform vec4 model[4];			// <x,y,u,v>

	layout (location = 0) in vec4 ver;  	// <x,y, size, rot>
	layout (location = 1) in vec4 color;    // <r, g, b, a>
	layout (location = 2) in float  index;    // <i>

	out vec2 fragTexCoord;
	out vec4 fragColor;

	void main() {
	    fragTexCoord = model[int(index)].zw;
	    fragColor = color;
	    gl_Position = projection * vec4(model[int(index)].xy * ver.z + ver.xy, 0, 1);
	}
	` + "\x00"

var pColor = `
	#version 330
	uniform sampler2D tex;

	in vec2 fragTexCoord;
	in vec4 fragColor;
	out vec4 outputColor;

	void main() {
	    outputColor = texture(tex, fragTexCoord);
	}
	` + "\x00"
