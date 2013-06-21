#version 130

in vec4 vertexPos;
in float occFactor;

out vec3 eyeNormal;
out float occFac;

uniform mat4 pv;
uniform mat4 model;
uniform vec3 normal;

void main() {
	occFac = occFactor;
	mat4 pvm = pv * model;
	eyeNormal = normal;
	gl_Position = pvm * vertexPos;
}
