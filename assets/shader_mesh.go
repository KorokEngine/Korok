package assets

// mesh shader
var vertex = `
	#version 330
	uniform mat4 projection;
	uniform mat4 model;

	layout (location = 0) in vec4 vert;  // <vec2 pos, vec2 tex>

	out vec2 fragTexCoord;
	void main() {
	    fragTexCoord = vert.zw;
	    gl_Position = projection * model * vec4(vert.xy, 0, 1);
	}
	` + "\x00"

var color = `
	#version 330
	uniform sampler2D tex;

	in vec2 fragTexCoord;
	out vec4 outputColor;

	void main() {
	    outputColor = texture(tex, fragTexCoord);
	}
	` + "\x00"
