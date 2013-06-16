#version 130

in vec4 vertexPos;

uniform mat4 pv;
uniform mat4 model;

void main() {
	mat4 pvm = pv * model;
	gl_Position = pvm * vertexPos;
}
