package chunkmanager

import (
	"bedrock/shader"
	"dwelling/camera"
	"dwelling/math/matrix"
)

var chunkShader *shader.ShaderProgram

func setUpRenderer() error {
	var err error
	chunkShader, err = shader.LoadShaderProgram("chunk", []shader.AttribLocation{
		{
			Position: 0,
			Location: "vertexPos",
		},
		{
			Position: 1,
			Location: "occFactor",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func Render(cam *camera.Camera) {
	chunkShader.Use()
	chunkShader.SetUniformMatrix("pv", cam.PVMatrix)
	if debugMode {
		chunkShader.SetUniformInt("onlyOccFac", 1)
	} else {
		chunkShader.SetUniformInt("onlyOccFac", 0)
	}

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
