package chunkmanager

import (
	"dwelling/math/simplex"
	"math"
	"math/rand"
)

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

func newCubeChunk(random bool) *Chunk {
	chunk := &Chunk{}

	chunk.data = map[BlockCoord]*Block{}
	for y := 0; y < ChunkBase; y++ {
		for x := 0; x < ChunkBase; x++ {
			for z := 0; z < ChunkBase; z++ {
				val := 1
				if random == false {
					val = rand.Intn(2)
				}

				if val == 1 {
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

func newSimplexChunk(pos ChunkCoord) *Chunk {
	chunk := &Chunk{}

	// Expecting 4x4x4 world
	worldMax := float64(ChunkBase * 4)
	chunk.data = map[BlockCoord]*Block{}
	for y := 0; y < ChunkBase; y++ {
		for x := 0; x < ChunkBase; x++ {
			for z := 0; z < ChunkBase; z++ {
				bX := float64((pos.X*ChunkBase)+x) / worldMax
				bY := float64((pos.Y*ChunkBase)+y) / worldMax
				bZ := float64((pos.Z*ChunkBase)+z) / worldMax
				noise := simplex.Noise(bX*3.0, bY*3.0, bZ*3.0)
				if noise > 1.1 {
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

func newFloatingRockChunk(pos ChunkCoord) *Chunk {
	chunk := &Chunk{}

	// Expecting 4x4x4 world
	worldMax := float64(ChunkBase * 4)
	chunk.data = map[BlockCoord]*Block{}
	for y := 0; y < ChunkBase; y++ {
		for x := 0; x < ChunkBase; x++ {
			for z := 0; z < ChunkBase; z++ {
				bX := float64((pos.X*ChunkBase)+x) / worldMax
				bY := float64((pos.Y*ChunkBase)+y) / worldMax
				bZ := float64((pos.Z*ChunkBase)+z) / worldMax

				plateauFallof := 0.0
				if bY <= 0.8 {
					plateauFallof = 1.0
				} else if 0.8 < bY && bY < 0.9 {
					plateauFallof = 1.0 - (bY-0.8)*10.0
				}

				centerFallof := 0.2 / (math.Pow((bX-0.5)*1.5, 2) + math.Pow((bY-1.0)*0.8, 2) + math.Pow((bZ-0.5)*1.5, 2))

				caves := math.Pow(simplex.Noise(bX*5.0, bY*5.0, bZ*5.0), 3)
				density := 0.0
				if caves >= 0.5 {
					density = simplex.NoiseOctave(5, bX, bY*0.5, bZ) * centerFallof * plateauFallof
					density *= math.Pow(simplex.Noise((bX+1.0)*3.0, (bY+1.0)*3.0, ((bZ+1.0)*3.0)+0.4), 1.8)
				}

				if density > 3.1 {
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
