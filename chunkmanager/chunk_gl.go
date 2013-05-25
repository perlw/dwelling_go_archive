package chunkmanager

import (
	"dwelling/math/matrix"
	"dwelling/math/vector"
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
	vertexBufferIds [6]gl.Uint
	numVertices     [6]gl.Sizei
}

func appendChunkFace(faceBuffer *[]float32, x, y, z float32, face int) {
	switch face {
	case FRONT:
		(*faceBuffer) = append((*faceBuffer),
			x, y, z,
			x+1.0, y, z,
			x+1.0, y+1.0, z,
			x+1.0, y+1.0, z,
			x, y+1.0, z,
			x, y, z,
		)

	case BACK:
		(*faceBuffer) = append((*faceBuffer),
			x+1.0, y+1.0, z-1.0,
			x+1.0, y, z-1.0,
			x, y, z-1.0,
			x, y, z-1.0,
			x, y+1.0, z-1.0,
			x+1.0, y+1.0, z-1.0,
		)

	case LEFT:
		(*faceBuffer) = append((*faceBuffer),
			x, y, z-1.0,
			x, y, z,
			x, y+1.0, z,
			x, y+1.0, z,
			x, y+1.0, z-1.0,
			x, y, z-1.0,
		)

	case RIGHT:
		(*faceBuffer) = append((*faceBuffer),
			x+1.0, y+1.0, z,
			x+1.0, y, z,
			x+1.0, y, z-1.0,
			x+1.0, y, z-1.0,
			x+1.0, y+1.0, z-1.0,
			x+1.0, y+1.0, z,
		)

	case TOP:
		(*faceBuffer) = append((*faceBuffer),
			x+1.0, y+1.0, z,
			x+1.0, y+1.0, z-1.0,
			x, y+1.0, z-1.0,
			x, y+1.0, z-1.0,
			x, y+1.0, z,
			x+1.0, y+1.0, z,
		)

	case BOTTOM:
		(*faceBuffer) = append((*faceBuffer),
			x, y, z-1.0,
			x+1.0, y, z-1.0,
			x+1.0, y, z,
			x+1.0, y, z,
			x, y, z,
			x, y, z-1.0,
		)

	default:
	}
}

