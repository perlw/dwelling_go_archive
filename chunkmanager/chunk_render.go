package chunkmanager

import (
	"bedrock/shader"
	"bedrock/math/matrix"
	"bedrock/math/vector"
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"unsafe"
)

const (
	FRONT int = iota
	BACK
	LEFT
	RIGHT
	TOP
	BOTTOM
)

var chunkNormals = [6]vector.Vector3f{
	{0.0, 0.0, 1.0},
	{0.0, 0.0, -1.0},
	{-1.0, 0.0, 0.0},
	{1.0, 0.0, 0.0},
	{0.0, 1.0, 0.0},
	{0.0, -1.0, 0.0},
}

type ChunkMesh struct {
	vao             [6]gl.Uint
	vertexBufferIds [6]gl.Uint
	indexBufferIds  [6]gl.Uint
	occBufferIds    [6]gl.Uint
	numVertices     [6]gl.Sizei
	numIndices      [6]gl.Sizei
}

func appendChunkFace(faceBuffer *[]float32, indexBuffer *[]uint32, occBuffer *[]float32, occFactor [6]float64, x, y, z float32, face int) {
	var vertices [4]vector.Vector4f

	switch face {
	case FRONT:
		avgOcc := [4]float64{
			(occFactor[FRONT] + occFactor[LEFT] + occFactor[BOTTOM]) / 3.0,
			(occFactor[FRONT] + occFactor[RIGHT] + occFactor[BOTTOM]) / 3.0,
			(occFactor[FRONT] + occFactor[RIGHT] + occFactor[TOP]) / 3.0,
			(occFactor[FRONT] + occFactor[LEFT] + occFactor[TOP]) / 3.0,
		}
		vertices = [4]vector.Vector4f{
			{float64(x), float64(y), float64(z) + 1.0, avgOcc[0]},
			{float64(x) + 1.0, float64(y), float64(z) + 1.0, avgOcc[1]},
			{float64(x) + 1.0, float64(y) + 1.0, float64(z) + 1.0, avgOcc[2]},
			{float64(x), float64(y) + 1.0, float64(z) + 1.0, avgOcc[3]},
		}

	case BACK:
		avgOcc := [4]float64{
			(occFactor[BACK] + occFactor[RIGHT] + occFactor[TOP]) / 3.0,
			(occFactor[BACK] + occFactor[RIGHT] + occFactor[BOTTOM]) / 3.0,
			(occFactor[BACK] + occFactor[LEFT] + occFactor[BOTTOM]) / 3.0,
			(occFactor[BACK] + occFactor[LEFT] + occFactor[TOP]) / 3.0,
		}
		vertices = [4]vector.Vector4f{
			{float64(x) + 1.0, float64(y) + 1.0, float64(z), avgOcc[0]},
			{float64(x) + 1.0, float64(y), float64(z), avgOcc[1]},
			{float64(x), float64(y), float64(z), avgOcc[2]},
			{float64(x), float64(y) + 1.0, float64(z), avgOcc[3]},
		}

	case LEFT:
		avgOcc := [4]float64{
			(occFactor[LEFT] + occFactor[BACK] + occFactor[BOTTOM]) / 3.0,
			(occFactor[LEFT] + occFactor[FRONT] + occFactor[BOTTOM]) / 3.0,
			(occFactor[LEFT] + occFactor[FRONT] + occFactor[TOP]) / 3.0,
			(occFactor[LEFT] + occFactor[BACK] + occFactor[TOP]) / 3.0,
		}
		vertices = [4]vector.Vector4f{
			{float64(x), float64(y), float64(z), avgOcc[0]},
			{float64(x), float64(y), float64(z) + 1.0, avgOcc[1]},
			{float64(x), float64(y) + 1.0, float64(z) + 1.0, avgOcc[2]},
			{float64(x), float64(y) + 1.0, float64(z), avgOcc[3]},
		}

	case RIGHT:
		avgOcc := [4]float64{
			(occFactor[RIGHT] + occFactor[FRONT] + occFactor[TOP]) / 3.0,
			(occFactor[RIGHT] + occFactor[FRONT] + occFactor[BOTTOM]) / 3.0,
			(occFactor[RIGHT] + occFactor[BACK] + occFactor[BOTTOM]) / 3.0,
			(occFactor[RIGHT] + occFactor[BACK] + occFactor[TOP]) / 3.0,
		}
		vertices = [4]vector.Vector4f{
			{float64(x) + 1.0, float64(y) + 1.0, float64(z) + 1.0, avgOcc[0]},
			{float64(x) + 1.0, float64(y), float64(z) + 1.0, avgOcc[1]},
			{float64(x) + 1.0, float64(y), float64(z), avgOcc[2]},
			{float64(x) + 1.0, float64(y) + 1.0, float64(z), avgOcc[3]},
		}

	case TOP:
		avgOcc := [4]float64{
			(occFactor[TOP] + occFactor[RIGHT] + occFactor[FRONT]) / 3.0,
			(occFactor[TOP] + occFactor[RIGHT] + occFactor[BACK]) / 3.0,
			(occFactor[TOP] + occFactor[LEFT] + occFactor[BACK]) / 3.0,
			(occFactor[TOP] + occFactor[LEFT] + occFactor[FRONT]) / 3.0,
		}
		vertices = [4]vector.Vector4f{
			{float64(x) + 1.0, float64(y) + 1.0, float64(z) + 1.0, avgOcc[0]},
			{float64(x) + 1.0, float64(y) + 1.0, float64(z), avgOcc[1]},
			{float64(x), float64(y) + 1.0, float64(z), avgOcc[2]},
			{float64(x), float64(y) + 1.0, float64(z) + 1.0, avgOcc[3]},
		}

	case BOTTOM:
		avgOcc := [4]float64{
			(occFactor[BOTTOM] + occFactor[LEFT] + occFactor[BACK]) / 3.0,
			(occFactor[BOTTOM] + occFactor[RIGHT] + occFactor[BACK]) / 3.0,
			(occFactor[BOTTOM] + occFactor[RIGHT] + occFactor[FRONT]) / 3.0,
			(occFactor[BOTTOM] + occFactor[LEFT] + occFactor[FRONT]) / 3.0,
		}
		vertices = [4]vector.Vector4f{
			{float64(x), float64(y), float64(z), avgOcc[0]},
			{float64(x) + 1.0, float64(y), float64(z), avgOcc[1]},
			{float64(x) + 1.0, float64(y), float64(z) + 1.0, avgOcc[2]},
			{float64(x), float64(y), float64(z) + 1.0, avgOcc[3]},
		}

	default:
		return
	}

	vertIds := [4]uint32{}
	for index, vertex := range vertices {
		(*faceBuffer) = append((*faceBuffer), float32(vertex.X), float32(vertex.Y), float32(vertex.Z))
		(*occBuffer) = append((*occBuffer), float32(vertex.W))
		vertIds[index] = uint32((len((*faceBuffer)) - 3) / 3)
	}

	a := vertIds[0]
	b := vertIds[1]
	c := vertIds[2]
	d := vertIds[3]

	(*indexBuffer) = append((*indexBuffer),
		a, b, c,
		c, d, a,
	)
}

