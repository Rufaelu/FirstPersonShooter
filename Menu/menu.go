package menu

import (
	"time"
	"unsafe"

	"FirstPersonShooter/gogl"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const winWidth, winHeight = 800, 600

func RunMainMenu(window *sdl.Window) bool {
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer ttf.Quit()

	font, err := ttf.OpenFont("assets/Doto-Regular.ttf", 28)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	menuShader, err := gogl.NewShader(
		"shaders/menu.vert",
		"shaders/menu.frag",
	)
	if err != nil {
		panic(err)
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {

			case *sdl.QuitEvent:
				return false

			case *sdl.MouseButtonEvent:
				if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == sdl.BUTTON_LEFT {
					mx := float32(e.X)
					my := float32(e.Y)

					if inside(mx, my, 490, 260, 300, 70) {
						time.Sleep(150 * time.Millisecond)
						return true
					}

					if inside(mx, my, 490, 360, 300, 70) {
						return false
					}
				}
			}
		}

		gl.ClearColor(0.08, 0.08, 0.12, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		DrawRect(menuShader, 490, 260, 300, 70, 0.2, 0.7, 0.2)
		DrawRect(menuShader, 490, 360, 300, 70, 0.7, 0.2, 0.2)
		DrawText("START", 585, 305, font, 255, 255, 255, window)
		DrawText("EXIT", 600, 405, font, 255, 255, 255, window)

		window.GLSwap()
		sdl.Delay(16)
	}
}
func DrawRect(shader *gogl.Shader, x, y, w, h float32, r, g, b float32) {
	nx := func(px float32) float32 {
		return (px/float32(winWidth))*2 - 1
	}
	ny := func(py float32) float32 {
		return 1 - (py/float32(winHeight))*2
	}

	vertices := []float32{
		nx(x), ny(y),
		nx(x + w), ny(y),
		nx(x + w), ny(y + h),
		nx(x), ny(y + h),
	}

	indices := []uint32{0, 1, 2, 2, 3, 0}

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	shader.Use()
	shader.SetVec3("color", r, g, b)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}
func inside(mx, my, x, y, w, h float32) bool {
	return mx >= x && mx <= x+w && my >= y && my <= y+h
}
func DrawText(
	text string,
	x, y int32,
	font *ttf.Font,
	r, g, b uint8,
	window *sdl.Window,
) {
	color := sdl.Color{R: r, G: g, B: b, A: 255}

	surface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return
	}
	defer surface.Free()

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(surface.W),
		int32(surface.H),
		0,
		gl.BGRA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&surface.Pixels()[0]),
	)

	// Convert to NDC
	nx := float32(x)/float32(winWidth)*2 - 1
	ny := 1 - float32(y)/float32(winHeight)*2
	nw := float32(surface.W) / float32(winWidth) * 2
	nh := float32(surface.H) / float32(winHeight) * 2

	vertices := []float32{
		nx, ny,
		nx + nw, ny,
		nx + nw, ny - nh,
		nx, ny - nh,
	}

	indices := []uint32{0, 1, 2, 2, 3, 0}

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

	gl.DeleteTextures(1, &texture)
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteBuffers(1, &ebo)
}
