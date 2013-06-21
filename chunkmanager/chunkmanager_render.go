package chunkmanager

import (
	"dwelling/camera"
	"dwelling/math/matrix"
	"dwelling/shader"
	gl "github.com/chsc/gogl/gl33"
)

var worldVAO gl.Uint = 0
var chunkShader *shader.ShaderProgram

func setUpRenderer() error {
	var err error
	chunkShader, err = shader.LoadShaderProgram("chunk")
	if err != nil {
		return err
	}

	gl.GenVertexArrays(1, &worldVAO)
	gl.BindVertexArray(worldVAO)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	return nil
}

func setRendererData() {
	gl.BindVertexArray(worldVAO)
}

func Render(cam *camera.Camera) {
	setRendererData()
	chunkShader.Use()
	chunkShader.SetUniformMatrix("pv", cam.PVMatrix)

	for pos, chnk := range renderChunks {
		posx := float64(pos.X * ChunkBase)
		posy := float64(pos.Y * ChunkBase)
		posz := float64(pos.Z * ChunkBase)

		modelMatrix := matrix.NewIdentityMatrix()
		modelMatrix.Translate(posx, posy, posz)
		chunkShader.SetUniformMatrix("model", modelMatrix)

		chnk.RenderChunk(chunkShader, cam.CullPos, modelMatrix, false)
	}

	if debugMode {
		for pos, chnk := range visibleChunks {
			posx := float64(pos.X * ChunkBase)
			posy := float64(pos.Y * ChunkBase)
			posz := float64(pos.Z * ChunkBase)

			modelMatrix := matrix.NewIdentityMatrix()
			modelMatrix.Translate(posx, posy, posz)
			chunkShader.SetUniformMatrix("model", modelMatrix)

			chnk.RenderChunk(chunkShader, cam.CullPos, modelMatrix, true)
		}
	}
}