func createMeshBuffer(faceBuffer *[]float32, size int) gl.Uint {
	var buffer gl.Uint
	sizeFloat := int(unsafe.Sizeof([1]float32{}))
	bufferPtr := (*faceBuffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*size), gl.Pointer(&bufferPtr[0]), gl.STATIC_DRAW)

	return buffer
}

func createIndexBuffer(faceBuffer *[]uint32, size int) gl.Uint {
	var buffer gl.Uint
	sizeInt := int(unsafe.Sizeof([1]uint32{}))
	bufferPtr := (*faceBuffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(sizeInt*size), gl.Pointer(&bufferPtr[0]), gl.STATIC_DRAW)

	return buffer
}

func (chunk *Chunk) CreateVertexData(rebuildCh chan<- RebuildData) {
	chunks := GetChunksAroundChunk(chunk.position)

	vertexBuffers := [6][]float32{}
	indexBuffers := [6][]uint32{}
	occBuffers := [6][]float32{}
	for pos := range chunk.data {
		x := float32(pos.X)
		y := float32(pos.Y)
		z := float32(pos.Z)

		sides := 0
		if _, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z + 1}]; !ok {
			skip := false
			if pos.Z == ChunkBase-1 && chunks[FRONT] != nil {
				if _, ok := chunks[FRONT].data[BlockCoord{pos.X, pos.Y, 0}]; ok {
					skip = true
				}
			}

			if !skip {
				occFactor := [6]float64{
					chunk.data[pos].occlusion[FRONT],
					chunk.data[pos].occlusion[FRONT],
					chunk.data[pos].occlusion[FRONT],
					chunk.data[pos].occlusion[FRONT],
					chunk.data[pos].occlusion[FRONT],
					chunk.data[pos].occlusion[FRONT],
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y + 1, pos.Z}]; ok {
					occFactor[TOP] = neighbor.occlusion[FRONT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y - 1, pos.Z}]; ok {
					occFactor[BOTTOM] = neighbor.occlusion[FRONT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X - 1, pos.Y, pos.Z}]; ok {
					occFactor[LEFT] = neighbor.occlusion[FRONT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X + 1, pos.Y, pos.Z}]; ok {
					occFactor[RIGHT] = neighbor.occlusion[FRONT]
				}

				sides++
				appendChunkFace(&vertexBuffers[FRONT], &indexBuffers[FRONT], &occBuffers[FRONT], occFactor, x, y, z, FRONT)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z - 1}]; !ok {
			skip := false
			if pos.Z == 0 && chunks[BACK] != nil {
				if _, ok := chunks[BACK].data[BlockCoord{pos.X, pos.Y, ChunkBase - 1}]; ok {
					skip = true
				}
			}

			if !skip {
				occFactor := [6]float64{
					chunk.data[pos].occlusion[BACK],
					chunk.data[pos].occlusion[BACK],
					chunk.data[pos].occlusion[BACK],
					chunk.data[pos].occlusion[BACK],
					chunk.data[pos].occlusion[BACK],
					chunk.data[pos].occlusion[BACK],
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y + 1, pos.Z}]; ok {
					occFactor[TOP] = neighbor.occlusion[BACK]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y - 1, pos.Z}]; ok {
					occFactor[BOTTOM] = neighbor.occlusion[BACK]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X - 1, pos.Y, pos.Z}]; ok {
					occFactor[LEFT] = neighbor.occlusion[BACK]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X + 1, pos.Y, pos.Z}]; ok {
					occFactor[RIGHT] = neighbor.occlusion[BACK]
				}

				sides++
				appendChunkFace(&vertexBuffers[BACK], &indexBuffers[BACK], &occBuffers[BACK], occFactor, x, y, z, BACK)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X - 1, pos.Y, pos.Z}]; !ok {
			skip := false
			if pos.X == 0 && chunks[LEFT] != nil {
				if _, ok := chunks[LEFT].data[BlockCoord{ChunkBase - 1, pos.Y, pos.Z}]; ok {
					skip = true
				}
			}

			if !skip {
				occFactor := [6]float64{
					chunk.data[pos].occlusion[LEFT],
					chunk.data[pos].occlusion[LEFT],
					chunk.data[pos].occlusion[LEFT],
					chunk.data[pos].occlusion[LEFT],
					chunk.data[pos].occlusion[LEFT],
					chunk.data[pos].occlusion[LEFT],
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y + 1, pos.Z}]; ok {
					occFactor[TOP] = neighbor.occlusion[LEFT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y - 1, pos.Z}]; ok {
					occFactor[BOTTOM] = neighbor.occlusion[LEFT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z - 1}]; ok {
					occFactor[BACK] = neighbor.occlusion[LEFT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z + 1}]; ok {
					occFactor[FRONT] = neighbor.occlusion[LEFT]
				}

				sides++
				appendChunkFace(&vertexBuffers[LEFT], &indexBuffers[LEFT], &occBuffers[LEFT], occFactor, x, y, z, LEFT)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X + 1, pos.Y, pos.Z}]; !ok {
			skip := false
			if pos.X == ChunkBase-1 && chunks[RIGHT] != nil {
				if _, ok := chunks[RIGHT].data[BlockCoord{0, pos.Y, pos.Z}]; ok {
					skip = true
				}
			}

			if !skip {
				occFactor := [6]float64{
					chunk.data[pos].occlusion[RIGHT],
					chunk.data[pos].occlusion[RIGHT],
					chunk.data[pos].occlusion[RIGHT],
					chunk.data[pos].occlusion[RIGHT],
					chunk.data[pos].occlusion[RIGHT],
					chunk.data[pos].occlusion[RIGHT],
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y + 1, pos.Z}]; ok {
					occFactor[TOP] = neighbor.occlusion[RIGHT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y - 1, pos.Z}]; ok {
					occFactor[BOTTOM] = neighbor.occlusion[RIGHT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z - 1}]; ok {
					occFactor[BACK] = neighbor.occlusion[RIGHT]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z + 1}]; ok {
					occFactor[FRONT] = neighbor.occlusion[RIGHT]
				}

				sides++
				appendChunkFace(&vertexBuffers[RIGHT], &indexBuffers[RIGHT], &occBuffers[RIGHT], occFactor, x, y, z, RIGHT)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X, pos.Y + 1, pos.Z}]; !ok {
			skip := false
			if pos.Y == ChunkBase-1 && chunks[TOP] != nil {
				if _, ok := chunks[TOP].data[BlockCoord{pos.X, 0, pos.Z}]; ok {
					skip = true
				}
			}

			if !skip {
				occFactor := [6]float64{
					chunk.data[pos].occlusion[TOP],
					chunk.data[pos].occlusion[TOP],
					chunk.data[pos].occlusion[TOP],
					chunk.data[pos].occlusion[TOP],
					chunk.data[pos].occlusion[TOP],
					chunk.data[pos].occlusion[TOP],
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z + 1}]; ok {
					occFactor[FRONT] = neighbor.occlusion[TOP]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z - 1}]; ok {
					occFactor[BACK] = neighbor.occlusion[TOP]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X - 1, pos.Y, pos.Z}]; ok {
					occFactor[LEFT] = neighbor.occlusion[TOP]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X + 1, pos.Y, pos.Z}]; ok {
					occFactor[RIGHT] = neighbor.occlusion[TOP]
				}

				sides++
				appendChunkFace(&vertexBuffers[TOP], &indexBuffers[TOP], &occBuffers[TOP], occFactor, x, y, z, TOP)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X, pos.Y - 1, pos.Z}]; !ok {
			skip := false
			if pos.Y == 0 && chunks[BOTTOM] != nil {
				if _, ok := chunks[BOTTOM].data[BlockCoord{pos.X, ChunkBase - 1, pos.Z}]; ok {
					skip = true
				}
			}

			if !skip {
				occFactor := [6]float64{
					chunk.data[pos].occlusion[BOTTOM],
					chunk.data[pos].occlusion[BOTTOM],
					chunk.data[pos].occlusion[BOTTOM],
					chunk.data[pos].occlusion[BOTTOM],
					chunk.data[pos].occlusion[BOTTOM],
					chunk.data[pos].occlusion[BOTTOM],
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z + 1}]; ok {
					occFactor[FRONT] = neighbor.occlusion[TOP]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z - 1}]; ok {
					occFactor[BACK] = neighbor.occlusion[TOP]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X - 1, pos.Y, pos.Z}]; ok {
					occFactor[LEFT] = neighbor.occlusion[TOP]
				}
				if neighbor, ok := chunk.data[BlockCoord{pos.X + 1, pos.Y, pos.Z}]; ok {
					occFactor[RIGHT] = neighbor.occlusion[TOP]
				}

				sides++
				appendChunkFace(&vertexBuffers[BOTTOM], &indexBuffers[BOTTOM], &occBuffers[BOTTOM], occFactor, x, y, z, BOTTOM)
			}
		}

		if sides > 0 {
			blk := chunk.data[pos]
			blk.visible = true
			chunk.data[pos] = blk
		}
	}

	rebuildData := RebuildData{
		vertexBuffers: vertexBuffers,
		indexBuffers:  indexBuffers,
		occBuffers:    occBuffers,
		chunk:         chunk,
	}
	rebuildCh <- rebuildData
}

