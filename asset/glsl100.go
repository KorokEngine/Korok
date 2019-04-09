// +build android ios windows js

package asset

// shader for batch-system
var bVertex = `
#version 100

uniform mat4 proj;

attribute vec4 xyuv;
attribute vec4 rgba;

varying vec4 outColor;
varying vec2 outTexCoord;

void main() {
    outColor = rgba;
	outTexCoord = xyuv.zw;
    gl_Position = proj * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var bColor = `
#version 100

#ifdef GL_ES
precision mediump float;
#endif

uniform sampler2D tex;

varying vec2 outTexCoord;
varying vec4 outColor;

void main() {
    gl_FragColor = texture2D(tex, outTexCoord) * outColor;
}
` + "\x00"

// mesh shader

var vertex = `
#version 100

uniform mat4 proj;
uniform mat4 model;

attribute vec4 xyuv;
attribute vec4 rgba;

varying vec4 outColor;
varying vec2 outTexCoord;

void main() {
    outColor = rgba;
	outTexCoord = xyuv.zw;
    gl_Position = proj * model * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var color = `
#version 100

#ifdef GL_ES
precision mediump float;
#endif

uniform sampler2D tex;

varying vec2 outTexCoord;
varying vec4 outColor;

void main() {
    gl_FragColor = texture2D(tex, outTexCoord) * outColor;
}
` + "\x00"
