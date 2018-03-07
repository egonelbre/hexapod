#version 330

// Input vertex attributes
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;
in vec4 vertexColor;

// Input uniform values
uniform mat4 mvp;
uniform mat4 mMatrix;

// Output vertex attributes (to fragment shader)
out vec3 fragPosition;
out vec3 fragCameraPosition;
out vec2 fragTexCoord;
out vec3 fragNormal;
out vec4 fragColor;

// NOTE: Add here your custom variables 

void main()
{
	// Calculate fragment normal based on normal transformations
    mat3 normalMatrix = transpose(inverse(mat3(mMatrix))); 

    // Send vertex attributes to fragment shader
    fragTexCoord = vertexTexCoord;
    fragColor = vertexColor;
    //fragPosition = vec3(mMatrix*vec4(vertexPosition, 1.0f));
    fragPosition = vec3(mMatrix*vec4(vertexPosition, 1.0f));
    fragNormal = normalize(normalMatrix*vertexNormal);
    
    // Calculate final vertex position
    fragCameraPosition = vec3(mvp*vec4(vertexPosition, 1.0));
    gl_Position = mvp*vec4(vertexPosition, 1.0);
}