func (chunk *Chunk) SetChunkMesh(rebuildData RebuildData) {
	vertexBuffers := rebuildData.vertexBuffers
	indexBuffers := rebuildData.indexBuffers
	occBuffers := rebuildData.occBuffers

	for t := 0; t < 6; t++ {
		chunk.mesh.numVertices[t] = gl.Sizei(len(vertexBuffers[t]))
		chunk.mesh.numIndices[t] = gl.Sizei(len(indexBuffers[t]))
		if chunk.mesh.numVertices[t] > 0 && chunk.mesh.numIndices[t] > 0 {
			if chunk.mesh.vao[t] == 0 {
				gl.GenVertexArrays(1, &chunk.mesh.vao[t])
				gl.BindVertexArray(chunk.mesh.vao[t])
				gl.EnableVertexAttribArray(0)
				gl.EnableVertexAttribArray(1)
			}

			gl.BindVertexArray(chunk.mesh.vao[t])
			if chunk.mesh.vertexBufferIds[t] > 0 {
				// Refactor this shit!
				sizeFloat := int(unsafe.Sizeof([1]float32{}))
				sizeInt := int(unsafe.Sizeof([1]uint32{}))
				size := gl.Sizeiptr(sizeFloat * len(vertexBuffers[t]))
				gl.BindBuffer(gl.ARRAY_BUFFER, chunk.mesh.vertexBufferIds[t])
				gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
				gl.BufferData(gl.ARRAY_BUFFER, size, gl.Pointer(&vertexBuffers[t][0]), gl.STATIC_DRAW)

				size = gl.Sizeiptr(sizeFloat * len(occBuffers[t]))
				gl.BindBuffer(gl.ARRAY_BUFFER, chunk.mesh.occBufferIds[t])
				gl.VertexAttribPointer(1, 1, gl.FLOAT, gl.FALSE, 0, nil)
				gl.BufferData(gl.ARRAY_BUFFER, size, gl.Pointer(&occBuffers[t][0]), gl.STATIC_DRAW)

				size = gl.Sizeiptr(sizeInt * len(indexBuffers[t]))
				gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, chunk.mesh.indexBufferIds[t])
				gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, size, gl.Pointer(&indexBuffers[t][0]), gl.STATIC_DRAW)
			} else {
				chunk.mesh.vertexBufferIds[t] = createMeshBuffer(&vertexBuffers[t], len(vertexBuffers[t]))
				chunk.mesh.indexBufferIds[t] = createIndexBuffer(&indexBuffers[t], len(indexBuffers[t]))
				chunk.mesh.occBufferIds[t] = createMeshBuffer(&occBuffers[t], len(occBuffers[t]))
			}

			// Vertices
			gl.BindBuffer(gl.ARRAY_BUFFER, chunk.mesh.vertexBufferIds[t])
			gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
			// Occlusion factor
			gl.BindBuffer(gl.ARRAY_BUFFER, chunk.mesh.occBufferIds[t])
			gl.VertexAttribPointer(1, 1, gl.FLOAT, gl.FALSE, 0, nil)
			// Indices
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, chunk.mesh.indexBufferIds[t])
		}
	}

	numVertices := 0
	numIndices := 0
	numFaces := 0
	for t := 0; t < 6; t++ {
		numVertices += int(chunk.mesh.numVertices[t])
		numIndices += int(chunk.mesh.numIndices[t])
		numFaces += int(chunk.mesh.numVertices[t] / 9.0)
	}
	worstCaseFaces := len(chunk.data) * 12
	fmt.Printf("%d vertices, %d indices, %d faces vs %d total, saved %d\n", numVertices, numIndices, numFaces, worstCaseFaces, worstCaseFaces-numFaces)
}

