package chunkmanager

import (
	"dwelling/camera"
	"dwelling/math/matrix"
	"dwelling/math/vector"
	gl "github.com/chsc/gogl/gl33"
)

func Render(program gl.Uint, cam *camera.Camera) {
	model := gl.GLString("model")
	modelId := gl.GetUniformLocation(program, model)
	gl.GLStringFree(model)
	normal := gl.GLString("normal")
	normalId := gl.GetUniformLocation(program, normal)
	gl.GLStringFree(normal)
	mouseHit := gl.GLString("mouseHit")
	mouseHitId := gl.GetUniformLocation(program, mouseHit)
	gl.GLStringFree(mouseHit)

	for pos, chnk := range renderChunks {
		posx := float64(pos.X * ChunkBase)
		posy := float64(pos.Y * ChunkBase)
		posz := float64(pos.Z * ChunkBase)

		modelMatrix := matrix.NewIdentityMatrix()
		modelMatrix.Translate(posx, posy, posz)
		glModelMatrix := modelMatrix.ToGL()
		gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])

		chnk.RenderChunk(normalId, mouseHitId, cam.CullPos, modelMatrix, false, vector.Vector3f{posx, posy, posz})
	}

	if debugMode {
		for pos, chnk := range visibleChunks {
			posx := float64(pos.X * ChunkBase)
			posy := float64(pos.Y * ChunkBase)
			posz := float64(pos.Z * ChunkBase)

			modelMatrix := matrix.NewIdentityMatrix()
			modelMatrix.Translate(posx, posy, posz)
			glModelMatrix := modelMatrix.ToGL()
			gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])

			chnk.RenderChunk(normalId, mouseHitId, cam.CullPos, modelMatrix, true, vector.Vector3f{posx, posy, posz})
		}
	}
}
