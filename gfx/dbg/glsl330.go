// +build !android,!ios,!windows

package dbg

var vsh = `
#version 330

uniform mat4 projection;

in vec4 xyuv;
in vec4 rgba;

out vec4 outColor;
out vec2 outTexCoord;

void main() {
    outColor = rgba;
	outTexCoord = xyuv.zw;

    gl_Position = projection * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var fsh = `
#version 330

uniform sampler2D tex;

in vec2 outTexCoord;
in vec4 outColor;

out vec4 outputColor;
void main() {
	if (outTexCoord.x == 2.0) {
		outputColor = outColor;
	} else {
	    outputColor = outColor * texture(tex, outTexCoord);
	}
}
` + "\x00"
