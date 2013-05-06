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

type Block struct {
}

type WorldMesh struct {
	frontFaceBuffer  gl.Uint
	numFrontFaces    gl.Sizei
	backFaceBuffer   gl.Uint
	numBackFaces     gl.Sizei
	leftFaceBuffer   gl.Uint
	numLeftFaces     gl.Sizei
	rightFaceBuffer  gl.Uint
	numRightFaces    gl.Sizei
	topFaceBuffer    gl.Uint
	numTopFaces      gl.Sizei
	bottomFaceBuffer gl.Uint
	numBottomFaces   gl.Sizei
}

var pyramidBase int = 12

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
	viewMatrix.Translate(-float64(pyramidBase*8), -float64(pyramidBase*2), -float64(pyramidBase*8))
	projMatrix := matrix.NewPerspectiveMatrix(53.13, 640.0/480.0, 0.1, 1000.0)

	var vao gl.Uint
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	worldData := buildWorldMesh()
	fmt.Println(worldData)

	program := readShaders()

	gl.UseProgram(program)

	pv := gl.GLString("pv")
	pvId := gl.GetUniformLocation(program, pv)
	gl.GLStringFree(pv)
	model := gl.GLString("model")
	modelId := gl.GetUniformLocation(program, model)
	gl.GLStringFree(model)
	normal := gl.GLString("normal")
	normalId := gl.GetUniformLocation(program, normal)
	gl.GLStringFree(normal)
	pyramid := gl.GLString("pyramid")
	pyramidId := gl.GetUniformLocation(program, pyramid)
	gl.GLStringFree(pyramid)
	waveGL := gl.GLString("wave")
	waveId := gl.GetUniformLocation(program, waveGL)
	gl.GLStringFree(waveGL)

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
		gl.Uniform1i(pyramidId, gl.Int(pyramidBase))

		for z := -5; z <= 5; z++ {
			for x := -5; x <= 5; x++ {
				delay := rot + float64(z*20)
				wave := math.Sin(delay * (math.Pi / 180.0))
				angle := wave * -15.0

				gl.Uniform1f(waveId, gl.Float((wave+1.0)/2.0))
				modelMatrix := matrix.NewIdentityMatrix()
				modelMatrix.Translate(float64(x*pyramidBase*2-pyramidBase), wave*float64(pyramidBase/4), float64(z*pyramidBase*2-pyramidBase))
				modelMatrix.RotateY(angle)
				glModelMatrix := matrixToGL(modelMatrix)
				gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])
				drawWorld(worldData, normalId)
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

