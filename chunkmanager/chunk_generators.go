package chunkmanager

import (
	"bedrock/math/simplex"
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

func newSimplexChunk(pos ChunkCoord, size int) *Chunk {
	chunk := &Chunk{}

	// Expecting a perfect cube world
	worldMax := float64(ChunkBase * size)
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

func newFloatingRockChunk(pos ChunkCoord, size int) *Chunk {
	chunk := &Chunk{}

	// Expecting a perfect cube world
	worldMax := float64(ChunkBase * size)
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

/*func newRayTestChunk() *Chunk {
	chunk := &Chunk{}

	chunk.data = map[BlockCoord]*Block{}
	startX := ChunkBase / 2
	startY := ChunkBase / 2
	startZ := ChunkBase / 2

	index := BlockCoord{startX, startY, startZ}
	chunk.data[index] = &Block{
		visible:  false,
		position: index,
	}

	rays := goldenSectionSpiralRays(16)
	for _, ray := range rays {
		rayStep := ray.MulScalar(0.2)
		currentStep := vector.Vector3f{float64(startX) + 0.5, float64(startY) + 0.5, float64(startZ) + 0.5}
		lastBlock := BlockCoord{startX, startY, startZ}
		for {
			currBlock := BlockCoord{int(currentStep.X), int(currentStep.Y), int(currentStep.Z)}

			if currBlock.X < 0 || currBlock.X >= ChunkBase || currBlock.Y < 0 || currBlock.Y >= ChunkBase || currBlock.Z < 0 || currBlock.Z >= ChunkBase {
				break
			}

			if currBlock.X != lastBlock.X || currBlock.Y != lastBlock.Y || currBlock.Z != lastBlock.Z {
				lastBlock = currBlock

				index := currBlock
				chunk.data[index] = &Block{
					visible:  false,
					position: index,
				}
			}

			currentStep = currentStep.Add(rayStep)
		}
	}

	chunk.IsLoaded = true
	chunk.MouseHit = false

	return chunk
}

// Port of Golden Section Spiral python code
// from http://www.softimageblog.com/archives/115
func goldenSectionSpiralRays(numRays int) []vector.Vector3f {
	rays := []vector.Vector3f{}

	increment := math.Pi * (3.0 - math.Sqrt(5.0))
	offset := 2.0 / float64(numRays)
	for t := 0; t < numRays; t++ {
		y := (float64(t) * offset) - 1.0 + (offset / 2.0)
		r := math.Sqrt(1 - (y * y))
		phi := float64(t) * increment

		rays = append(rays, vector.Vector3f{math.Cos(phi) * r, y, math.Sin(phi) * r})
	}

	return rays
}

func occlusion(chunk *Chunk, pos BlockCoord) float64 {
	occFactor := 0.0
	numRays := 16
	rays := goldenSectionSpiralRays(numRays)
	for _, ray := range rays {
		rayStep := ray.MulScalar(0.2)
		currentStep := vector.Vector3f{float64(pos.X) + 0.5, float64(pos.Y) + 0.5, float64(pos.Z) + 0.5}
		lastBlock := pos
		for {
			currBlock := BlockCoord{int(currentStep.X), int(currentStep.Y), int(currentStep.Z)}

			if currBlock.X < 0 || currBlock.X >= ChunkBase || currBlock.Y < 0 || currBlock.Y >= ChunkBase || currBlock.Z < 0 || currBlock.Z >= ChunkBase {
				occFactor += 1.0
				break
			}

			if currBlock.X != lastBlock.X || currBlock.Y != lastBlock.Y || currBlock.Z != lastBlock.Z {
				lastBlock = currBlock
				if _, ok := chunk.data[currBlock]; ok {
					break
				}
			}

			currentStep = currentStep.Add(rayStep)
		}
	}
	return occFactor / float64(numRays)
}
*/
