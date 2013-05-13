#version 330

in vec3 eyeNormal;
out vec4 fragment;

vec3 lightDir = vec3(1.0, 1.0, 1.0);
vec4 lightColor = vec4(1.0, 1.0, 1.0, 1.0);

void main() {
	float NdotL = max(dot(eyeNormal, lightDir), 0.0);

	vec4 color = (lightColor + vec4(eyeNormal, 1.0)) / 2;
	vec4 ambient = color / 2;
	fragment = ambient + (color * NdotL);
}