func buildWorldMesh() WorldMesh {
	worldData := WorldMesh{}
	sizeFloat := int(unsafe.Sizeof([1]float32{}))

	chunk := map[[3]int]Block{}

	for y := 0; y < pyramidBase; y++ {
		modBase := pyramidBase - y
		for x := -modBase; x < modBase; x++ {
			for z := -modBase; z < modBase; z++ {
				index := [3]int{x, y, z}
				chunk[index] = Block{}
			}
		}
	}

	frontFaces := []float32{}
	backFaces := []float32{}
	leftFaces := []float32{}
	rightFaces := []float32{}
	topFaces := []float32{}
	bottomFaces := []float32{}
	for pos, _ := range chunk {
		x := float32(pos[0])
		y := float32(pos[1])
		z := float32(pos[2])

		if _, ok := chunk[[3]int{pos[0], pos[1], pos[2] + 1}]; !ok {
			frontFaces = append(frontFaces,
				x, y, z,
				x+1.0, y, z,
				x+1.0, y+1.0, z,
				x+1.0, y+1.0, z,
				x, y+1.0, z,
				x, y, z,
			)
		}

		if _, ok := chunk[[3]int{pos[0], pos[1], pos[2] - 1}]; !ok {
			backFaces = append(backFaces,
				x+1.0, y+1.0, z-1.0,
				x+1.0, y, z-1.0,
				x, y, z-1.0,
				x, y, z-1.0,
				x, y+1.0, z-1.0,
				x+1.0, y+1.0, z-1.0,
			)
		}

		if _, ok := chunk[[3]int{pos[0] - 1, pos[1], pos[2]}]; !ok {
			leftFaces = append(leftFaces,
				x, y, z-1.0,
				x, y, z,
				x, y+1.0, z,
				x, y+1.0, z,
				x, y+1.0, z-1.0,
				x, y, z-1.0,
			)
		}

		if _, ok := chunk[[3]int{pos[0] + 1, pos[1], pos[2]}]; !ok {
			rightFaces = append(rightFaces,
				x+1.0, y+1.0, z,
				x+1.0, y, z,
				x+1.0, y, z-1.0,
				x+1.0, y, z-1.0,
				x+1.0, y+1.0, z-1.0,
				x+1.0, y+1.0, z,
			)
		}

		if _, ok := chunk[[3]int{pos[0], pos[1] + 1, pos[2]}]; !ok {
			topFaces = append(topFaces,
				x+1.0, y+1.0, z,
				x+1.0, y+1.0, z-1.0,
				x, y+1.0, z-1.0,
				x, y+1.0, z-1.0,
				x, y+1.0, z,
				x+1.0, y+1.0, z,
			)
		}

		if _, ok := chunk[[3]int{pos[0], pos[1] - 1, pos[2]}]; !ok {
			bottomFaces = append(bottomFaces,
				x, y, z-1.0,
				x+1.0, y, z-1.0,
				x+1.0, y, z,
				x+1.0, y, z,
				x, y, z,
				x, y, z-1.0,
			)
		}
	}

	worldData.numFrontFaces = gl.Sizei(len(frontFaces))
	if worldData.numFrontFaces > 0 {
		worldData.frontFaceBuffer = makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&frontFaces[0]), sizeFloat*len(frontFaces))
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.frontFaceBuffer)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	worldData.numBackFaces = gl.Sizei(len(backFaces))
	if worldData.numBackFaces > 0 {
		worldData.backFaceBuffer = makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&backFaces[0]), sizeFloat*len(backFaces))
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.backFaceBuffer)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	worldData.numLeftFaces = gl.Sizei(len(leftFaces))
	if worldData.numLeftFaces > 0 {
		worldData.leftFaceBuffer = makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&leftFaces[0]), sizeFloat*len(leftFaces))
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.leftFaceBuffer)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	worldData.numRightFaces = gl.Sizei(len(rightFaces))
	if worldData.numRightFaces > 0 {
		worldData.rightFaceBuffer = makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&rightFaces[0]), sizeFloat*len(rightFaces))
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.rightFaceBuffer)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	worldData.numTopFaces = gl.Sizei(len(topFaces))
	if worldData.numTopFaces > 0 {
		worldData.topFaceBuffer = makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&topFaces[0]), sizeFloat*len(topFaces))
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.topFaceBuffer)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	worldData.numBottomFaces = gl.Sizei(len(bottomFaces))
	if worldData.numBottomFaces > 0 {
		worldData.bottomFaceBuffer = makeBuffer(gl.ARRAY_BUFFER, gl.Pointer(&bottomFaces[0]), sizeFloat*len(bottomFaces))
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.bottomFaceBuffer)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	numFaces := 0
	numFaces += int(worldData.numFrontFaces / 9.0)
	numFaces += int(worldData.numBackFaces / 9.0)
	numFaces += int(worldData.numLeftFaces / 9.0)
	numFaces += int(worldData.numRightFaces / 9.0)
	numFaces += int(worldData.numTopFaces / 9.0)
	numFaces += int(worldData.numBottomFaces / 9.0)

	worstCaseFaces := len(chunk) * 12
	fmt.Printf("%d faces vs %d total, saved %d\n", numFaces, worstCaseFaces, worstCaseFaces-numFaces)

	return worldData
}

func drawWorld(worldData WorldMesh, normalId gl.Int) {
	frontNormal := [3]gl.Float{0.0, 0.0, 1.0}
	backNormal := [3]gl.Float{0.0, 0.0, -1.0}
	leftNormal := [3]gl.Float{-1.0, 0.0, 0.0}
	rightNormal := [3]gl.Float{1.0, 0.0, 0.0}
	topNormal := [3]gl.Float{0.0, 1.0, 0.0}
	bottomNormal := [3]gl.Float{0.0, -1.0, 0.0}

	if worldData.numFrontFaces > 0 {
		gl.Uniform3fv(normalId, 1, &frontNormal[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.frontFaceBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, worldData.numFrontFaces)
	}

	if worldData.numBackFaces > 0 {
		gl.Uniform3fv(normalId, 1, &backNormal[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.backFaceBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, worldData.numBackFaces)
	}

	if worldData.numLeftFaces > 0 {
		gl.Uniform3fv(normalId, 1, &leftNormal[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.leftFaceBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, worldData.numLeftFaces)
	}

	if worldData.numRightFaces > 0 {
		gl.Uniform3fv(normalId, 1, &rightNormal[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.rightFaceBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, worldData.numRightFaces)
	}

	if worldData.numTopFaces > 0 {
		gl.Uniform3fv(normalId, 1, &topNormal[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.topFaceBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, worldData.numTopFaces)
	}

	if worldData.numBottomFaces > 0 {
		gl.Uniform3fv(normalId, 1, &bottomNormal[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, worldData.bottomFaceBuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, worldData.numBottomFaces)
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
