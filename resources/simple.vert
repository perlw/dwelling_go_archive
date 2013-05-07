#version 330

layout(location = 0) in vec4 vertexPos;

out vec3 eyeNormal;
out float height;

uniform mat4 pv;
uniform mat4 model;
uniform vec3 normal;
uniform float maxHeight;

void main() {
	height = vertexPos.y / maxHeight;
	mat4 pvm = pv * model;
	eyeNormal = normalize(vec3(pvm * vec4(normal, 0.0)));
	gl_Position = pvm * vertexPos;
}
