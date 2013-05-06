#version 330

in vec3 eyeNormal;
in float height;
out vec4 fragment;

vec3 lightDir = vec3(0.0, 1.0, 0.0);
//vec4 ambient = vec4(0.25, 0.25, 0.25, 1.0);
vec4 lightColor = vec4(1.0, 1.0, 1.0, 1.0);

uniform int pyramid;
uniform float wave;

void main() {
	float NdotL = max(dot(eyeNormal, lightDir), 0.0);

	vec4 color = lightColor;
	float modifier = 4.0 * (height * wave) / pyramid;
	color.r = clamp(min(modifier - 1.5, -modifier + 4.5), 0.0, 1.0);
	color.g = clamp(min(modifier - 0.5, -modifier + 3.5), 0.0, 1.0);
	color.b = clamp(min(modifier + 0.5, -modifier + 2.5), 0.0, 1.0);
	vec4 ambient = color / 4;
	fragment = ambient + (color * NdotL);
}
