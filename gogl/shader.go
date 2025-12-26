package gogl

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	id             ProgramID
	vertexPath     string
	fragmentPath   string
	vertexmodified time.Time
	fragmodified   time.Time
}

// var loadedShaders = make(map[ProgramID]*Shader)

func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	pid, err := CreateProgram(vertexPath, fragmentPath)
	if err != nil {
		return nil, err
	}

	result := &Shader{pid, vertexPath, fragmentPath, getModifiedTime(vertexPath), getModifiedTime(fragmentPath)}

	// loadedShaders[pid] = result

	return result, nil
}

func (shader *Shader) Use() {
	// shader := loadedShaders[id]
	UseProgram(shader.id)
}

func (shader *Shader) SetFloat(name string, f float32) {
	name_cstring := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(uint32(shader.id), name_cstring)

	gl.Uniform1f(location, f)
}

func (shader *Shader) SetMat4(name string, mat mgl32.Mat4) {
	name_cstring := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(uint32(shader.id), name_cstring)
	m4 := [16]float32(mat)
	gl.UniformMatrix4fv(location, 1, false, &m4[0])
}

func getModifiedTime(filePath string) time.Time {
	file, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	return file.ModTime()

}

func (shader *Shader) CheckShadersForChanges() {

	vertexmodTime := getModifiedTime(shader.vertexPath)
	fragmodTime := getModifiedTime(shader.vertexPath)
	if !vertexmodTime.Equal(shader.vertexmodified) || !fragmodTime.Equal(shader.fragmodified) {
		pid, err := CreateProgram(shader.vertexPath, shader.fragmentPath)
		if err != nil {
			fmt.Println(err)
		} else {
			gl.DeleteProgram(uint32(shader.id))
			shader.id = pid

		}
	}

}

func (shader *Shader) SetVec3(name string, x, y, z float32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(uint32(shader.id), name_cstr)
	gl.Uniform3f(location, x, y, z)
}
