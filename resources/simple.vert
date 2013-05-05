#version 330

layout(location = 0) in vec4 vertexPos;
layout(location = 1) in vec4 vertexNormal;

out vec4 fragNormal;

uniform mat4 pvm;

void main() {
	fragNormal = vertexNormal;
	gl_Position = pvm * vertexPos;
}
