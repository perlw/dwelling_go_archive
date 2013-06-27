package camera

import (
	"bedrock/shader"
	"bedrock/math/matrix"
	"bedrock/math/vector"
	gl "github.com/chsc/gogl/gl33"
	"unsafe"
)

type DebugData struct {
	debugVao           gl.Uint
	debugShader        *shader.ShaderProgram
	debugFrustumBuffer gl.Uint
	debugGridBuffer    gl.Uint
	debugMouseBuffer   gl.Uint
}

func (cam *Camera) setUpDebugRenderer() error {
	var err error
	cam.debugData.debugShader, err = shader.LoadShaderProgram("debug", nil)
	if err != nil {
		return err
	}

	gl.GenVertexArrays(1, &cam.debugData.debugVao)
	gl.BindVertexArray(cam.debugData.debugVao)
	gl.EnableVertexAttribArray(0)

	cam.createFrustumMesh()
	cam.createGridMesh()

	return nil
}

func (cam *Camera) RenderDebugMeshes() {
	gl.BindVertexArray(cam.debugData.debugVao)
	cam.debugData.debugShader.Use()
	cam.debugData.debugShader.SetUniformMatrix("pv", cam.PVMatrix)

	// Render frustum
	cam.debugData.debugShader.SetUniformVector3f("flatColor", vector.Vector3f{X: 0.5, Y: 0.5, Z: 1.0})
	modelMatrix := matrix.NewIdentityMatrix()
	modelMatrix.TranslateVector(cam.FrustumPos)
	modelMatrix.RotateY(cam.FrustumRot.Y)
	modelMatrix.RotateX(cam.FrustumRot.X)
	cam.debugData.debugShader.SetUniformMatrix("model", modelMatrix)
	cam.renderFrustumMesh()
	// Render frustum

	// Render Mouse ray
	cam.debugData.debugShader.SetUniformVector3f("flatColor", vector.Vector3f{X: 1.0, Y: 0.5, Z: 0.5})
	cam.createMouseMesh()
	modelMatrix = matrix.NewIdentityMatrix()
	cam.debugData.debugShader.SetUniformMatrix("model", modelMatrix)
	cam.renderMouseMesh()
	// Render Mouse ray

	// Render Grid
	cam.debugData.debugShader.SetUniformVector3f("flatColor", vector.Vector3f{X: 0.0, Y: 0.0, Z: 0.0})
	modelMatrix = matrix.NewIdentityMatrix()
	cam.debugData.debugShader.SetUniformMatrix("model", modelMatrix)
	cam.renderGridMesh()
	// Render Mouse ray
}

func (cam *Camera) createFrustumMesh() {
	sizeFloat := int(unsafe.Sizeof([1]float32{}))

	proj := cam.ProjectionMatrix.Values
	near := proj[11] / (proj[10] - 1.0)
	far := 100.0 //proj[11] / (1.0 + proj[10])
	nLeft := float32(near * (proj[2] - 1.0) / proj[0])
	nRight := float32(near * (1.0 + proj[2]) / proj[0])
	nTop := float32(near * (1.0 + proj[6]) / proj[5])
	nBottom := float32(near * (proj[6] - 1.0) / proj[5])
	fLeft := float32(far * (proj[2] - 1.0) / proj[0])
	fRight := float32(far * (1.0 + proj[2]) / proj[0])
	fTop := float32(far * (1.0 + proj[6]) / proj[5])
	fBottom := float32(far * (proj[6] - 1.0) / proj[5])

	vertices := [...]float32{
		0.0, 0.0, 0.0,
		fLeft, fBottom, float32(-far),

		0.0, 0.0, 0.0,
		fRight, fBottom, float32(-far),

		0.0, 0.0, 0.0,
		fRight, fTop, float32(-far),

		0.0, 0.0, 0.0,
		fLeft, fTop, float32(-far),

		fLeft, fBottom, float32(-far),
		fRight, fBottom, float32(-far),

		fRight, fTop, float32(-far),
		fLeft, fTop, float32(-far),

		fRight, fTop, float32(-far),
		fRight, fBottom, float32(-far),

		fLeft, fTop, float32(-far),
		fLeft, fBottom, float32(-far),

		nLeft, nBottom, float32(-near),
		nRight, nBottom, float32(-near),

		nRight, nTop, float32(-near),
		nLeft, nTop, float32(-near),

		nLeft, nTop, float32(-near),
		nLeft, nBottom, float32(-near),

		nRight, nTop, float32(-near),
		nRight, nBottom, float32(-near),
	}

	gl.GenBuffers(1, &cam.debugData.debugFrustumBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugFrustumBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*len(vertices)), gl.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugFrustumBuffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (cam *Camera) renderFrustumMesh() {
	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugFrustumBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.DrawArrays(gl.LINES, 0, 24)
}

func (cam *Camera) createMouseMesh() {
	sizeFloat := int(unsafe.Sizeof([1]float32{}))

	startPos := cam.MousePos
	endPos := startPos.Add(cam.MouseDir.MulScalar(1000.0))
	vertices := [...]float32{
		float32(startPos.X), float32(startPos.Y), float32(startPos.Z),
		float32(endPos.X), float32(endPos.Y), float32(endPos.Z),
	}

	if cam.debugData.debugMouseBuffer == 0 {
		gl.GenBuffers(1, &cam.debugData.debugMouseBuffer)

	}
	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugMouseBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*len(vertices)), gl.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugMouseBuffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (cam *Camera) renderMouseMesh() {
	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugMouseBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.DrawArrays(gl.LINES, 0, 2)
}

func (cam *Camera) createGridMesh() {
	sizeFloat := int(unsafe.Sizeof([1]float32{}))

	var vertices []float32
	gridSize := 128
	for t := -gridSize; t < gridSize; t++ {
		i := float32(t)
		g := float32(gridSize)

		// XZ
		vertices = append(vertices, -g, 0.0, i)
		vertices = append(vertices, g, 0.0, i)

		vertices = append(vertices, i, 0.0, -g)
		vertices = append(vertices, i, 0.0, g)

		// XY
		vertices = append(vertices, -g, i, 0.0)
		vertices = append(vertices, g, i, 0.0)

		vertices = append(vertices, i, -g, 0.0)
		vertices = append(vertices, i, g, 0.0)

		// YZ
		vertices = append(vertices, 0.0, -g, i)
		vertices = append(vertices, 0.0, g, i)

		vertices = append(vertices, 0.0, i, -g)
		vertices = append(vertices, 0.0, i, g)
	}

	gl.GenBuffers(1, &cam.debugData.debugGridBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugGridBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*len(vertices)), gl.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugGridBuffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (cam *Camera) renderGridMesh() {
	gl.BindBuffer(gl.ARRAY_BUFFER, cam.debugData.debugGridBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.DrawArrays(gl.LINES, 0, 3072)
}
