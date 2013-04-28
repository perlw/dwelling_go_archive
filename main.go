package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
	"io/ioutil"
	"math"
	"runtime"
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
	glfw.SetWindowTitle("AH MAH GAHD IT WORKZ")

	if err := gl.Init(); err != nil {
		fmt.Printf("gl: %s\n", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.5, 0.5, 0.5, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Viewport(0, 0, 640, 480)

	vMatrix := viewMatrix()
	projMatrix := perspectiveMatrix(53.13, 640.0/480.0, 0.1, 1000.0)
	modelMatrix := translationMatrix(0.0, 0.0, -5.0)
	fmt.Println(projMatrix)
	fmt.Println(modelMatrix)

	var vao gl.Uint
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	//sizeFloat := int(unsafe.Sizeof([1]float32{}))
	vertexData := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	vertexBuffer := makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&vertexData[0]), 4*len(vertexData)) // 4 == sizeof float32

	program := readShaders()

	gl.UseProgram(program)

	view := gl.GLString("view")
	defer gl.GLStringFree(view)
	viewId := gl.GetUniformLocation(program, view)
	proj := gl.GLString("proj")
	defer gl.GLStringFree(proj)
	projId := gl.GetUniformLocation(program, proj)
	model := gl.GLString("model")
	defer gl.GLStringFree(model)
	modelId := gl.GetUniformLocation(program, model)

	projMatrix = flipMatrix(projMatrix)
	vMatrix = flipMatrix(vMatrix)
	modelMatrix = flipMatrix(modelMatrix)

	xpos := 0.0
	gl.ClearColor(0.5, 0.5, 1.0, 1.0)
	for glfw.WindowParam(glfw.Opened) == 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		xpos += 0.0005
		if xpos > math.Pi*2 {
			xpos = 0.0
		}
		modelMatrix = flipMatrix(translationMatrix(gl.Float(math.Sin(xpos)), 0.0, -5.0+gl.Float(math.Cos(xpos))))

		gl.UseProgram(program)
		gl.UniformMatrix4fv(projId, 1, gl.FALSE, &projMatrix[0])
		gl.UniformMatrix4fv(viewId, 1, gl.FALSE, &vMatrix[0])
		gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &modelMatrix[0])

		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

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

func identityMatrix() [16]gl.Float {
	return [16]gl.Float{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func viewMatrix() [16]gl.Float {
	matrix := identityMatrix()

	matrix[0] = gl.Float(math.Cos(-0.83))
	matrix[2] = gl.Float(-math.Sin(-0.83))
	matrix[8] = gl.Float(math.Sin(-0.83))
	matrix[10] = gl.Float(math.Cos(-0.83))

	matrix[3] = 4.0
	matrix[7] = 0.0
	matrix[11] = 0.0

	return matrix
}

func perspectiveMatrix(fov, ratio, nearp, farp float64) [16]gl.Float {
	f := 1.0 / math.Tan((fov/2.0)*(math.Pi/180.0))

	matrix := identityMatrix()
	matrix[0] = gl.Float(f / ratio)
	matrix[5] = gl.Float(f)
	matrix[10] = -gl.Float((farp + nearp) / (farp - nearp))
	matrix[11] = -gl.Float((2.0 * farp * nearp) / (farp - nearp))
	matrix[14] = -1.0
	matrix[15] = 0.0

	return matrix
}

func translationMatrix(x, y, z gl.Float) [16]gl.Float {
	matrix := identityMatrix()
	matrix[3] = x
	matrix[7] = y
	matrix[11] = z
	return matrix
}

func flipMatrix(matrix [16]gl.Float) [16]gl.Float {
	var newMatrix [16]gl.Float

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			newMatrix[(x*4)+y] = matrix[(y*4)+x]
		}
	}

	return newMatrix
}

/*func multMatrix(a, b gl.Float) [16]gl.Float {
	matrix := identityMatrix()

}*/
