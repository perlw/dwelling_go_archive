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
	hitChunks := map[float64]ChunkCoord{}
	intersectPoints := map[ChunkCoord]vector.Vector3f{}
	for pos, chnk := range renderChunks {
		x := float64(pos.X * CHUNK_BASE)
		y := float64(pos.Y * CHUNK_BASE)
		z := float64(pos.Z * CHUNK_BASE)
		x1 := x + float64(CHUNK_BASE)
		y1 := y + float64(CHUNK_BASE)
		z1 := z + float64(CHUNK_BASE)
		planePos := [6]vector.Vector3f{
			{x, y, z},
			{x, y, z1},
			{x1, y, z},
			{x, y, z},
			{x, y, z},
			{x, y1, z},
		}

		chnk.MouseHit = false
		for dist := 0.0; dist < 128.0; dist += 1.0 {
			inside := 0
			rayPos := vector.Vector3f{}
			for t := 0; t < 6; t++ {
				d := -(vector.DotProduct(planeNormals[t], planePos[t]))
				rayPos = cam.MousePos.Add(cam.MouseDir.MulScalar(dist))
				deep := vector.DotProduct(planeNormals[t], rayPos) + d

				if deep > 0.0 {
					inside++
				}
			}
			if inside >= 6 {
				chnk.MouseHit = true
				hitChunks[dist] = pos
				intersectPoints[pos] = rayPos
				break
			}
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
				if blk.visible {
					x := float64((chnkPos.X * CHUNK_BASE) + pos.X)
					y := float64((chnkPos.Y * CHUNK_BASE) + pos.Y)
					z := float64((chnkPos.Z * CHUNK_BASE) + pos.Z)
					x1 := x + 1.0
					y1 := y + 1.0
					z1 := z + 1.0
					planePos := [6]vector.Vector3f{
						{x, y, z},
						{x, y, z1},
						{x1, y, z},
						{x, y, z},
						{x, y, z},
						{x, y1, z},
					}

					rayOrig := intersectPoints[chnkPos]
					for dist := 0.0; dist < 16.0; dist += 0.5 {
						inside := 0
						for t := 0; t < 6; t++ {
							d := -(vector.DotProduct(planeNormals[t], planePos[t]))
							rayPos := rayOrig.Add(cam.MouseDir.MulScalar(dist))
							deep := vector.DotProduct(planeNormals[t], rayPos) + d

							if deep > 0.0 {
								inside++
							}
						}
						if inside >= 6 {
							hitBlocks[dist] = pos
							break
						}
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
