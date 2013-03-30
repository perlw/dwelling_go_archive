package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/jteeuwen/glfw"
	"runtime"
)

func main() {
	var (
		rotx, roty gl.Float
	)

	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		fmt.Printf("glfw: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(640, 480, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		fmt.Printf("glfw: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(0)
	glfw.SetWindowTitle("AH MAH GAHD IT WORKZ")

	if err := gl.Init(); err != nil {
		fmt.Printf("gl: %s\n", err)
	}

	//gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.DEPTH_TEST)

	gl.ClearColor(0.5, 0.5, 0.5, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	gl.Viewport(0, 0, 640, 480)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.ClearColor(0.5, 0.5, 1.0, 1.0)
	for glfw.WindowParam(glfw.Opened) == 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Translated(0, 0, -3.0)
		gl.Rotatef(rotx, 1, 0, 0)
		gl.Rotatef(roty, 0, 1, 0)

		rotx += 0.5
		roty += 0.5

		gl.Begin(gl.TRIANGLES)
		gl.Normal3f(0.0, 0.0, 1.0)
		gl.Color4f(1.0, 0.0, 0.0, 1.0)
		gl.Vertex3f(0.0, 1.0, 0.0)
		gl.Color4f(0.0, 1.0, 0.0, 1.0)
		gl.Vertex3f(0.5, 0.0, 0.0)
		gl.Color4f(0.0, 0.0, 1.0, 1.0)
		gl.Vertex3f(-0.5, 0.0, 0.0)
		gl.End()

		glfw.SwapBuffers()
	}
}
