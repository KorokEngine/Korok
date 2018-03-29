package asset

// mesh shader

var vertex = `
#version 330

uniform mat4 proj;
// uniform mat4 camera;
uniform mat4 model;

in vec4 xyuv;
in vec4 rgba;

out vec4 outColor;
out vec2 fragTexCoord;

void main() {
    outColor = rgba;
	fragTexCoord = xyuv.zw;
    gl_Position = proj * model * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var color = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec4 outColor;
out vec4 outputColor;
void main() {
    outputColor = texture(tex, fragTexCoord) * outColor;
}
` + "\x00"