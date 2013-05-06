package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
	"io/ioutil"
	"math"
	"runtime"
	"time"
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
	viewMatrix.RotateX(-25)
	viewMatrix.RotateY(45)
	viewMatrix.Translate(-6.0, -5.0, -6.0)
	projMatrix := matrix.NewPerspectiveMatrix(53.13, 640.0/480.0, 0.1, 1000.0)
	fmt.Println(projMatrix)
	fmt.Println(viewMatrix)

	var vao gl.Uint
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	sizeFloat := int(unsafe.Sizeof([1]float32{}))
	vertexData := []float32{
		// Front
		-0.5, -0.5, 0.5, // 0
		0.5, -0.5, 0.5, // 1
		0.5, 0.5, 0.5, // 2
		0.5, 0.5, 0.5, // 3
		-0.5, 0.5, 0.5, // 4
		-0.5, -0.5, 0.5, // 5

		// Back
		0.5, 0.5, -0.5, // 6
		0.5, -0.5, -0.5, // 7
		-0.5, -0.5, -0.5, // 8
		-0.5, -0.5, -0.5, // 9
		-0.5, 0.5, -0.5, // 10
		0.5, 0.5, -0.5, // 11

		// Left
		-0.5, -0.5, -0.5, // 12
		-0.5, -0.5, 0.5, // 13
		-0.5, 0.5, 0.5, // 14
		-0.5, 0.5, 0.5, // 15
		-0.5, 0.5, -0.5, // 16
		-0.5, -0.5, -0.5, // 17

		// Right
		0.5, 0.5, 0.5, // 20
		0.5, -0.5, 0.5, // 19
		0.5, -0.5, -0.5, // 18
		0.5, -0.5, -0.5, // 23
		0.5, 0.5, -0.5, // 22
		0.5, 0.5, 0.5, // 21

		// Top
		0.5, 0.5, 0.5, // 24
		0.5, 0.5, -0.5, // 25
		-0.5, 0.5, -0.5, // 26
		-0.5, 0.5, -0.5, // 27
		-0.5, 0.5, 0.5, // 28
		0.5, 0.5, 0.5, // 29

		// Bottom
		-0.5, -0.5, -0.5, // 30
		0.5, -0.5, -0.5, // 31
		0.5, -0.5, 0.5, // 32
		0.5, -0.5, 0.5, // 33
		-0.5, -0.5, 0.5, // 34
		-0.5, -0.5, -0.5, // 35
	}
	normalData := []float32{
		// Front
		0.0, 0.0, 1.0, // 0
		0.0, 0.0, 1.0, // 1
		0.0, 0.0, 1.0, // 2
		0.0, 0.0, 1.0, // 3
		0.0, 0.0, 1.0, // 4
		0.0, 0.0, 1.0, // 5

		// Back
		0.0, 0.0, -1.0, // 6
		0.0, 0.0, -1.0, // 7
		0.0, 0.0, -1.0, // 8
		0.0, 0.0, -1.0, // 9
		0.0, 0.0, -1.0, // 10
		0.0, 0.0, -1.0, // 11

		// Left
		-1.0, 0.0, 0.0, // 12
		-1.0, 0.0, 0.0, // 13
		-1.0, 0.0, 0.0, // 14
		-1.0, 0.0, 0.0, // 15
		-1.0, 0.0, 0.0, // 16
		-1.0, 0.0, 0.0, // 17

		// Right
		1.0, 0.0, 0.0, // 18
		1.0, 0.0, 0.0, // 19
		1.0, 0.0, 0.0, // 20
		1.0, 0.0, 0.0, // 21
		1.0, 0.0, 0.0, // 22
		1.0, 0.0, 0.0, // 23

		// Top
		0.0, 1.0, 0.0, // 24
		0.0, 1.0, 0.0, // 25
		0.0, 1.0, 0.0, // 26
		0.0, 1.0, 0.0, // 27
		0.0, 1.0, 0.0, // 28
		0.0, 1.0, 0.0, // 29

		// Bottom
		0.0, -1.0, 0.0, // 30
		0.0, -1.0, 0.0, // 31
		0.0, -1.0, 0.0, // 32
		0.0, -1.0, 0.0, // 33
		0.0, -1.0, 0.0, // 34
		0.0, -1.0, 0.0, // 35
	}

	vertexBuffer := makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&vertexData[0]), sizeFloat*len(vertexData)) // 4 == sizeof float32
	normalBuffer := makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&normalData[0]), sizeFloat*len(normalData))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	program := readShaders()

	gl.UseProgram(program)

	pv := gl.GLString("pv")
	pvId := gl.GetUniformLocation(program, pv)
	gl.GLStringFree(pv)
	model := gl.GLString("model")
	modelId := gl.GetUniformLocation(program, model)
	gl.GLStringFree(model)

	heightStr := gl.GLString("height")
	heightVal := gl.GetUniformLocation(program, heightStr)
	gl.GLStringFree(heightStr)

	pvMatrix := matrix.MultiplyMatrix(projMatrix, viewMatrix)

	logicCh := make(chan float64)
	go logicLoop(logicCh)

	rot := 0.0
	gl.ClearColor(0.5, 0.5, 1.0, 1.0)
	currentTick := time.Now().UnixNano() / 1000000.0
	frameCount := 0
	for glfw.WindowParam(glfw.Opened) == 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		select {
		case rot = <-logicCh:
		default:
		}

		glPVMatrix := matrixToGL(pvMatrix)
		gl.UseProgram(program)
		gl.UniformMatrix4fv(pvId, 1, gl.FALSE, &glPVMatrix[0])

		for z := -10; z <= 10; z++ {
			for x := -10; x <= 10; x++ {
				delay := rot + float64(z*10) + float64(x*10)
				wave := math.Sin(delay * (math.Pi / 180.0))

				angle := wave * -15.0
				modelMatrix := matrix.NewIdentityMatrix()
				modelMatrix.Translate(float64(x), wave, float64(z))
				modelMatrix.RotateX(angle)

				glModelMatrix := matrixToGL(modelMatrix)

				gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])
				gl.Uniform1f(heightVal, gl.Float((wave+1.0)/2.0))

				gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
				gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
				gl.DrawArrays(gl.TRIANGLES, 0, gl.Sizei(len(vertexData)))
			}
		}

		if err := gl.GetError(); err != 0 {
			fmt.Printf("Err: %d\n", err)
			break
		}

		glfw.SwapBuffers()
		frameCount++

		newTick := time.Now().UnixNano() / 1000000.0
		if newTick-currentTick >= 1000.0 {
			fmt.Printf("FPS: %d\n", frameCount)
			frameCount = 0
			currentTick = newTick
		}
	}
}

func logicLoop(logicCh chan<- float64) {
	currentTick := time.Now().UnixNano() / 1000000.0

	rot := 0.0
	remainder := 0.0
	for {
		newTick := time.Now().UnixNano() / 1000000.0
		elapsedTick := float64(newTick-currentTick) + remainder
		if elapsedTick >= 16.0 {
			for elapsedTick >= 16.0 {
				elapsedTick -= 16.0

				rot += 0.5
				if rot >= 360 {
					rot = 0.0
				}
			}
			remainder = math.Max(elapsedTick, 0.0)
			currentTick = newTick
			logicCh <- rot
		}

		time.Sleep(1)
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
