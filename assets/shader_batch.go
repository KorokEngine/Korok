package assets


// shader for batch-system
var bVertex = `
	#version 330
	uniform mat4 projection;

	layout (location = 0) in vec2 position;  // <vec2 pos, vec2 tex>
	layout (location = 1) in vec2 texCoord;  // <vec2 pos, vec2 tex>
	layout (location = 2) in vec4 color;     // <vec2 pos, vec2 tex>

	out vec2 fragTexCoord;
	out vec4 fragColor;

	void main() {
	    fragTexCoord = texCoord;
	    fragColor = color;
	    gl_Position = projection * vec4(position, 0, 1);
	}
	` + "\x00"

var bColor = `
	#version 330
	uniform sampler2D tex;

	in vec2 fragTexCoord;
	in vec4 fragColor;
	out vec4 outputColor;

	void main() {
	    outputColor = texture(tex, fragTexCoord);
	}
	` + "\x00"
