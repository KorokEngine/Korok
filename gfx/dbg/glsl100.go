// +build android ios windows

package dbg

var vsh = `
#version 100

uniform mat4 projection;

attribute vec4 xyuv;
attribute vec4 rgba;

varying vec4 outColor;
varying vec2 outTexCoord;

void main() {
    outColor = rgba;
	outTexCoord = xyuv.zw;
    gl_Position = projection * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var fsh = `
#version 100

#ifdef GL_ES
	precision mediump float;
#endif

uniform sampler2D tex;

varying vec4 outColor;
varying vec2 outTexCoord;

void main() {
	gl_FragColor = outColor * texture2D(tex, outTexCoord);
}
` + "\x00"
