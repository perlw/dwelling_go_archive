package chunkmanager

import "math"

func newPyramidChunk(invert bool) *Chunk {
	chunk := &Chunk{}

	chunk.data = map[BlockCoord]*Block{}
	for y := 0; y < ChunkBase/2; y++ {
		for x := y; x < ChunkBase-y; x++ {
			for z := y; z < ChunkBase-y; z++ {
				index := BlockCoord{x, y, z}
				if invert {
					index.Y = ChunkBase - index.Y - 1
				}
				chunk.data[index] = &Block{
					visible:  false,
					position: index,
				}
			}
		}
	}

	chunk.IsLoaded = true
	chunk.MouseHit = false

	return chunk
}

func newCubeChunk() *Chunk {
	chunk := &Chunk{}

	chunk.data = map[BlockCoord]*Block{}
	for y := 0; y < ChunkBase; y++ {
		for x := 0; x < ChunkBase; x++ {
			for z := 0; z < ChunkBase; z++ {
				index := BlockCoord{x, y, z}
				chunk.data[index] = &Block{
					visible:  false,
					position: index,
				}
			}
		}
	}

	chunk.IsLoaded = true
	chunk.MouseHit = false

	return chunk
}

func newSphereChunk() *Chunk {
	chunk := &Chunk{}

	halfChunk := float64(ChunkBase) / 2.0
	chunk.data = map[BlockCoord]*Block{}
	for y := 0; y < ChunkBase; y++ {
		for x := 0; x < ChunkBase; x++ {
			for z := 0; z < ChunkBase; z++ {
				xx := math.Pow(float64(x)-halfChunk, 2)
				yy := math.Pow(float64(y)-halfChunk, 2)
				zz := math.Pow(float64(z)-halfChunk, 2)
				dist := math.Sqrt(xx + yy + zz)
				if dist > halfChunk-1.0 && dist <= halfChunk {
					index := BlockCoord{x, y, z}
					chunk.data[index] = &Block{
						visible:  false,
						position: index,
					}
				}
			}
		}
	}

	chunk.IsLoaded = true
	chunk.MouseHit = false

	return chunk
}

func newWireCubeChunk() *Chunk {
	chunk := &Chunk{}

	chunk.data = map[BlockCoord]*Block{}
	for y := 0; y < ChunkBase; y++ {
		for x := 0; x < ChunkBase; x++ {
			for z := 0; z < ChunkBase; z++ {
				aa := ((x+1)%ChunkBase-1 == 0 || (x+1)%ChunkBase == 0)
				bb := ((y+1)%ChunkBase-1 == 0 || (y+1)%ChunkBase == 0)
				cc := ((z+1)%ChunkBase-1 == 0 || (z+1)%ChunkBase == 0)
				if (aa && bb) || (bb && cc) || (aa && cc) {
					index := BlockCoord{x, y, z}
					chunk.data[index] = &Block{
						visible:  false,
						position: index,
					}
				}
			}
		}
	}

	chunk.IsLoaded = true
	chunk.MouseHit = false

	return chunk
}
