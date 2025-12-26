package main

import (
	"fmt"
	// "strings"

	"FirstPersonShooter/gogl"

	"github.com/go-gl/gl/v3.3-core/gl"
	// "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
	// "FirstPersonShooter/menu"
)

const width, height = 800, 600

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	window, err := sdl.CreateWindow("Triangle", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)

	if err != nil {
		panic(err.Error())

	}
	window.GLCreateContext()
	defer window.Destroy()
	gl.Init()
	// if !menu.RunMainMenu(window) {
	// 	return
	// }
	fmt.Println(gogl.GetVersion())

	shaderProgram, err := gogl.NewShader("./Shaders/Vert.vert", "./Shaders/quadtexture.frag")
	if err != nil {
		panic(err)
	}
	texture := gogl.LoadTextureAlpha("./assets/te.png")

	// vertices := []float32{
	// 	0.5, 0.5, 0.0, 1.0, 1.0,
	// 	0.5, -0.5, 0.0, 1.0, 0.0,
	// 	-0.5, -0.5, 0.0, 0.0, 0.0,
	// 	-0.5, 0.5, 0.0, 0.0, 1.0,
	// }

	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0}

	// indices := []uint32{
	// 	0, 1, 3,
	// 	1, 2, 3,
	// }
	// cubePos := []mgl32.Vec3{
	// 	mgl32.Vec3{0.0, 0.0, 0.0},
	// 	mgl32.Vec3{2.0, 5.0, -10.0},
	// }

	gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	VAO := gogl.GenBindVertexArray()

	gogl.BufferDataFloat(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	// gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)
	// gogl.BufferDataInt(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	gogl.UnBindVertexArray()

	var x float32 = 0.0

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// gogl.UseProgram(shaderProgram)
		shaderProgram.Use()
		shaderProgram.SetFloat("x", x)
		shaderProgram.SetFloat("y", 0.0)
		gogl.BindTexture(texture)
		gogl.BindVertexArray(VAO)
		// projectionMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 100.0)
		// viewMatrix := mgl32.Ident4()
		// viewMatrix = mgl32.Translate3D(0.0, 0.0, -3.0)
		// shaderProgram.SetMat4("projection", projectionMatrix)
		// shaderProgram.SetMat4("view", viewMatrix)
		// gogl.BindTexture(tex)

		// for i,pos := range cubePos{
		// 	modelMatrix:= mgl32.Ident4()
		// 	modelMatrix = mgl32.Translate3D(pos.X(),pos.Y(),pos.Z()). Mul4(modelMatrix)
		// 	// angle:= 20.0* float(i)

		// 	shaderProgram.SetMat4("model",modelMatrix)
		// }

		// gl.DrawArrays(gl.TRIANGLES, 0, 3)
		// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		window.GLSwap()
		shaderProgram.CheckShadersForChanges()
		x += 0.1

	}

}
