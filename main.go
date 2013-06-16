package main

import (
	"dwelling/camera"
	"dwelling/chunkmanager"
	"dwelling/shader"
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
	"math"
	"runtime"
	"time"
)

var cam = camera.Camera{}

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

	if err := cam.Init(); err != nil {
		fmt.Println(err)
		return
	}

	chunkmanager.Start()

	simpleShader, err := shader.LoadShaderProgram("simple")
	if err != nil {
		fmt.Println(err)
		return
	}

	camCh := make(chan bool)
	debugCh := make(chan bool)
	logicCh := make(chan bool)
	exitCh := make(chan bool)
	go logicLoop(camCh, debugCh, logicCh, exitCh, &cam)

	gl.ClearColor(0.25, 0.25, 0.25, 1.0)
	currentTick := time.Now().UnixNano() / 1000000.0
	frameCount := 0
	debugMode := false
	running := true
	for glfw.WindowParam(glfw.Opened) == 1 && running {
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
		case <-exitCh:
			running = false
		default:
		}

		simpleShader.Use()
		simpleShader.SetUniformMatrix("pv", cam.PVMatrix)
		chunkmanager.Render(simpleShader.GetProgramId(), &cam)

		if debugMode {
			cam.RenderDebugMeshes()
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

func logicLoop(camCh chan<- bool, debugCh chan<- bool, logicCh chan<- bool, exitCh chan<- bool, cam *camera.Camera) {
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
				if glfw.Key(glfw.KeyEsc) == glfw.KeyPress {
					exitCh <- true
				}
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

				if glfw.MouseButton(glfw.MouseLeft) == glfw.KeyPress {
					mx, my := glfw.MousePos()
					chunkmanager.ClickedInChunk(mx, my, cam)
				}

				if debugMode {
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
