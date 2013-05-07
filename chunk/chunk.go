package chunk

const CHUNK_BASE int = 16

type ChunkCoord struct {
	x, y, z int
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
				index := ChunkCoord{x: x, y: y, z: z}
				chunk.data[index] = Block{}
			}
		}
	}

	chunk.UpdateChunkMesh()

	return chunk
}