func createMeshBuffer(faceBuffer *[]float32, size int) gl.Uint {
	var buffer gl.Uint
	sizeFloat := int(unsafe.Sizeof([1]float32{}))
	bufferPtr := (*faceBuffer)

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*size), gl.Pointer(&bufferPtr[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return buffer
}

func (chunk *Chunk) UpdateChunkMesh(chunkPos ChunkCoord) {
	chunks := GetChunksAroundChunk(chunkPos)

	vertexBuffers := [6][]float32{}
	for pos := range chunk.data {
		x := float32(pos.X)
		y := float32(pos.Y)
		z := float32(pos.Z)

		if _, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z + 1}]; !ok {
			skip := false
			if pos.Z == CHUNK_BASE-1 && chunks[FRONT] != nil {
				if _, ok := chunks[FRONT].data[BlockCoord{pos.X, pos.Y, 0}]; ok {
					skip = true
				}
			}

			if !skip {
				appendChunkFace(&vertexBuffers[FRONT], x, y, z, FRONT)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X, pos.Y, pos.Z - 1}]; !ok {
			skip := false
			if pos.Z == 0 && chunks[BACK] != nil {
				if _, ok := chunks[BACK].data[BlockCoord{pos.X, pos.Y, CHUNK_BASE - 1}]; ok {
					skip = true
				}
			}

			if !skip {
				appendChunkFace(&vertexBuffers[BACK], x, y, z, BACK)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X - 1, pos.Y, pos.Z}]; !ok {
			skip := false
			if pos.X == 0 && chunks[LEFT] != nil {
				if _, ok := chunks[LEFT].data[BlockCoord{CHUNK_BASE - 1, pos.Y, pos.Z}]; ok {
					skip = true
				}
			}

			if !skip {
				appendChunkFace(&vertexBuffers[LEFT], x, y, z, LEFT)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X + 1, pos.Y, pos.Z}]; !ok {
			skip := false
			if pos.X == CHUNK_BASE-1 && chunks[RIGHT] != nil {
				if _, ok := chunks[RIGHT].data[BlockCoord{0, pos.Y, pos.Z}]; ok {
					skip = true
				}
			}

			if !skip {
				appendChunkFace(&vertexBuffers[RIGHT], x, y, z, RIGHT)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X, pos.Y + 1, pos.Z}]; !ok {
			skip := false
			if pos.Y == CHUNK_BASE-1 && chunks[TOP] != nil {
				if _, ok := chunks[TOP].data[BlockCoord{pos.X, 0, pos.Z}]; ok {
					skip = true
				}
			}

			if !skip {
				appendChunkFace(&vertexBuffers[TOP], x, y, z, TOP)
			}
		}
		if _, ok := chunk.data[BlockCoord{pos.X, pos.Y - 1, pos.Z}]; !ok {
			skip := false
			if pos.Y == 0 && chunks[BOTTOM] != nil {
				if _, ok := chunks[BOTTOM].data[BlockCoord{pos.X, CHUNK_BASE - 1, pos.Z}]; ok {
					skip = true
				}
			}

			if !skip {
				appendChunkFace(&vertexBuffers[BOTTOM], x, y, z, BOTTOM)
			}
		}
	}

	for t := 0; t < 6; t++ {
		chunk.mesh.numVertices[t] = gl.Sizei(len(vertexBuffers[t]))
		if chunk.mesh.numVertices[t] > 0 {
			if chunk.mesh.vertexBufferIds[t] > 0 {
				gl.DeleteBuffers(1, &chunk.mesh.vertexBufferIds[t])
			}
			chunk.mesh.vertexBufferIds[t] = createMeshBuffer(&vertexBuffers[t], len(vertexBuffers[t]))
		}
	}

	numFaces := 0
	for t := 0; t < 6; t++ {
		numFaces += int(chunk.mesh.numVertices[t] / 9.0)
	}
	worstCaseFaces := len(chunk.data) * 12
	fmt.Printf("%d faces vs %d total, saved %d\n", numFaces, worstCaseFaces, worstCaseFaces-numFaces)

	chunk.IsSetup = true
}

var facePos = [6]vector.Vector3f{
	{0.0, 0.0, 0.0},
	{0.0, 0.0, 0.0 + float64(CHUNK_BASE)},
	{0.0 + float64(CHUNK_BASE), 0.0, 0.0},
	{0.0, 0.0, 0.0},
	{0.0, 0.0, 0.0},
	{0.0, 0.0 + float64(CHUNK_BASE), 0.0},
}

func (chunk *Chunk) RenderChunk(normalId gl.Int, cam vector.Vector3f, world *matrix.Matrix, wireframe bool, chunkPos vector.Vector3f) {
	invModel, _ := matrix.InvertMatrix(world)
	invModel = invModel.Transpose()
	for t := 0; t < 6; t++ {
		if chunk.mesh.numVertices[t] > 0 {
			normal := chunkNormals[t]
			normal = matrix.MultiplyVector3f(invModel, normal)
			face := matrix.MultiplyVector3f(world, facePos[t])
			camDir := cam.Sub(face)
			dot := vector.DotProduct(camDir, normal)

			if dot > 0.0 {
				chunk.renderMeshBuffer(t, normalId, wireframe)
			}
		}
	}
}

func (chunk *Chunk) renderMeshBuffer(side int, normalId gl.Int, wireframe bool) {
	normal := chunkNormals[side].ToGL()
	gl.Uniform3fv(normalId, 1, &normal[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, chunk.mesh.vertexBufferIds[side])
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	if wireframe {
		gl.DrawArrays(gl.LINES, 0, chunk.mesh.numVertices[side])
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, chunk.mesh.numVertices[side])

	}
}
