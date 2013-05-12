package main

import (
	"dwelling/camera"
	"dwelling/chunk"
	"dwelling/math/matrix"
	"dwelling/math/vector"
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
	"io/ioutil"
	"math"
	"runtime"
	"time"
)

var cam = camera.Camera{X: 8.0, Y: 8.0, Z: 16.0, Rx: 0.0, Ry: 0.0, Rz: 0.0}

func main() {
	runtime.LockOSThread()

	tmp, tmp2, tmp3 := vector.Vector2i{1, 2}, vector.Vector2i{1, 2}, vector.Vector3f{1, 2, 3}
	fmt.Println(tmp, tmp2, tmp3)
	fmt.Println(tmp.Add(tmp2), tmp.Mul(tmp2))
	fmt.Println(tmp3.Normalize(), tmp3.Normalize().Length())

	return

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

	cam.ProjectionMatrix = matrix.NewPerspectiveMatrix(53.13, 640.0/480.0, 0.1, 1000.0)
	cam.UpdateViewMatrix()
	cam.UpdatePVMatrix()
	cam.UpdateFrustum()

	var vao gl.Uint
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	cubed := 4
	chunks := map[chunk.ChunkCoord]*chunk.Chunk{}
	for x := 0; x < cubed; x++ {
		for z := 0; z < cubed; z++ {
			for y := 0; y < cubed; y++ {
				chunks[chunk.ChunkCoord{x, y, z}] = chunk.NewCubeChunk()
			}
		}
	}

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
	maxHeight := gl.GLString("maxHeight")
	maxHeightId := gl.GetUniformLocation(program, maxHeight)
	gl.GLStringFree(maxHeight)
	chunkHeight := gl.GLString("chunkHeight")
	chunkHeightId := gl.GetUniformLocation(program, chunkHeight)
	gl.GLStringFree(chunkHeight)

	frustumBuffer := camera.CreateFrustumMesh(&cam)

	camCh := make(chan bool)
	go logicLoop(camCh, &cam)

	gl.ClearColor(0.5, 0.5, 1.0, 1.0)
	currentTick := time.Now().UnixNano() / 1000000.0
	frameCount := 0
	for glfw.WindowParam(glfw.Opened) == 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		select {
		case <-camCh:
			cam.UpdateViewMatrix()
			cam.UpdatePVMatrix()
		default:
		}

		glPVMatrix := matrixToGL(cam.PVMatrix)
		gl.UseProgram(program)
		gl.UniformMatrix4fv(pvId, 1, gl.FALSE, &glPVMatrix[0])

		for pos, chnk := range chunks {
			posx := float64(pos.X * chunk.CHUNK_BASE)
			posy := float64(pos.Y * chunk.CHUNK_BASE)
			posz := float64(pos.Z * chunk.CHUNK_BASE)
			if cam.CubeInView([3]float64{posx, posy, posz}, float64(chunk.CHUNK_BASE)) != 2 {
				modelMatrix := matrix.NewIdentityMatrix()
				modelMatrix.Translate(posx, posy, posz)
				glModelMatrix := matrixToGL(modelMatrix)
				gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])
				gl.Uniform1f(maxHeightId, gl.Float(chunk.CHUNK_BASE*cubed))
				gl.Uniform1f(chunkHeightId, gl.Float(posy))

				chnk.RenderChunk(normalId, [3]float64{cam.CullX, cam.CullY, cam.CullZ}, [3]float64{posx, 0.0, posz}, modelMatrix)
			}
		}

		modelMatrix := matrix.NewIdentityMatrix()
		modelMatrix.Translate(cam.Fx, cam.Fy, cam.Fz)
		modelMatrix.RotateY(cam.Fry)
		modelMatrix.RotateX(cam.Frx)
		glModelMatrix := matrixToGL(modelMatrix)
		gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])
		camera.RenderFrustumMesh(&cam, frustumBuffer)

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

func logicLoop(camCh chan<- bool, cam *camera.Camera) {
	currentTick := time.Now().UnixNano() / 1000000.0

	rotSpeed := 1.0
	camSpeed := 0.25

	remainder := 0.0
	for {
		newTick := time.Now().UnixNano() / 1000000.0
		elapsedTick := float64(newTick-currentTick) + remainder
		if elapsedTick >= 16.0 {
			update := false
			for elapsedTick >= 16.0 {
				elapsedTick -= 16.0

				// Execute logic
				if glfw.Key(glfw.KeyUp) == glfw.KeyPress {
					cam.Rx = math.Max(cam.Rx-rotSpeed, -90.0)
					update = true
				}
				if glfw.Key(glfw.KeyDown) == glfw.KeyPress {
					cam.Rx = math.Min(cam.Rx+rotSpeed, 90.0)
					update = true
				}
				if glfw.Key(glfw.KeyLeft) == glfw.KeyPress {
					cam.Ry -= rotSpeed
					update = true
				}
				if glfw.Key(glfw.KeyRight) == glfw.KeyPress {
					cam.Ry += rotSpeed
					update = true
				}

				if glfw.Key('W') == glfw.KeyPress {
					xRadii := -cam.Rx * (math.Pi / 180.0)
					yRadii := -cam.Ry * (math.Pi / 180.0)
					xMove := math.Sin(yRadii) * camSpeed
					yMove := math.Sin(xRadii) * camSpeed
					zMove := math.Cos(yRadii) * camSpeed
					cam.X -= xMove
					cam.Y += yMove
					cam.Z -= zMove
					update = true
				}
				if glfw.Key('S') == glfw.KeyPress {
					xRadii := -cam.Rx * (math.Pi / 180.0)
					yRadii := -cam.Ry * (math.Pi / 180.0)
					xMove := math.Sin(yRadii) * camSpeed
					yMove := math.Sin(xRadii) * camSpeed
					zMove := math.Cos(yRadii) * camSpeed
					cam.X += xMove
					cam.Y -= yMove
					cam.Z += zMove
					update = true
				}
				if glfw.Key('A') == glfw.KeyPress {
					yRadii := -(cam.Ry - 90.0) * (math.Pi / 180.0)
					xMove := math.Sin(yRadii) * camSpeed
					zMove := math.Cos(yRadii) * camSpeed
					cam.X -= xMove
					cam.Z -= zMove
					update = true
				}
				if glfw.Key('D') == glfw.KeyPress {
					yRadii := -(cam.Ry + 90.0) * (math.Pi / 180.0)
					xMove := math.Sin(yRadii) * camSpeed
					zMove := math.Cos(yRadii) * camSpeed
					cam.X -= xMove
					cam.Z -= zMove
					update = true
				}

				if glfw.Key('F') == glfw.KeyPress {
					cam.UpdateFrustum()
				}
				if glfw.Key('C') == glfw.KeyPress {
					cam.CullX = cam.X
					cam.CullY = cam.Y
					cam.CullZ = cam.Z
				}
			}
			remainder = math.Max(elapsedTick, 0.0)
			currentTick = newTick

			if update {
				camCh <- true
			}
		}

		time.Sleep(1)
	}
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

func multiplyMatrixVector(matrix *matrix.Matrix, vector [4]float64) [4]float64 {
	values := [4]float64{0.0, 0.0, 0.0}

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			i := (y * 4) + x
			values[y] += matrix.Values[i] * vector[x]
		}
	}

	return values
}
