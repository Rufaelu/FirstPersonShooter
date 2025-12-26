package gogl

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"strings"

	// "time"

	"github.com/go-gl/gl/v3.3-core/gl"
	// "github.com/veandco/go-sdl2/sdl"
)

type ShaderID uint32
type ProgramID uint32
type BufferID uint32
type TextureID uint32

// type VAOID uint32
// type VBOID uint32
// type BufferID uint32

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))

}

func CreateShader(shaderSource string, shaderType uint32) (ShaderID, error) {
	shaderid := gl.CreateShader(shaderType)
	shaderSource += "\x00"

	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderid, 1, csource, nil)

	free()

	gl.CompileShader(shaderid)
	var status int32
	gl.GetShaderiv(shaderid, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var loglength int32
		gl.GetShaderiv(shaderid, gl.INFO_LOG_LENGTH, &loglength)
		log := strings.Repeat("\x00", int(loglength+1))
		gl.GetShaderInfoLog(shaderid, loglength, nil, gl.Str(log))
		return 0, errors.New("Failed To Compile shader: \n" + log)
	}
	return ShaderID(shaderid), nil
}

func LoadShader(path string, shaderType uint32) (ShaderID, error) {

	shaderFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	shaderFileStr := string(shaderFile)
	shaderid, err := CreateShader(shaderFileStr, shaderType)
	if err != nil {
		return 0, err
	}

	return shaderid, nil

}

func CreateProgram(vertpath, fragpath string) (ProgramID, error) {

	vert, err := LoadShader(vertpath, gl.VERTEX_SHADER)
	if err != nil {

		return 0, err
	}
	frag, err := LoadShader(fragpath, gl.FRAGMENT_SHADER)
	if err != nil {

		return 0, err
	}
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.LinkProgram(shaderProgram)
	var sucess int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &sucess)

	if sucess == gl.FALSE {
		var loglength int32
		gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &loglength)
		log := strings.Repeat("\x00", int(loglength+1))
		gl.GetProgramInfoLog(shaderProgram, loglength, nil, gl.Str(log))
		return 0, fmt.Errorf(fmt.Sprintf("Failed To Link Program: %s\n", log))

	}
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	return ProgramID(shaderProgram), nil

}

func GenBindBuffer(target uint32) BufferID {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)
	return BufferID(buffer)
}
func GenEBO() BufferID {
	var EBO uint32
	gl.GenBuffers(1, &EBO)
	return BufferID(EBO)
}

func GenBindVertexArray() BufferID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return BufferID(VAO)
}

func BindVertexArray(vaid BufferID) {
	gl.BindVertexArray(uint32(vaid))

}

func BufferDataFloat(target uint32, data []float32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)

}
func BufferDataInt(target uint32, data []uint32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)

}

func UnBindVertexArray() {
	gl.BindVertexArray(0)
}

func UseProgram(pId ProgramID) {
	gl.UseProgram(uint32(pId))

}

func LoadTextureAlpha(filename string) TextureID {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}

	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {

		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 246)
			bIndex++
		}
	}

	texture := GenBindTexture()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}
func BindTexture(id TextureID) {
	gl.BindTexture(gl.TEXTURE_2D, uint32(id))

}
func GenBindTexture() TextureID {
	var texid uint32
	gl.GenTextures(1, &texid)
	gl.BindTexture(gl.TEXTURE_2D, texid)
	return TextureID(texid)

}
