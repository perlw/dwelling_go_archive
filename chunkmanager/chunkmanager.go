package chunkmanager

import (
	"dwelling/camera"
	"dwelling/math/matrix"
	"dwelling/math/vector"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const CHUNK_BASE int = 16

type ChunkCoord struct {
	X, Y, Z int
}

type BlockCoord struct {
	X, Y, Z int
}

type Chunk struct {
	data         map[BlockCoord]Block
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

// Note: This should really be rebuilt to take steps first and check against cube map instead
// of other way around.
func ClickedInChunk(mx, my int, cam *camera.Camera) {
	mouseNear, _ := matrix.Unproject(vector.Vector3f{float64(mx), float64(480 - my), 0.0}, cam.ViewMatrix, cam.ProjectionMatrix, 640, 480)
	mouseFar, _ := matrix.Unproject(vector.Vector3f{float64(mx), float64(480 - my), 1.0}, cam.ViewMatrix, cam.ProjectionMatrix, 640, 480)
	cam.MousePos = cam.Pos
	cam.MouseDir = mouseFar.Sub(mouseNear).Normalize()

	hitChunks := map[float64]ChunkCoord{}
	intersectPoints := map[ChunkCoord]vector.Vector3f{}
	for pos, chnk := range renderChunks {
		chnk.MouseHit = false
		boxPos := vector.Vector3f{float64(pos.X * CHUNK_BASE), float64(pos.Y * CHUNK_BASE), float64(pos.Z * CHUNK_BASE)}
		if hit, dist, hitPoint := CastRayAtBox(cam.MousePos, cam.MouseDir, boxPos, float64(CHUNK_BASE), 1.0, 128.0); hit {
			chnk.MouseHit = true
			hitChunks[dist] = pos
			intersectPoints[pos] = hitPoint
		}
	}

	if len(hitChunks) > 0 {
		chunkKeys := make([]float64, len(hitChunks))
		t := 0
		for k, _ := range hitChunks {
			chunkKeys[t] = k
			t++
		}
		sort.Float64s(chunkKeys)

		for t := 0; t < len(chunkKeys); t++ {
			chnkPos := hitChunks[chunkKeys[t]]
			chnk := chunkMap[chnkPos]
			hitBlocks := map[float64]BlockCoord{}
			for pos, blk := range chnk.data {
				// Note: Checking blk.visible seems to be reason for skipping to other side
				// at times.
				if true || blk.visible {
					boxPos := vector.Vector3f{
						float64(chnkPos.X*CHUNK_BASE + pos.X),
						float64(chnkPos.Y*CHUNK_BASE + pos.Y),
						float64(chnkPos.Z*CHUNK_BASE + pos.Z),
					}
					rayOrig := intersectPoints[chnkPos]
					if hit, dist, _ := CastRayAtBox(rayOrig, cam.MouseDir, boxPos, 1.0, 0.5, 16.0); hit {
						hitBlocks[dist] = pos
						break
					}
				}
			}

			if len(hitBlocks) > 0 {
				if !chnk.IsRebuilding {
					blockKeys := make([]float64, len(hitBlocks))
					t := 0
					for k, _ := range hitBlocks {
						blockKeys[t] = k
						t++
					}
					sort.Float64s(blockKeys)

					k := blockKeys[0]
					blockPos := hitBlocks[k]

					delete(chnk.data, blockPos)
					rebuildChunks[chnkPos] = chnk
					fmt.Printf("Found hit in chunk %v\n", chnkPos)
				}

				break
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
		posx := float64(pos.X * CHUNK_BASE)
		posy := float64(pos.Y * CHUNK_BASE)
		posz := float64(pos.Z * CHUNK_BASE)
		if cam.CubeInView(vector.Vector3f{posx, posy, posz}, float64(CHUNK_BASE)) != 2 {
			renderChunks[pos] = visibleChunks[pos]
		}
	}
}
