#version 130

in vec3 eyeNormal;
out vec4 fragment;

vec3 lightDir = vec3(-1.0, 1.0, -1.0);
vec4 lightColor = vec4(1.0, 1.0, 1.0, 1.0);

uniform int mouseHit;
uniform vec3 flatColor;
uniform int skipLight;

void main() {
	if (skipLight == 0) {
		float NdotL = max(dot(eyeNormal, lightDir), 0.0);

		vec4 color = (lightColor + vec4(eyeNormal, 1.0)) / 4;
		if (mouseHit > 0) {
			//color.g = 1.0f;
		}
		vec4 ambient = color;

		fragment = ambient + (color * NdotL);
	} else {
		fragment = vec4(flatColor, 1.0f);
	}
}
