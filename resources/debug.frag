#version 130

out vec4 fragment;
uniform vec3 flatColor;

void main() {
	fragment = vec4(flatColor, 1.0f);
}
