package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
	"io/ioutil"
	"runtime"
	"unsafe"
)

import (
	"dwelling/matrix"
)

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		fmt.Printf("glfw: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 3)
	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(640, 480, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		fmt.Printf("glfw: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(0)
	glfw.SetWindowTitle("Dwelling")

	if err := gl.Init(); err != nil {
		fmt.Printf("gl: %s\n", err)
	}

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.5, 0.5, 0.5, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Viewport(0, 0, 640, 480)

	viewMatrix := matrix.NewIdentityMatrix()
	viewMatrix.Translate(4.0, 0.0, 0.0)
	viewMatrix.RotateY(-45.0)
	projMatrix := matrix.NewPerspectiveMatrix(53.13, 640.0/480.0, 0.1, 1000.0)
	fmt.Println(projMatrix)
	fmt.Println(viewMatrix)

	var vao gl.Uint
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	sizeFloat := int(unsafe.Sizeof([1]float32{}))
	sizeInt := int(unsafe.Sizeof([1]int{}))
	vertexData := []float32{
		// Front
		-0.5, -0.5, 0.5, // 0
		0.5, -0.5, 0.5, // 1
		0.5, 0.5, 0.5, // 2
		-0.5, 0.5, 0.5, // 3

		// Back
		0.5, -0.5, -0.5, // 4
		-0.5, -0.5, -0.5, // 5
		-0.5, 0.5, -0.5, // 6
		0.5, 0.5, -0.5, // 7
	}
	normalData := []float32{
		// Front
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,

		// Back
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
	}
	indexData := []uint32{
		// Front
		0, 1, 2,
		2, 3, 0,

		// Back
		4, 5, 6,
		6, 7, 4,
	}
	vertexBuffer := makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&vertexData[0]), sizeFloat*len(vertexData)) // 4 == sizeof float32
	normalBuffer := makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&normalData[0]), sizeFloat*len(normalData))
	indexBuffer := makeBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Pointer(&indexData[0]), sizeInt*len(indexData))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	program := readShaders()

	gl.UseProgram(program)

	pvm := gl.GLString("pvm")
	defer gl.GLStringFree(pvm)
	pvmId := gl.GetUniformLocation(program, pvm)

	pvmMatrix := matrix.MultiplyMatrix(projMatrix, viewMatrix)

	rot := 0.0
	gl.ClearColor(0.5, 0.5, 1.0, 1.0)
	for glfw.WindowParam(glfw.Opened) == 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		rot += 0.005
		if rot >= 360 {
			rot = 0.0
		}
		modelMatrix := matrix.NewIdentityMatrix()
		modelMatrix.Translate(0.0, 0.0, -5.0)
		modelMatrix.RotateY(rot)

		glPVMMatrix := matrixToGL(matrix.MultiplyMatrix(pvmMatrix, modelMatrix))

		gl.UseProgram(program)
		gl.UniformMatrix4fv(pvmId, 1, gl.FALSE, &glPVMMatrix[0])

		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
		gl.DrawElements(gl.TRIANGLES, gl.Sizei(len(indexData)), gl.UNSIGNED_INT, nil)

		if err := gl.GetError(); err != 0 {
			fmt.Printf("Err: %d\n", err)
			break
		}

		glfw.SwapBuffers()
	}
}

func makeBuffer(target gl.Enum, buffer_data gl.Pointer, size int) gl.Uint {
	var buffer gl.Uint

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)
	gl.BufferData(target, gl.Sizeiptr(size), buffer_data, gl.STATIC_DRAW)

	return buffer
}

func readShaders() gl.Uint {
	// Vertex shader
	vertexFile, err := ioutil.ReadFile("simple.vert")
	if err != nil {
		return 0
	}
	vertexSource := gl.GLString(string(vertexFile))
	defer gl.GLStringFree(vertexSource)

	vertexObj := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexObj, 1, &vertexSource, nil)
	gl.CompileShader(vertexObj)
	defer gl.DeleteShader(vertexObj)
	printShaderLog(vertexObj)

	// Fragment shader
	fragmentFile, err := ioutil.ReadFile("simple.frag")
	if err != nil {
		return 0
	}
	fragmentSource := gl.GLString(string(fragmentFile))
	defer gl.GLStringFree(fragmentSource)

	fragmentObj := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentObj, 1, &fragmentSource, nil)
	gl.CompileShader(fragmentObj)
	defer gl.DeleteShader(fragmentObj)
	printShaderLog(fragmentObj)

	// Program
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexObj)
	gl.AttachShader(program, fragmentObj)

	// Attribs
	fragmentOut := gl.GLString("outputF")
	defer gl.GLStringFree(fragmentOut)
	gl.BindFragDataLocation(program, 0, fragmentOut)

	gl.LinkProgram(program)

	printProgramLog(program)

	return program
}

func printShaderLog(shader gl.Uint) {
	var length gl.Int
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
	glString := gl.GLStringAlloc(gl.Sizei(length))
	defer gl.GLStringFree(glString)
	gl.GetShaderInfoLog(shader, gl.Sizei(length), nil, glString)
	fmt.Println("shader log: ", gl.GoString(glString))
}

func printProgramLog(program gl.Uint) {
	var length gl.Int
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
	glString := gl.GLStringAlloc(gl.Sizei(length))
	defer gl.GLStringFree(glString)
	gl.GetProgramInfoLog(program, gl.Sizei(length), nil, glString)
	fmt.Println("program log: ", gl.GoString(glString))
}

func matrixToGL(matrix *matrix.Matrix) [16]gl.Float {
	var newMatrix [16]gl.Float

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			newMatrix[(x*4)+y] = gl.Float(matrix.Values[(y*4)+x])
		}
	}

	return newMatrix
}
