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
	"dwelling/chunk"
	"dwelling/matrix"
)

type Camera struct {
	X, Y, Z             float64
	Rx, Ry, Rz          float64
	CullX, CullY, CullZ float64
	Fx, Fy, Fz          float64
	Frx, Fry, Frz       float64

	ViewMatrix       *matrix.Matrix
	ProjectionMatrix *matrix.Matrix
	PVMatrix         *matrix.Matrix

	Planes [6]Plane
}

type Plane struct {
	A, B, C, D float64
}

func (cam *Camera) updateViewMatrix() {
	view := matrix.NewIdentityMatrix()
	view.RotateX(-cam.Rx)
	view.RotateY(-cam.Ry)
	view.RotateZ(-cam.Rz)
	view.Translate(-cam.X, -cam.Y, -cam.Z)

	cam.ViewMatrix = view
}

func (cam *Camera) updatePVMatrix() {
	cam.PVMatrix = matrix.MultiplyMatrix(cam.ProjectionMatrix, cam.ViewMatrix)
}

func (cam *Camera) updateFrustum() {
	cam.Planes = [6]Plane{}

	// Left
	cam.Planes[0].A = cam.PVMatrix.Values[12] + cam.PVMatrix.Values[0]
	cam.Planes[0].B = cam.PVMatrix.Values[13] + cam.PVMatrix.Values[1]
	cam.Planes[0].C = cam.PVMatrix.Values[14] + cam.PVMatrix.Values[2]
	cam.Planes[0].D = cam.PVMatrix.Values[15] + cam.PVMatrix.Values[3]

	// Right
	cam.Planes[1].A = cam.PVMatrix.Values[12] - cam.PVMatrix.Values[0]
	cam.Planes[1].B = cam.PVMatrix.Values[13] - cam.PVMatrix.Values[1]
	cam.Planes[1].C = cam.PVMatrix.Values[14] - cam.PVMatrix.Values[2]
	cam.Planes[1].D = cam.PVMatrix.Values[15] - cam.PVMatrix.Values[3]

	// Top
	cam.Planes[2].A = cam.PVMatrix.Values[12] - cam.PVMatrix.Values[4]
	cam.Planes[2].B = cam.PVMatrix.Values[13] - cam.PVMatrix.Values[5]
	cam.Planes[2].C = cam.PVMatrix.Values[14] - cam.PVMatrix.Values[6]
	cam.Planes[2].D = cam.PVMatrix.Values[15] - cam.PVMatrix.Values[7]

	// Bottom
	cam.Planes[3].A = cam.PVMatrix.Values[12] + cam.PVMatrix.Values[4]
	cam.Planes[3].B = cam.PVMatrix.Values[13] + cam.PVMatrix.Values[5]
	cam.Planes[3].C = cam.PVMatrix.Values[14] + cam.PVMatrix.Values[6]
	cam.Planes[3].D = cam.PVMatrix.Values[15] + cam.PVMatrix.Values[7]

	// Near
	cam.Planes[4].A = cam.PVMatrix.Values[12] + cam.PVMatrix.Values[8]
	cam.Planes[4].B = cam.PVMatrix.Values[13] + cam.PVMatrix.Values[9]
	cam.Planes[4].C = cam.PVMatrix.Values[14] + cam.PVMatrix.Values[10]
	cam.Planes[4].D = cam.PVMatrix.Values[15] + cam.PVMatrix.Values[11]

	// Far
	cam.Planes[5].A = cam.PVMatrix.Values[12] - cam.PVMatrix.Values[8]
	cam.Planes[5].B = cam.PVMatrix.Values[13] - cam.PVMatrix.Values[9]
	cam.Planes[5].C = cam.PVMatrix.Values[14] - cam.PVMatrix.Values[10]
	cam.Planes[5].D = cam.PVMatrix.Values[15] - cam.PVMatrix.Values[11]

	/*for t := range cam.Planes {
		cam.Planes[t].Normalize()
	}*/

	cam.Fx = cam.X
	cam.Fy = cam.Y
	cam.Fz = cam.Z
	cam.Frx = cam.Rx
	cam.Fry = cam.Ry
	cam.Frz = cam.Rz
}

