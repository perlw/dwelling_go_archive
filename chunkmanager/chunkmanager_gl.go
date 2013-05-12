package chunkmanager

import (
	"dwelling/camera"
	"dwelling/math/matrix"
	gl "github.com/chsc/gogl/gl33"
)

func Render(program gl.Uint, cam *camera.Camera) {
	model := gl.GLString("model")
	modelId := gl.GetUniformLocation(program, model)
	gl.GLStringFree(model)
	normal := gl.GLString("normal")
	normalId := gl.GetUniformLocation(program, normal)
	gl.GLStringFree(normal)

	for pos, chnk := range renderChunks {
		posx := float64(pos.X * CHUNK_BASE)
		posy := float64(pos.Y * CHUNK_BASE)
		posz := float64(pos.Z * CHUNK_BASE)

		modelMatrix := matrix.NewIdentityMatrix()
		modelMatrix.Translate(posx, posy, posz)
		glModelMatrix := modelMatrix.ToGL()
		gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])

		chnk.RenderChunk(normalId, cam.CullPos, modelMatrix, false)
	}

	if debugMode {
		for pos, chnk := range visibleChunks {
			posx := float64(pos.X * CHUNK_BASE)
			posy := float64(pos.Y * CHUNK_BASE)
			posz := float64(pos.Z * CHUNK_BASE)

			modelMatrix := matrix.NewIdentityMatrix()
			modelMatrix.Translate(posx, posy, posz)
			glModelMatrix := modelMatrix.ToGL()
			gl.UniformMatrix4fv(modelId, 1, gl.FALSE, &glModelMatrix[0])

			chnk.RenderChunk(normalId, cam.CullPos, modelMatrix, true)
		}
	}
}
