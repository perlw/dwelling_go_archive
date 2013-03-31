#version 330

layout(location = 0) in vec4 vert;

uniform mat4 view;
uniform mat4 proj;
uniform mat4 model;
mat4 pvm = view * proj;

void main() {
	gl_Position = view * proj * (vert * model);
}