func (cam *Camera) CubeInView(origo [3]float64, size float64) int {
	corners := [8][3]float64{
		{origo[0], origo[1], origo[2]},
		{origo[0] + size, origo[1], origo[2]},
		{origo[0] + size, origo[1], origo[2] + size},
		{origo[0], origo[1], origo[2] + size},
		{origo[0], origo[1] + size, origo[2] + size},
		{origo[0] + size, origo[1] + size, origo[2] + size},
		{origo[0] + size, origo[1] + size, origo[2]},
		{origo[0], origo[1] + size, origo[2]},
	}

	status := 0 // 0 inside, 1 partly, 2 outside
	for t := range cam.Planes {
		in, out := 0, 0

		for u := range corners {
			if cam.Planes[t].ClassifyPoint(corners[u]) < 0.0 {
				out++
			} else {
				in++
			}
		}

		if in == 0 {
			return 2
		} else if out > 0 {
			status = 1
		}
	}

	return status
}

func (plane *Plane) Normalize() {
	magnitude := math.Sqrt((plane.A * plane.A) + (plane.B * plane.B) + (plane.C * plane.C) + (plane.D * plane.D))

	plane.A /= magnitude
	plane.B /= magnitude
	plane.C /= magnitude
	plane.D /= magnitude
}

func (plane *Plane) ClassifyPoint(vector [3]float64) float64 {
	return (plane.A * vector[0]) + (plane.B * vector[1]) + (plane.C * vector[2]) + plane.D
}

var cam = Camera{X: 8.0, Y: 8.0, Z: 16.0, Rx: 0.0, Ry: 0.0, Rz: 0.0}

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

	cam.ProjectionMatrix = matrix.NewPerspectiveMatrix(53.13, 640.0/480.0, 0.1, 1000.0)
	cam.updateViewMatrix()
	cam.updatePVMatrix()
	cam.updateFrustum()

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

	frustumBuffer := createFrustumMesh(&cam)

	camCh := make(chan bool)
	go logicLoop(camCh, &cam)

	gl.ClearColor(0.5, 0.5, 1.0, 1.0)
	currentTick := time.Now().UnixNano() / 1000000.0
	frameCount := 0
	for glfw.WindowParam(glfw.Opened) == 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		select {
		case <-camCh:
			cam.updateViewMatrix()
			cam.updatePVMatrix()
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
		renderFrustumMesh(&cam, frustumBuffer)

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

func createFrustumMesh(cam *Camera) gl.Uint {
	var buffer gl.Uint
	sizeFloat := int(unsafe.Sizeof([1]float32{}))

	fmt.Println(cam.ProjectionMatrix)
	proj := cam.ProjectionMatrix.Values
	near := proj[11] / (proj[10] - 1.0)
	far := 100.0 //proj[11] / (1.0 + proj[10])
	nLeft := float32(near * (proj[2] - 1.0) / proj[0])
	nRight := float32(near * (1.0 + proj[2]) / proj[0])
	nTop := float32(near * (1.0 + proj[6]) / proj[5])
	nBottom := float32(near * (proj[6] - 1.0) / proj[5])
	fLeft := float32(far * (proj[2] - 1.0) / proj[0])
	fRight := float32(far * (1.0 + proj[2]) / proj[0])
	fTop := float32(far * (1.0 + proj[6]) / proj[5])
	fBottom := float32(far * (proj[6] - 1.0) / proj[5])

	vertices := [...]float32{
		0.0, 0.0, 0.0,
		fLeft, fBottom, float32(-far),

		0.0, 0.0, 0.0,
		fRight, fBottom, float32(-far),

		0.0, 0.0, 0.0,
		fRight, fTop, float32(-far),

		0.0, 0.0, 0.0,
		fLeft, fTop, float32(-far),

		fLeft, fBottom, float32(-far),
		fRight, fBottom, float32(-far),

		fRight, fTop, float32(-far),
		fLeft, fTop, float32(-far),

		fRight, fTop, float32(-far),
		fRight, fBottom, float32(-far),

		fLeft, fTop, float32(-far),
		fLeft, fBottom, float32(-far),

		nLeft, nBottom, float32(-near),
		nRight, nBottom, float32(-near),

		nRight, nTop, float32(-near),
		nLeft, nTop, float32(-near),

		nLeft, nTop, float32(-near),
		nLeft, nBottom, float32(-near),

		nRight, nTop, float32(-near),
		nRight, nBottom, float32(-near),
	}

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*len(vertices)), gl.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return buffer
}

func renderFrustumMesh(cam *Camera, meshBuffer gl.Uint) {
	gl.BindBuffer(gl.ARRAY_BUFFER, meshBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.DrawArrays(gl.LINES, 0, 24)
}

func logicLoop(camCh chan<- bool, cam *Camera) {
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
					cam.updateFrustum()
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
