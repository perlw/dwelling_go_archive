#version 330

layout(location = 0) in vec4 vertexPos;

out vec3 eyeNormal;

uniform mat4 pv;
uniform mat4 model;
uniform vec3 normal;

void main() {
	mat4 pvm = pv * model;
	eyeNormal = normal;
	gl_Position = pvm * vertexPos;
}
