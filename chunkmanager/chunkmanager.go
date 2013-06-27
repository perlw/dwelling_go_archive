package chunkmanager

import (
	"bedrock/math/matrix"
	"bedrock/math/vector"
	"dwelling/camera"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const ChunkBase int = 16

type ChunkCoord struct {
	X, Y, Z int
}

type BlockCoord struct {
	X, Y, Z int
}

type Chunk struct {
	data         map[BlockCoord]*Block
	mesh         ChunkMesh
	IsLoaded     bool
	IsSetup      bool
	IsRebuilding bool
	MouseHit     bool
	position     ChunkCoord
}

type Block struct {
	visible   bool
	position  BlockCoord
	occlusion [6]float64
}

var chunkMap = map[ChunkCoord]*Chunk{}
var rebuildChunks = map[ChunkCoord]*Chunk{}
var visibleChunks = map[ChunkCoord]*Chunk{}
var renderChunks = map[ChunkCoord]*Chunk{}

var camPos = vector.Vector3f{0.0, 0.0, 0.0}
var camView = vector.Vector3f{0.0, 0.0, -1.0}

var debugMode = false

func Start() error {
	rand.Seed(time.Now().Unix())

	cubed := 4
	for x := 0; x < cubed; x++ {
		for z := 0; z < cubed; z++ {
			for y := 0; y < cubed; y++ {
				val := rand.Intn(6)

				var chunk *Chunk
				switch val {
				case 0:
					chunk = newCubeChunk(false)
				case 1:
					chunk = newCubeChunk(true)
				case 2:
					chunk = newPyramidChunk(false)
				case 3:
					chunk = newPyramidChunk(true)
				case 4:
					chunk = newSphereChunk()
				case 5:
					chunk = newWireCubeChunk()
				default:
					chunk = newCubeChunk(false)
				}
				//chunk := newFloatingRockChunk(ChunkCoord{x, y, z}, cubed)
				//chunk := newSimplexChunk(ChunkCoord{x, y, z}, cubed)
				chunk.position = ChunkCoord{x, y, z}
				chunkMap[chunk.position] = chunk
			}
		}
	}

	for chnkPos, chunk := range chunkMap {
		for blkPos, block := range chunk.data {
			block.occlusion = occlusion(chnkPos, blkPos, cubed)
		}
	}

	if err := setUpRenderer(); err != nil {
		return err
	}

	return nil
}

func SetDebug(mode bool) {
	debugMode = mode
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

func occlusion(chnkPos ChunkCoord, blkPos BlockCoord, size int) [6]float64 {
	occFactor := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	numRays := 16
	rays := goldenSectionSpiralRays(numRays)
	for t := 0; t < 6; t++ {
		for _, ray := range rays {
			currentStep := vector.Vector3f{}
			switch t {
			case FRONT:
				if ray.Z < 0.0 {
					continue
				}
				currentStep.X = float64((chnkPos.X*ChunkBase)+blkPos.X) + 0.5
				currentStep.Y = float64((chnkPos.Y*ChunkBase)+blkPos.Y) + 0.5
				currentStep.Z = float64((chnkPos.Z*ChunkBase)+blkPos.Z) + 1.0
			case BACK:
				if ray.Z > 0.0 {
					continue
				}
				currentStep.X = float64((chnkPos.X*ChunkBase)+blkPos.X) + 0.5
				currentStep.Y = float64((chnkPos.Y*ChunkBase)+blkPos.Y) + 0.5
				currentStep.Z = float64((chnkPos.Z*ChunkBase)+blkPos.Z) + 0.0
			case LEFT:
				if ray.X > 0.0 {
					continue
				}
				currentStep.X = float64((chnkPos.X*ChunkBase)+blkPos.X) + 0.0
				currentStep.Y = float64((chnkPos.Y*ChunkBase)+blkPos.Y) + 0.5
				currentStep.Z = float64((chnkPos.Z*ChunkBase)+blkPos.Z) + 0.5
			case RIGHT:
				if ray.X < 0.0 {
					continue
				}
				currentStep.X = float64((chnkPos.X*ChunkBase)+blkPos.X) + 1.0
				currentStep.Y = float64((chnkPos.Y*ChunkBase)+blkPos.Y) + 0.5
				currentStep.Z = float64((chnkPos.Z*ChunkBase)+blkPos.Z) + 0.5
			case TOP:
				if ray.Y < 0.0 {
					continue
				}
				currentStep.X = float64((chnkPos.X*ChunkBase)+blkPos.X) + 0.5
				currentStep.Y = float64((chnkPos.Y*ChunkBase)+blkPos.Y) + 1.0
				currentStep.Z = float64((chnkPos.Z*ChunkBase)+blkPos.Z) + 0.5
			case BOTTOM:
				if ray.Y > 0.0 {
					continue
				}
				currentStep.X = float64((chnkPos.X*ChunkBase)+blkPos.X) + 0.5
				currentStep.Y = float64((chnkPos.Y*ChunkBase)+blkPos.Y) + 0.0
				currentStep.Z = float64((chnkPos.Z*ChunkBase)+blkPos.Z) + 0.5
			default:
				continue
			}

			rayStep := ray.MulScalar(0.2)
			lastBlock := blkPos
			currentChnkPos := chnkPos
			currentChunk, _ := chunkMap[currentChnkPos]
			for {
				if currentStep.X < 0.0 || currentStep.X >= float64(ChunkBase*size) || currentStep.Y < 0.0 || currentStep.Y >= float64(ChunkBase*size) || currentStep.Z < 0.0 || currentStep.Z >= float64(ChunkBase*size) {
					occFactor[t] += 1.0
					break
				}

				recalc := false
				currBlock := BlockCoord{int(currentStep.X) - (currentChnkPos.X * ChunkBase), int(currentStep.Y) - (currentChnkPos.Y * ChunkBase), int(currentStep.Z) - (currentChnkPos.Z * ChunkBase)}
				if currBlock.X < 0 {
					currentChnkPos.X -= 1
					recalc = true
				}
				if currBlock.X >= ChunkBase {
					currentChnkPos.X += 1
					recalc = true
				}
				if currBlock.Y < 0 {
					currentChnkPos.Y -= 1
					recalc = true
				}
				if currBlock.Y >= ChunkBase {
					currentChnkPos.Y += 1
					recalc = true
				}
				if currBlock.Z < 0 {
					currentChnkPos.Z -= 1
					recalc = true
				}
				if currBlock.Z >= ChunkBase {
					currentChnkPos.Z += 1
					recalc = true
				}
				if recalc {
					currentChunk, _ = chunkMap[currentChnkPos]
					currBlock = BlockCoord{int(currentStep.X) - (currentChnkPos.X * ChunkBase), int(currentStep.Y) - (currentChnkPos.Y * ChunkBase), int(currentStep.Z) - (currentChnkPos.Z * ChunkBase)}
				}

				if currBlock.X != lastBlock.X || currBlock.Y != lastBlock.Y || currBlock.Z != lastBlock.Z {
					lastBlock = currBlock
					if _, ok := currentChunk.data[currBlock]; ok {
						break
					}
				}

				currentStep = currentStep.Add(rayStep)
			}
		}
	}

	for t := range occFactor {
		occFactor[t] = occFactor[t] / float64(numRays/2)
	}

	return occFactor
}

// Traces a ray against a box.
// Returns:
// bool Intersected?
// float64 Intersect distance
// vector.Vector3f Intersection point
func CastRayAtBox(rayOrig, rayDir, boxPos vector.Vector3f, boxSize, stepSize, maxDepth float64) (bool, float64, vector.Vector3f) {
	planeNormals := [6]vector.Vector3f{
		{0.0, 0.0, 1.0},
		{0.0, 0.0, -1.0},
		{-1.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{0.0, -1.0, 0.0},
	}

	planePos := [6]vector.Vector3f{
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z + boxSize},
		{boxPos.X + boxSize, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y + boxSize, boxPos.Z},
	}

	for depth := 0.0; depth < maxDepth; depth += stepSize {
		inside := 0
		rayStep := vector.Vector3f{}
		for t := 0; t < 6; t++ {
			d := -(vector.DotProduct(planeNormals[t], planePos[t]))
			rayStep = rayOrig.Add(rayDir.MulScalar(depth))
			deep := vector.DotProduct(planeNormals[t], rayStep) + d

			if deep > 0.0 {
				inside++
			}
		}
		if inside >= 6 {
			return true, depth, rayStep
		}
	}

	return false, 0.0, vector.Vector3f{}
}

// Checks if point is inside box
// Returns:
// bool In box?
func PointInBox(point, boxPos vector.Vector3f, boxSize float64) bool {
	planeNormals := [6]vector.Vector3f{
		{0.0, 0.0, 1.0},
		{0.0, 0.0, -1.0},
		{-1.0, 0.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{0.0, -1.0, 0.0},
	}

	planePos := [6]vector.Vector3f{
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z + boxSize},
		{boxPos.X + boxSize, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y, boxPos.Z},
		{boxPos.X, boxPos.Y + boxSize, boxPos.Z},
	}

	inside := 0
	for t := 0; t < 6; t++ {
		d := -(vector.DotProduct(planeNormals[t], planePos[t]))
		deep := vector.DotProduct(planeNormals[t], point) + d

		if deep > 0.0 {
			inside++
		}
	}
	if inside >= 6 {
		return true
	}

	return false
}

// TODO: Optmise chunk lookup by adding checked chunks to map copy and check against to avoid double and triple checking
// TODO: Rename func
func ClickedInChunk(mx, my int, cam *camera.Camera) {
	mouseNear, _ := matrix.Unproject(vector.Vector3f{float64(mx), float64(480 - my), 0.0}, cam.ViewMatrix, cam.ProjectionMatrix, 640, 480)
	mouseFar, _ := matrix.Unproject(vector.Vector3f{float64(mx), float64(480 - my), 1.0}, cam.ViewMatrix, cam.ProjectionMatrix, 640, 480)
	cam.MousePos = cam.Pos
	cam.MouseDir = mouseFar.Sub(mouseNear).Normalize()

	for _, chnk := range renderChunks {
		chnk.MouseHit = false
	}

	fmt.Println("---")
	chunkBase := float64(ChunkBase)
	for dist := 0.0; dist < 128.0; dist += 8.0 {
		rayStep := cam.MousePos.Add(cam.MouseDir.MulScalar(dist))

		pos := ChunkCoord{
			int(math.Trunc(rayStep.X / chunkBase)),
			int(math.Trunc(rayStep.Y / chunkBase)),
			int(math.Trunc(rayStep.Z / chunkBase)),
		}

		if chnk, ok := renderChunks[pos]; ok {
			fmt.Printf("Found potential chunk at %v...", pos)
			boxPos := vector.Vector3f{float64(pos.X * ChunkBase), float64(pos.Y * ChunkBase), float64(pos.Z * ChunkBase)}
			boxSize := chunkBase
			if PointInBox(rayStep, boxPos, boxSize) {
				fmt.Printf("hit!\n")
				if debugMode {
					chnk.MouseHit = true
				}

				startDist := dist - 8.0
				for dist := startDist; dist < startDist+24.0; dist += 0.5 {
					rayStep := cam.MousePos.Add(cam.MouseDir.MulScalar(dist))
					blkPos := BlockCoord{
						int(math.Trunc(rayStep.X)) - (pos.X * ChunkBase),
						int(math.Trunc(rayStep.Y)) - (pos.Y * ChunkBase),
						int(math.Trunc(rayStep.Z)) - (pos.Z * ChunkBase),
					}
					// Note: Checking blk.visible seems to be reason for skipping to other side
					// at times.
					if blk, ok := chnk.data[blkPos]; ok {
						if blk.visible {
							fmt.Printf("\tFound potential block at %v...", blkPos)

							boxPos := vector.Vector3f{
								float64(pos.X*ChunkBase + blkPos.X),
								float64(pos.Y*ChunkBase + blkPos.Y),
								float64(pos.Z*ChunkBase + blkPos.Z),
							}
							if PointInBox(rayStep, boxPos, 1.0) {
								fmt.Println("hit!\n")

								if !chnk.IsRebuilding {
									delete(chnk.data, blkPos)
									rebuildChunks[pos] = chnk
									recalcOcclusion(chnk, blkPos)
									rebuildNeighborsCheck(chnk.position, blkPos)
								}

								return
							} else {
								fmt.Println("nope...\n")
							}
						}
					}
				}
			} else {
				fmt.Println("nope...\n")
			}
		}
	}
}

func recalcOcclusion(chunk *Chunk, blkPos BlockCoord) {
	for index, block := range chunk.data {
		if block.visible {
			block.occlusion = occlusion(chunk.position, index, 4)
		}
	}
}

func rebuildNeighborsCheck(chnkPos ChunkCoord, blkPos BlockCoord) {
	if blkPos.X == 0 {
		neighborPos := ChunkCoord{chnkPos.X - 1, chnkPos.Y, chnkPos.Z}
		if chnk, ok := chunkMap[neighborPos]; ok {
			if _, ok := chnk.data[BlockCoord{ChunkBase - 1, blkPos.Y, blkPos.Z}]; ok {
				rebuildChunks[neighborPos] = chnk
			}
		}
	} else if blkPos.X == ChunkBase-1 {
		neighborPos := ChunkCoord{chnkPos.X + 1, chnkPos.Y, chnkPos.Z}
		if chnk, ok := chunkMap[neighborPos]; ok {
			if _, ok := chnk.data[BlockCoord{0, blkPos.Y, blkPos.Z}]; ok {
				rebuildChunks[neighborPos] = chnk
			}
		}
	}
	if blkPos.Y == 0 {
		neighborPos := ChunkCoord{chnkPos.X, chnkPos.Y - 1, chnkPos.Z}
		if chnk, ok := chunkMap[neighborPos]; ok {
			if _, ok := chnk.data[BlockCoord{blkPos.X, ChunkBase - 1, blkPos.Z}]; ok {
				rebuildChunks[neighborPos] = chnk
			}
		}
	} else if blkPos.Y == ChunkBase-1 {
		neighborPos := ChunkCoord{chnkPos.X, chnkPos.Y + 1, chnkPos.Z}
		if chnk, ok := chunkMap[neighborPos]; ok {
			if _, ok := chnk.data[BlockCoord{blkPos.X, 0, blkPos.Z}]; ok {
				rebuildChunks[neighborPos] = chnk
			}
		}
	}
	if blkPos.Z == 0 {
		neighborPos := ChunkCoord{chnkPos.X, chnkPos.Y, chnkPos.Z - 1}
		if chnk, ok := chunkMap[neighborPos]; ok {
			if _, ok := chnk.data[BlockCoord{blkPos.X, blkPos.Y, ChunkBase - 1}]; ok {
				rebuildChunks[neighborPos] = chnk
			}
		}
	} else if blkPos.Z == ChunkBase-1 {
		neighborPos := ChunkCoord{chnkPos.X, chnkPos.Y, chnkPos.Z + 1}
		if chnk, ok := chunkMap[neighborPos]; ok {
			if _, ok := chnk.data[BlockCoord{blkPos.X, blkPos.Y, 0}]; ok {
				rebuildChunks[neighborPos] = chnk
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

type RebuildData struct {
	vertexBuffers [6][]float32
	indexBuffers  [6][]uint32
	occBuffers    [6][]float32
	chunk         *Chunk
}

var rebuildCh = make(chan RebuildData)
var numRebuilding = 0

func updateRebuildList() {
	select {
	case rebuildData := <-rebuildCh:
		rebuildData.chunk.SetChunkMesh(rebuildData)
		rebuildData.chunk.IsRebuilding = false
		fmt.Printf("rebuilds: %v rebuilt.\n", rebuildData.chunk.position)

		numRebuilding--
		if numRebuilding <= 0 {
			numRebuilding = 0
			fmt.Println("rebuilds: All done")
		}
	default:
	}

	for index, chnk := range rebuildChunks {
		// TODO: Performance throttling?
		if numRebuilding < 2 {
			numRebuilding++
			chnk.IsRebuilding = true
			fmt.Printf("rebuilds: (%d/%d) - Adding %v to rebuild queue.\n", numRebuilding, 2, index)
			go chnk.CreateVertexData(rebuildCh)
			delete(rebuildChunks, index)
		}
	}
}

func updateVisibilityList(cam *camera.Camera) {
	// TODO: Add chunk range limit
	for t, chnk := range chunkMap {
		if chnk.IsLoaded && chnk.IsSetup {
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
		posx := float64(pos.X * ChunkBase)
		posy := float64(pos.Y * ChunkBase)
		posz := float64(pos.Z * ChunkBase)
		if cam.CubeInView(vector.Vector3f{posx, posy, posz}, float64(ChunkBase)) != 2 {
			renderChunks[pos] = visibleChunks[pos]
		}
	}
}
