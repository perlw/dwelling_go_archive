#version 330

in vec3 eyeNormal;
out vec4 fragment;

vec3 lightDir = vec3(0.0, 1.0, 0.0);
vec4 ambient = vec4(0.25, 0.25, 0.25, 1.0);
vec4 material = vec4(1.0, 1.0, 1.0, 1.0);


void main() {
	float NdotL = max(dot(eyeNormal, lightDir), 0.0);

	fragment = ambient + (material * NdotL);
}
