package chunk

const CHUNK_BASE int = 16

type BlockCoord struct {
	X, Y, Z int
}

type Chunk struct {
	data map[BlockCoord]Block
	mesh ChunkMesh
}

type Block struct {
}

func NewPyramidChunk() *Chunk {
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

	chunk.UpdateChunkMesh()

	return chunk
}

func NewCubeChunk() *Chunk {
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

	chunk.UpdateChunkMesh()

	return chunk
}
