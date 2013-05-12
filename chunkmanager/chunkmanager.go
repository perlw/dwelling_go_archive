package chunkmanager

import (
	"dwelling/camera"
	"dwelling/chunk"
	"dwelling/math/vector"
	"fmt"
)

type ChunkCoord struct {
	X, Y, Z int
}

var chunkMap = map[ChunkCoord]*chunk.Chunk{}
var visibleChunks = map[ChunkCoord]*chunk.Chunk{}
var renderChunks = map[ChunkCoord]*chunk.Chunk{}

var camPos = vector.Vector3f{0.0, 0.0, 0.0}
var camView = vector.Vector3f{0.0, 0.0, -1.0}

func Start() {
	cubed := 4
	for x := 0; x < cubed; x++ {
		for z := 0; z < cubed; z++ {
			for y := 0; y < cubed; y++ {
				chunkMap[ChunkCoord{x, y, z}] = chunk.NewCubeChunk()
			}
		}
	}
}

func Update(cam *camera.Camera) {
	updateVisibilityList(cam)

	if camPos != cam.Pos || camView != cam.Rot {
		updateRenderList(cam)

		camPos = cam.Pos
		camView = cam.Rot
	}
}

func updateVisibilityList(cam *camera.Camera) {
	// TODO: Add chunk range limit
	for t := range chunkMap {
		if _, ok := visibleChunks[t]; !ok {
			fmt.Printf("Added chunk at %v to visible list.\n", t)
			visibleChunks[t] = chunkMap[t]
		}
	}
}

func updateRenderList(cam *camera.Camera) {
	renderChunks = map[ChunkCoord]*chunk.Chunk{}

	for pos := range visibleChunks {
		posx := float64(pos.X * chunk.CHUNK_BASE)
		posy := float64(pos.Y * chunk.CHUNK_BASE)
		posz := float64(pos.Z * chunk.CHUNK_BASE)
		if cam.CubeInView(vector.Vector3f{posx, posy, posz}, float64(chunk.CHUNK_BASE)) != 2 {
			renderChunks[pos] = visibleChunks[pos]
		}
	}
}
