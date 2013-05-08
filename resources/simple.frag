#version 330

in vec3 eyeNormal;
out vec4 fragment;

vec3 lightDir = vec3(1.0, 1.0, 0.0);
vec4 lightColor = vec4(1.0, 1.0, 1.0, 1.0);

in float height;

void main() {
	float NdotL = max(dot(eyeNormal, lightDir), 0.0);

	vec4 color = lightColor;
	float modifier = 4.0 * height;
	color.r = clamp(min(modifier - 1.5, -modifier + 4.5), 0.0, 1.0);
	color.g = clamp(min(modifier - 0.5, -modifier + 3.5), 0.0, 1.0);
	color.b = clamp(min(modifier + 0.5, -modifier + 2.5), 0.0, 1.0);
	vec4 ambient = color / 4;
	fragment = ambient + (color * NdotL);
}
