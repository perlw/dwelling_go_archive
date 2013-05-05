#version 330

in vec4 fragNormal;
out vec4 fragment;

float correctNormalColor(float val) {
	if (val < 0.0) {
		return -val;
	}

	return val;
}

void main() {
	vec4 color = vec4(correctNormalColor(fragNormal.x), correctNormalColor(fragNormal.y), correctNormalColor(fragNormal.z), 1.0);
	fragment = color;
}
