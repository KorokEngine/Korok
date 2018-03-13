package assets


// shader for batch-system
var bVertex = `
#version 330

uniform mat4 proj;
// uniform mat4 camera;

in vec4 xyuv;
in vec4 rgba;

out vec4 outColor;
out vec2 fragTexCoord;

void main() {
    outColor = rgba;
	fragTexCoord = xyuv.zw;
    gl_Position = proj * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var bColor = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec4 outColor;
out vec4 outputColor;
void main() {
    outputColor = texture(tex, fragTexCoord) * outColor;
}
` + "\x00"

