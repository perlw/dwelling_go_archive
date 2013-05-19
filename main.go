package main

import (
	"dwelling/camera"
	"dwelling/chunkmanager"
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

var cam = camera.Camera{Pos: vector.Vector3f{-48.0, 32.0, -48.0}, Rot: vector.Vector3f{0.0, 135, 0.0}}

func main() {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Using %d cpus for concurrency\n", runtime.NumCPU())

	if err := glfw.Init(); err != nil {
		fmt.Printf("glfw: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 0)
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
	cam.FrustumPos = cam.Pos
	cam.FrustumRot = cam.Rot
	cam.CullPos = cam.Pos
	cam.UpdateViewMatrix()
	cam.UpdatePVMatrix()
	cam.UpdateFrustum()

	var vao gl.Uint
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	chunkmanager.Start()

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

	frustumBuffer := camera.CreateFrustumMesh(&cam)

	camCh := make(chan bool)
	debugCh := make(chan bool)
	logicCh := make(chan bool)
	delCh := make(chan bool)
	go logicLoop(camCh, debugCh, logicCh, delCh, &cam)

	gl.ClearColor(0.25, 0.25, 0.25, 1.0)
	currentTick := time.Now().UnixNano() / 1000000.0
	frameCount := 0
	debugMode := false
	for glfw.WindowParam(glfw.Opened) == 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		select {
		case <-camCh:
			cam.UpdateViewMatrix()
			cam.UpdatePVMatrix()
		case debugMode = <-debugCh:
			if debugMode {
				chunkmanager.SetDebug(true)
			} else {
				chunkmanager.SetDebug(false)
			}
		case <-logicCh:
			chunkmanager.Update(&cam)
		case <-delCh:
			chunkmanager.DebugDeleteRandomBlock()
		default:
		}

		glPVMatrix := cam.PVMatrix.ToGL()
		gl.UseProgram(program)
		gl.UniformMatrix4fv(pvId, 1, gl.FALSE, &glPVMatrix[0])

		chunkmanager.Render(program, &cam)

		if debugMode {
			// Render frustum
			glNormal := vector.Vector3f{0.0, 1.0, 0.0}.ToGL()
			gl.Uniform3fv(normalId, 1, &glNormal[0])
			modelMatrix := matrix.NewIdentityMatrix()
			modelMatrix.TranslateVector(cam.FrustumPos)
			modelMatrix.RotateY(cam.FrustumRot.Y)
			modelMatrix.RotateX(cam.FrustumRot.X)
			glModelMatrix := modelMatrix.ToGL()
			gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])
			camera.RenderFrustumMesh(&cam, frustumBuffer)
			// Render frustum
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

		runtime.Gosched()
	}
}

func logicLoop(camCh chan<- bool, debugCh chan<- bool, logicCh chan<- bool, delCh chan<- bool, cam *camera.Camera) {
	currentTick := time.Now().UnixNano() / 1000000.0

	rotSpeed := 1.0
	camSpeed := 0.25

	keyF1Held := false
	debugMode := false

	remainder := 0.0
	for {
		newTick := time.Now().UnixNano() / 1000000.0
		elapsedTick := float64(newTick-currentTick) + remainder
		if elapsedTick >= 16.0 {
			update := false
			// Catch up loop
			for elapsedTick >= 16.0 {
				elapsedTick -= 16.0

				// Execute logic
				if !keyF1Held && glfw.Key(glfw.KeyF1) == glfw.KeyPress {
					keyF1Held = true
				}
				if keyF1Held && glfw.Key(glfw.KeyF1) == glfw.KeyRelease {
					keyF1Held = false

					debugMode = !debugMode
					debugCh <- debugMode
					fmt.Printf("Debug mode: %v.\n", debugMode)
				}

				if glfw.Key(glfw.KeyUp) == glfw.KeyPress {
					cam.Rot.X = math.Max(cam.Rot.X-rotSpeed, -90.0)
					update = true
				}
				if glfw.Key(glfw.KeyDown) == glfw.KeyPress {
					cam.Rot.X = math.Min(cam.Rot.X+rotSpeed, 90.0)
					update = true
				}
				if glfw.Key(glfw.KeyLeft) == glfw.KeyPress {
					cam.Rot.Y -= rotSpeed
					update = true
				}
				if glfw.Key(glfw.KeyRight) == glfw.KeyPress {
					cam.Rot.Y += rotSpeed
					update = true
				}

				if glfw.Key('W') == glfw.KeyPress {
					xRadii := -cam.Rot.X * (math.Pi / 180.0)
					yRadii := -cam.Rot.Y * (math.Pi / 180.0)
					xMove := math.Sin(yRadii) * camSpeed
					yMove := math.Sin(xRadii) * camSpeed
					zMove := math.Cos(yRadii) * camSpeed
					cam.Pos.X -= xMove
					cam.Pos.Y += yMove
					cam.Pos.Z -= zMove
					update = true
				}
				if glfw.Key('S') == glfw.KeyPress {
					xRadii := -cam.Rot.X * (math.Pi / 180.0)
					yRadii := -cam.Rot.Y * (math.Pi / 180.0)
					xMove := math.Sin(yRadii) * camSpeed
					yMove := math.Sin(xRadii) * camSpeed
					zMove := math.Cos(yRadii) * camSpeed
					cam.Pos.X += xMove
					cam.Pos.Y -= yMove
					cam.Pos.Z += zMove
					update = true
				}
				if glfw.Key('A') == glfw.KeyPress {
					yRadii := -(cam.Rot.Y - 90.0) * (math.Pi / 180.0)
					xMove := math.Sin(yRadii) * camSpeed
					zMove := math.Cos(yRadii) * camSpeed
					cam.Pos.X -= xMove
					cam.Pos.Z -= zMove
					update = true
				}
				if glfw.Key('D') == glfw.KeyPress {
					yRadii := -(cam.Rot.Y + 90.0) * (math.Pi / 180.0)
					xMove := math.Sin(yRadii) * camSpeed
					zMove := math.Cos(yRadii) * camSpeed
					cam.Pos.X -= xMove
					cam.Pos.Z -= zMove
					update = true
				}

				if debugMode {
					if glfw.Key(glfw.KeyDel) == glfw.KeyPress {
						delCh <- true
					}

					if glfw.Key('F') == glfw.KeyPress {
						cam.UpdateFrustum()
					}
					if glfw.Key('C') == glfw.KeyPress {
						cam.CullPos.X = cam.Pos.X
						cam.CullPos.Y = cam.Pos.Y
						cam.CullPos.Z = cam.Pos.Z
					}
				}
			}
			remainder = math.Max(elapsedTick, 0.0)
			currentTick = newTick

			if update {
				if !debugMode {
					cam.UpdateFrustum()
					cam.CullPos.X = cam.Pos.X
					cam.CullPos.Y = cam.Pos.Y
					cam.CullPos.Z = cam.Pos.Z
				}
				camCh <- true
			}
			logicCh <- true
		}
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
	fmt.Println("vertex")
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
	fmt.Println("fragment")
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
