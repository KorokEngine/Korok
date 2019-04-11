// +build !android,!ios,!windows,!js

package asset

// shader for batch-system
var bVertex = `
#version 330

uniform mat4 proj;

in vec4 xyuv;
in vec4 rgba;

out vec4 outColor;
out vec2 outTexCoord;

void main() {
    outColor = rgba;
	outTexCoord = xyuv.zw;
    gl_Position = proj * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var bColor = `
#version 330

uniform sampler2D tex;

in vec2 outTexCoord;
in vec4 outColor;

out vec4 outputColor;
void main() {
    outputColor = texture(tex, outTexCoord) * outColor;
}
` + "\x00"

// mesh shader

var vertex = `
#version 330

uniform mat4 proj;
uniform mat4 model;

in vec4 xyuv;
in vec4 rgba;

out vec4 outColor;
out vec2 outTexCoord;

void main() {
    outColor = rgba;
	outTexCoord = xyuv.zw;

    gl_Position = proj * model * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var color = `
#version 330

uniform sampler2D tex;

in vec2 outTexCoord;
in vec4 outColor;

out vec4 outputColor;
void main() {
    outputColor = texture(tex, outTexCoord) * outColor;
}
` + "\x00"
