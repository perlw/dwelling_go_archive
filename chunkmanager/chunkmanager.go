package chunkmanager

import (
	"dwelling/camera"
	"dwelling/math/matrix"
	"dwelling/math/vector"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const ChunkBase int = 16

type ChunkCoord struct {
	X, Y, Z int
}

type BlockCoord struct {
	X, Y, Z int
}

type Chunk struct {
	data         map[BlockCoord]*Block
	mesh         ChunkMesh
	IsLoaded     bool
	IsSetup      bool
	IsRebuilding bool
	MouseHit     bool
}

type Block struct {
	visible bool
}

var chunkMap = map[ChunkCoord]*Chunk{}
var rebuildChunks = map[ChunkCoord]*Chunk{}
var visibleChunks = map[ChunkCoord]*Chunk{}
var renderChunks = map[ChunkCoord]*Chunk{}

var camPos = vector.Vector3f{0.0, 0.0, 0.0}
var camView = vector.Vector3f{0.0, 0.0, -1.0}

var debugMode = false

func Start() {
	rand.Seed(time.Now().Unix())

	cubed := 4
	for x := 0; x < cubed; x++ {
		for z := 0; z < cubed; z++ {
			for y := 0; y < cubed; y++ {
				val := rand.Intn(5)

				var chunk *Chunk
				switch val {
				case 0:
					chunk = newCubeChunk()
				case 1:
					chunk = newPyramidChunk(false)
				case 2:
					chunk = newPyramidChunk(true)
				case 3:
					chunk = newSphereChunk()
				case 4:
					chunk = newWireCubeChunk()
				default:
					chunk = newCubeChunk()
				}
				chunkMap[ChunkCoord{x, y, z}] = chunk
			}
		}
	}
}

func SetDebug(mode bool) {
	debugMode = mode
}

func DebugDeleteRandomBlock() {
	deleted := false

	for !deleted {
		cx := rand.Intn(4)
		cy := rand.Intn(4)
		cz := rand.Intn(4)

		if chunk, ok := chunkMap[ChunkCoord{cx, cy, cz}]; ok {
			x := rand.Intn(16)
			y := rand.Intn(16)
			z := rand.Intn(16)

			if _, ok := chunk.data[BlockCoord{x, y, z}]; ok {
				fmt.Printf("Deleted block at [%d,%d,%d]:[%d,%d,%d]\n", cx, cy, cz, x, y, z)

				delete(chunk.data, BlockCoord{x, y, z})
				rebuildChunks[ChunkCoord{cx, cy, cz}] = chunk

				deleted = true
			}
		}
	}

}

func GetChunksAroundChunk(chunkPos ChunkCoord) [6]*Chunk {
	chunks := [6]*Chunk{nil, nil, nil, nil, nil, nil}

	if chnk, ok := chunkMap[ChunkCoord{chunkPos.X, chunkPos.Y, chunkPos.Z + 1}]; ok {
		chunks[FRONT] = chnk
	}
	if chnk, ok := chunkMap[ChunkCoord{chunkPos.X, chunkPos.Y, chunkPos.Z - 1}]; ok {
		chunks[BACK] = chnk
	}
	if chnk, ok := chunkMap[ChunkCoord{chunkPos.X - 1, chunkPos.Y, chunkPos.Z}]; ok {
		chunks[LEFT] = chnk
	}
	if chnk, ok := chunkMap[ChunkCoord{chunkPos.X + 1, chunkPos.Y, chunkPos.Z}]; ok {
		chunks[RIGHT] = chnk
	}
	if chnk, ok := chunkMap[ChunkCoord{chunkPos.X, chunkPos.Y + 1, chunkPos.Z}]; ok {
		chunks[TOP] = chnk
	}
	if chnk, ok := chunkMap[ChunkCoord{chunkPos.X, chunkPos.Y - 1, chunkPos.Z}]; ok {
		chunks[BOTTOM] = chnk
	}

	return chunks
}

// Traces a ray against a box.
// Returns:
// bool Intersected?
// float64 Intersect distance
// vector.Vector3f Intersection point
func CastRayAtBox(rayOrig, rayDir, boxPos vector.Vector3f, boxSize, stepSize, maxDepth float64) (bool, float64, vector.Vector3f) {
	planeNormals := [6]vector.Vector3f{
		{0.0, 0.0, 1.0},
		{0.0, 0.0, -1.0},
		{-1.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{0.0, -1.0, 0.0},
	}

	planePos := [6]vector.Vector3f{
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z + boxSize},
		{boxPos.X + boxSize, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y + boxSize, boxPos.Z},
	}

	for depth := 0.0; depth < maxDepth; depth += stepSize {
		inside := 0
		rayStep := vector.Vector3f{}
		for t := 0; t < 6; t++ {
			d := -(vector.DotProduct(planeNormals[t], planePos[t]))
			rayStep = rayOrig.Add(rayDir.MulScalar(depth))
			deep := vector.DotProduct(planeNormals[t], rayStep) + d

			if deep > 0.0 {
				inside++
			}
		}
		if inside >= 6 {
			return true, depth, rayStep
		}
	}

	return false, 0.0, vector.Vector3f{}
}

// Checks if point is inside box
// Returns:
// bool In box?
func PointInBox(point, boxPos vector.Vector3f, boxSize float64) bool {
	planeNormals := [6]vector.Vector3f{
		{0.0, 0.0, 1.0},
		{0.0, 0.0, -1.0},
		{-1.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{0.0, -1.0, 0.0},
	}

	planePos := [6]vector.Vector3f{
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z + boxSize},
		{boxPos.X + boxSize, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y + boxSize, boxPos.Z},
	}

	inside := 0
	for t := 0; t < 6; t++ {
		d := -(vector.DotProduct(planeNormals[t], planePos[t]))
		deep := vector.DotProduct(planeNormals[t], point) + d

		if deep > 0.0 {
			inside++
		}
	}
	if inside >= 6 {
		return true
	}

	return false
}

// TODO: Optmise chunk lookup by adding checked chunks to map copy and check against to avoid double and triple checking
// TODO: Rename func
func ClickedInChunk(mx, my int, cam *camera.Camera) {
	mouseNear, _ := matrix.Unproject(vector.Vector3f{float64(mx), float64(480 - my), 0.0}, cam.ViewMatrix, cam.ProjectionMatrix, 640, 480)
	mouseFar, _ := matrix.Unproject(vector.Vector3f{float64(mx), float64(480 - my), 1.0}, cam.ViewMatrix, cam.ProjectionMatrix, 640, 480)
	cam.MousePos = cam.Pos
	cam.MouseDir = mouseFar.Sub(mouseNear).Normalize()

	for _, chnk := range renderChunks {
		chnk.MouseHit = false
	}

	fmt.Println("---")
	chunkBase := float64(ChunkBase)
	for dist := 0.0; dist < 128.0; dist += 8.0 {
		rayStep := cam.MousePos.Add(cam.MouseDir.MulScalar(dist))

		pos := ChunkCoord{
			int(math.Trunc(rayStep.X / chunkBase)),
			int(math.Trunc(rayStep.Y / chunkBase)),
			int(math.Trunc(rayStep.Z / chunkBase)),
		}

		if chnk, ok := renderChunks[pos]; ok {
			fmt.Printf("Found potential chunk at %v...", pos)
			boxPos := vector.Vector3f{float64(pos.X * ChunkBase), float64(pos.Y * ChunkBase), float64(pos.Z * ChunkBase)}
			boxSize := chunkBase
			if PointInBox(rayStep, boxPos, boxSize) {
				fmt.Printf("hit!\n")
				chnk.MouseHit = true

				startDist := dist - 8.0
				for dist := startDist; dist < startDist+24.0; dist += 0.5 {
					rayStep := cam.MousePos.Add(cam.MouseDir.MulScalar(dist))
					blkPos := BlockCoord{
						int(math.Trunc(rayStep.X)) - (pos.X * ChunkBase),
						int(math.Trunc(rayStep.Y)) - (pos.Y * ChunkBase),
						int(math.Trunc(rayStep.Z)) - (pos.Z * ChunkBase),
					}
					// Note: Checking blk.visible seems to be reason for skipping to other side
					// at times.
					if blk, ok := chnk.data[blkPos]; ok {
						if blk.visible {
							fmt.Printf("\tFound potential block at %v...", blkPos)

							boxPos := vector.Vector3f{
								float64(pos.X*ChunkBase + blkPos.X),
								float64(pos.Y*ChunkBase + blkPos.Y),
								float64(pos.Z*ChunkBase + blkPos.Z),
							}
							if PointInBox(rayStep, boxPos, 1.0) {
								fmt.Println("hit!\n")

								if !chnk.IsRebuilding {
									delete(chnk.data, blkPos)
									rebuildChunks[pos] = chnk
								}

								return
							} else {
								fmt.Println("nope...\n")
							}
						}
					}
				}
			} else {
				fmt.Println("nope...\n")
			}
		}
	}
}

func Update(cam *camera.Camera) {
	updateSetupList()
	updateRebuildList()
	updateVisibilityList(cam)

	if camPos != cam.Pos || camView != cam.Rot {
		updateRenderList(cam)

		camPos = cam.Pos
		camView = cam.Rot
	}
}

func updateSetupList() {
	for pos, chnk := range chunkMap {
		if chnk.IsLoaded && !chnk.IsSetup {
			rebuildChunks[pos] = chnk
			chnk.IsSetup = true
		}
	}
}

func updateRebuildList() {
	for pos, chnk := range rebuildChunks {
		chnk.IsRebuilding = true
		chnk.UpdateChunkMesh(pos)
		chnk.IsRebuilding = false
	}

	rebuildChunks = map[ChunkCoord]*Chunk{}
}

func updateVisibilityList(cam *camera.Camera) {
	// TODO: Add chunk range limit
	for t, chnk := range chunkMap {
		if chnk.IsLoaded && chnk.IsSetup && !chnk.IsRebuilding {
			if _, ok := visibleChunks[t]; !ok {
				fmt.Printf("Added chunk at %v to visible list.\n", t)
				visibleChunks[t] = chunkMap[t]
			}
		} else {
			if _, ok := visibleChunks[t]; ok {
				delete(visibleChunks, t)
			}
		}
	}
}

func updateRenderList(cam *camera.Camera) {
	renderChunks = map[ChunkCoord]*Chunk{}

	for pos := range visibleChunks {
		posx := float64(pos.X * ChunkBase)
		posy := float64(pos.Y * ChunkBase)
		posz := float64(pos.Z * ChunkBase)
		if cam.CubeInView(vector.Vector3f{posx, posy, posz}, float64(ChunkBase)) != 2 {
			renderChunks[pos] = visibleChunks[pos]
		}
	}
}
