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

const CHUNK_BASE int = 16

type ChunkCoord struct {
	X, Y, Z int
}

type BlockCoord struct {
	X, Y, Z int
}

type Chunk struct {
	data     map[BlockCoord]Block
	mesh     ChunkMesh
	IsLoaded bool
	IsSetup  bool
	MouseHit bool
}

type Block struct {
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
				val := rand.Intn(2)

				if val == 0 {
					chunkMap[ChunkCoord{x, y, z}] = newCubeChunk()
				} else {
					chunkMap[ChunkCoord{x, y, z}] = newPyramidChunk()
				}
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

func ClickedInChunk(mx, my int, cam *camera.Camera) {
	mouseNear, _ := matrix.Unproject(vector.Vector3f{float64(mx), float64(480 - my), 0.0}, cam.ViewMatrix, cam.ProjectionMatrix, 640, 480)
	mouseFar, _ := matrix.Unproject(vector.Vector3f{float64(mx), float64(480 - my), 1.0}, cam.ViewMatrix, cam.ProjectionMatrix, 640, 480)
	cam.MousePos = cam.Pos
	cam.MouseDir = mouseFar.Sub(mouseNear).Normalize()

	planeNormals := [6]vector.Vector3f{
		{0.0, 0.0, 1.0},
		{0.0, 0.0, -1.0},
		{-1.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{0.0, -1.0, 0.0},
	}
	for pos, chnk := range renderChunks {
		x := float64(pos.X * CHUNK_BASE)
		y := float64(pos.Y * CHUNK_BASE)
		z := float64(pos.Z * CHUNK_BASE)
		x1 := x + float64(CHUNK_BASE)
		y1 := y + float64(CHUNK_BASE)
		z1 := z + float64(CHUNK_BASE)
		halfBase := float64(CHUNK_BASE / 2)
		planePos := [6]vector.Vector3f{
			{x, y, z},
			{x, y, z1},
			{x1, y, z},
			{x, y, z},
			{x, y, z},
			{x, y1, z},
		}

		chnk.MouseHit = false
		inside := 0
		for t := 0; t < 6; t++ {
			cubeMid := vector.Vector3f{x + halfBase, y + halfBase, z + halfBase}
			tmp := cubeMid.Sub(cam.MousePos)
			dist := math.Sqrt((tmp.X * tmp.X) + (tmp.Y * tmp.Y) + (tmp.Z * tmp.Z))

			d := -(vector.DotProduct(planeNormals[t], planePos[t]))
			deep := vector.DotProduct(planeNormals[t], cam.MousePos.Add(cam.MouseDir.MulScalar(dist))) + d

			if deep > 0.0 {
				inside++
			}
		}
		if inside >= 6 {
			chnk.MouseHit = true
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
			chnk.UpdateChunkMesh(pos)
		}
	}
}

func updateRebuildList() {
	for pos, chnk := range rebuildChunks {
		chnk.UpdateChunkMesh(pos)
	}

	rebuildChunks = map[ChunkCoord]*Chunk{}
}

func updateVisibilityList(cam *camera.Camera) {
	// TODO: Add chunk range limit
	for t, chnk := range chunkMap {
		if chnk.IsLoaded && chnk.IsSetup {
			if _, ok := visibleChunks[t]; !ok {
				fmt.Printf("Added chunk at %v to visible list.\n", t)
				visibleChunks[t] = chunkMap[t]
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
