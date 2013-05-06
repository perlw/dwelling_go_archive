#version 330

layout(location = 0) in vec4 vertexPos;
//layout(location = 1) in vec3 vertexNormal;

out vec3 eyeNormal;
out float height;

uniform mat4 pv;
uniform mat4 model;
uniform vec3 normal;

void main() {
	height = vertexPos.y;
	mat4 pvm = pv * model;
	eyeNormal = normalize(vec3(pvm * vec4(normal, 0.0)));
	gl_Position = pvm * vertexPos;
}
