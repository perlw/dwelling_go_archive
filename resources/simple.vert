#version 330

layout(location = 0) in vec4 vertexPos;
layout(location = 1) in vec3 vertexNormal;

out vec3 eyeNormal;

uniform mat4 pvm;

void main() {
	eyeNormal = normalize(vec3(pvm * vec4(vertexNormal, 0.0)));
	gl_Position = pvm * vertexPos;
}