var facePos = [6]vector.Vector3f{
	{0.0, 0.0, 0.0},
	{0.0, 0.0, 0.0 + float64(ChunkBase)},
	{0.0 + float64(ChunkBase), 0.0, 0.0},
	{0.0, 0.0, 0.0},
	{0.0, 0.0, 0.0},
	{0.0, 0.0 + float64(ChunkBase), 0.0},
}

func (chunk *Chunk) RenderChunk(chunkShader *shader.ShaderProgram, cam vector.Vector3f, world *matrix.Matrix, wireframe bool) {
	invModel, _ := matrix.InvertMatrix(world)
	invModel = invModel.Transpose()

	for t := 0; t < 6; t++ {
		if chunk.mesh.numVertices[t] > 0 && chunk.mesh.numIndices[t] > 0 {
			normal := chunkNormals[t]
			normal = matrix.MultiplyVector3f(invModel, normal)
			face := matrix.MultiplyVector3f(world, facePos[t])
			camDir := cam.Sub(face)
			dot := vector.DotProduct(camDir, normal)

			if dot > 0.0 {
				mouseHit := 0
				if chunk.MouseHit {
					mouseHit = 1
				}
				chunkShader.SetUniformVector3f("normal", normal)
				chunkShader.SetUniformInt("mouseHit", mouseHit)
				chunk.renderMeshBuffer(t, wireframe)
			}
		}
	}
}

func (chunk *Chunk) renderMeshBuffer(side int, wireframe bool) {
	gl.BindVertexArray(chunk.mesh.vao[side])

	if wireframe {
		gl.DrawElements(gl.LINES, chunk.mesh.numIndices[side], gl.UNSIGNED_INT, nil)
	} else {
		gl.DrawElements(gl.TRIANGLES, chunk.mesh.numIndices[side], gl.UNSIGNED_INT, nil)
	}

	gl.BindVertexArray(0)
}
