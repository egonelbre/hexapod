#version 330

// Input vertex attributes (from vertex shader)
in vec3 fragPosition;
in vec3 fragCameraPosition;
in vec2 fragTexCoord;
in vec3 fragNormal;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform vec3 viewPos;

// Output fragment color
out vec4 finalColor;

// NOTE: Add here your custom variables

const vec3 directionalLight = vec3(1, 1, 1);

void main()
{
   // Texel color fetching from texture sampler
    vec4 texelColor = texture(texture0, fragTexCoord)*colDiffuse*fragColor;
    float cosTheta = clamp(dot(fragNormal, directionalLight), 0.5, 1);

    // Calculate final fragment color
    vec3 color = texelColor.rgb * cosTheta;

    finalColor = vec4(color, 1);
}