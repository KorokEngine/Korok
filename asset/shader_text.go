package asset


// GLShader for TextRender
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
