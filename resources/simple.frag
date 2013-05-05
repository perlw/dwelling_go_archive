#version 330

in vec4 fragNormal;
out vec4 outputF;

void main() {
	outputF = fragNormal;
}
