#version 330 core

	out vec4 FragColor;
	in vec2 TexCoord;
	uniform sampler2D texture1;
	
	void main(){
	// FragColor =  vec4(.0f,.0f,1.0f,1.0f);
	// vec4 tex = texture(texture1, TexCoord);
	FragColor = texture(texture1, TexCoord);
	// FragColor = vec4(tex.x,tex.y,1.0,1.0);
	
	}