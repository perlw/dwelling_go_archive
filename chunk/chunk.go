package chunk

const CHUNK_BASE int = 16

type ChunkCoord struct {
	X, Y, Z int
}

type Chunk struct {
	data map[ChunkCoord]Block
	mesh ChunkMesh
}

type Block struct {
}

func NewPyramidChunk() *Chunk {
	chunk := &Chunk{}

	chunk.data = map[ChunkCoord]Block{}
	for y := 0; y < CHUNK_BASE/2; y++ {
		for x := y; x < CHUNK_BASE-y; x++ {
			for z := y; z < CHUNK_BASE-y; z++ {
				index := ChunkCoord{x, y, z}
				chunk.data[index] = Block{}
			}
		}
	}

	chunk.UpdateChunkMesh()

	return chunk
}

func NewCubeChunk() *Chunk {
	chunk := &Chunk{}

	chunk.data = map[ChunkCoord]Block{}
	for y := 0; y < CHUNK_BASE; y++ {
		for x := 0; x < CHUNK_BASE; x++ {
			for z := 0; z < CHUNK_BASE; z++ {
				index := ChunkCoord{x, y, z}
				chunk.data[index] = Block{}
			}
		}
	}

	chunk.UpdateChunkMesh()

	return chunk
}
