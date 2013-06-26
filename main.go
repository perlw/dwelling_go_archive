package main

import (
	"bedrock"
	"dwelling/camera"
	"dwelling/chunkmanager"
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

	if err := bedrock.Init(); err != nil {
		fmt.Println(err)
		return
	}

	if err := cam.Init(); err != nil {
		fmt.Println(err)
		return
	}

	if err := chunkmanager.Start(); err != nil {
		fmt.Println(err)
		return
	}

	camCh := make(chan bool)
	debugCh := make(chan bool)
	logicCh := make(chan bool)
	exitCh := make(chan bool)
	go logicLoop(camCh, debugCh, logicCh, exitCh, &cam)

	gl.ClearColor(0.8, 0.8, 0.8, 1.0)
	currentTick := time.Now().UnixNano() / 1e6
	frameCount := 0
	debugMode := false
	running := true
	for glfw.WindowParam(glfw.Opened) == 1 && running {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		select {
		case <-camCh:
			cam.UpdateViewMatrix()
			cam.UpdatePVMatrix()
			if !debugMode {
				cam.UpdateFrustum()
				cam.CullPos.X = cam.Pos.X
				cam.CullPos.Y = cam.Pos.Y
				cam.CullPos.Z = cam.Pos.Z
			}
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

		chunkmanager.Render(&cam)

		if debugMode {
			cam.RenderDebugMeshes()
		}

		if err := gl.GetError(); err != 0 {
			fmt.Printf("Err: %d\n", err)
			break
		}

		glfw.SwapBuffers()
		frameCount++

		newTick := time.Now().UnixNano() / 1e6
		if newTick-currentTick >= 1000.0 {
			fmt.Printf("FPS: %d\n", frameCount)
			frameCount = 0
			currentTick = newTick
		}

		runtime.Gosched()
	}

	bedrock.Cleanup()
}

func logicLoop(camCh chan<- bool, debugCh chan<- bool, logicCh chan<- bool, exitCh chan<- bool, cam *camera.Camera) {
	currentTick := time.Now().UnixNano() / 1e6

	rotSpeed := 1.0
	camSpeed := 0.25

	keyF1Held := false
	debugMode := false

	remainder := 0.0
	for {
		newTick := time.Now().UnixNano() / 1e6
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
					// Dangerous, race condition!
					mx, my := glfw.MousePos()
					chunkmanager.ClickedInChunk(mx, my, cam)
				}

				if debugMode {
					// Dangerous, race condition!
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
				camCh <- true
			}
			logicCh <- true
		}
	}
}
