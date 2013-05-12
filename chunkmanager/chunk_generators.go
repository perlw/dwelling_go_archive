package chunkmanager

func newPyramidChunk() *Chunk {
	chunk := &Chunk{}

	chunk.data = map[BlockCoord]Block{}
	for y := 0; y < CHUNK_BASE/2; y++ {
		for x := y; x < CHUNK_BASE-y; x++ {
			for z := y; z < CHUNK_BASE-y; z++ {
				index := BlockCoord{x, y, z}
				chunk.data[index] = Block{}
			}
		}
	}

	chunk.IsLoaded = true

	return chunk
}

func newCubeChunk() *Chunk {
	chunk := &Chunk{}

	chunk.data = map[BlockCoord]Block{}
	for y := 0; y < CHUNK_BASE; y++ {
		for x := 0; x < CHUNK_BASE; x++ {
			for z := 0; z < CHUNK_BASE; z++ {
				index := BlockCoord{x, y, z}
				chunk.data[index] = Block{}
			}
		}
	}

	chunk.IsLoaded = true

	return chunk
}
