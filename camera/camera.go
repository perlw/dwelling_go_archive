package camera

import (
	"dwelling/math/matrix"
	"dwelling/math/vector"
	gl "github.com/chsc/gogl/gl33"
	"math"
	"unsafe"
)

type Camera struct {
	Pos, Rot               vector.Vector3f
	CullPos                vector.Vector3f
	FrustumPos, FrustumRot vector.Vector3f
	MousePos, MouseDir     vector.Vector3f

	ViewMatrix       *matrix.Matrix
	ProjectionMatrix *matrix.Matrix
	PVMatrix         *matrix.Matrix

	Planes [6]Plane
}

type Plane struct {
	A, B, C, D float64
}

func (cam *Camera) UpdateViewMatrix() {
	view := matrix.NewIdentityMatrix()
	view.RotateX(-cam.Rot.X)
	view.RotateY(-cam.Rot.Y)
	view.RotateZ(-cam.Rot.Z)
	view.Translate(-cam.Pos.X, -cam.Pos.Y, -cam.Pos.Z)

	cam.ViewMatrix = view
}

func (cam *Camera) UpdatePVMatrix() {
	cam.PVMatrix = matrix.MultiplyMatrix(cam.ProjectionMatrix, cam.ViewMatrix)
}

func (cam *Camera) UpdateFrustum() {
	cam.Planes = [6]Plane{}

	// Left
	cam.Planes[0].A = cam.PVMatrix.Values[12] + cam.PVMatrix.Values[0]
	cam.Planes[0].B = cam.PVMatrix.Values[13] + cam.PVMatrix.Values[1]
	cam.Planes[0].C = cam.PVMatrix.Values[14] + cam.PVMatrix.Values[2]
	cam.Planes[0].D = cam.PVMatrix.Values[15] + cam.PVMatrix.Values[3]

	// Right
	cam.Planes[1].A = cam.PVMatrix.Values[12] - cam.PVMatrix.Values[0]
	cam.Planes[1].B = cam.PVMatrix.Values[13] - cam.PVMatrix.Values[1]
	cam.Planes[1].C = cam.PVMatrix.Values[14] - cam.PVMatrix.Values[2]
	cam.Planes[1].D = cam.PVMatrix.Values[15] - cam.PVMatrix.Values[3]

	// Top
	cam.Planes[2].A = cam.PVMatrix.Values[12] - cam.PVMatrix.Values[4]
	cam.Planes[2].B = cam.PVMatrix.Values[13] - cam.PVMatrix.Values[5]
	cam.Planes[2].C = cam.PVMatrix.Values[14] - cam.PVMatrix.Values[6]
	cam.Planes[2].D = cam.PVMatrix.Values[15] - cam.PVMatrix.Values[7]

	// Bottom
	cam.Planes[3].A = cam.PVMatrix.Values[12] + cam.PVMatrix.Values[4]
	cam.Planes[3].B = cam.PVMatrix.Values[13] + cam.PVMatrix.Values[5]
	cam.Planes[3].C = cam.PVMatrix.Values[14] + cam.PVMatrix.Values[6]
	cam.Planes[3].D = cam.PVMatrix.Values[15] + cam.PVMatrix.Values[7]

	// Near
	cam.Planes[4].A = cam.PVMatrix.Values[12] + cam.PVMatrix.Values[8]
	cam.Planes[4].B = cam.PVMatrix.Values[13] + cam.PVMatrix.Values[9]
	cam.Planes[4].C = cam.PVMatrix.Values[14] + cam.PVMatrix.Values[10]
	cam.Planes[4].D = cam.PVMatrix.Values[15] + cam.PVMatrix.Values[11]

	// Far
	cam.Planes[5].A = cam.PVMatrix.Values[12] - cam.PVMatrix.Values[8]
	cam.Planes[5].B = cam.PVMatrix.Values[13] - cam.PVMatrix.Values[9]
	cam.Planes[5].C = cam.PVMatrix.Values[14] - cam.PVMatrix.Values[10]
	cam.Planes[5].D = cam.PVMatrix.Values[15] - cam.PVMatrix.Values[11]

	/*for t := range cam.Planes {
		cam.Planes[t].Normalize()
	}*/

	cam.FrustumPos = cam.Pos
	cam.FrustumRot = cam.Rot
}

func (cam *Camera) CubeInView(origo vector.Vector3f, size float64) int {
	corners := [8]vector.Vector3f{
		{origo.X, origo.Y, origo.Z},
		{origo.X + size, origo.Y, origo.Z},
		{origo.X + size, origo.Y, origo.Z + size},
		{origo.X, origo.Y, origo.Z + size},
		{origo.X, origo.Y + size, origo.Z + size},
		{origo.X + size, origo.Y + size, origo.Z + size},
		{origo.X + size, origo.Y + size, origo.Z},
		{origo.X, origo.Y + size, origo.Z},
	}

	status := 0 // 0 inside, 1 partly, 2 outside
	for t := range cam.Planes {
		in, out := 0, 0

		for u := range corners {
			if cam.Planes[t].ClassifyPoint(corners[u]) < 0.0 {
				out++
			} else {
				in++
			}
		}

		if in == 0 {
			return 2
		} else if out > 0 {
			status = 1
		}
	}

	return status
}

func (plane *Plane) Normalize() {
	magnitude := math.Sqrt((plane.A * plane.A) + (plane.B * plane.B) + (plane.C * plane.C) + (plane.D * plane.D))

	plane.A /= magnitude
	plane.B /= magnitude
	plane.C /= magnitude
	plane.D /= magnitude
}

func (plane *Plane) ClassifyPoint(v vector.Vector3f) float64 {
	return (plane.A * v.X) + (plane.B * v.Y) + (plane.C * v.Z) + plane.D
}

func CreateFrustumMesh(cam *Camera) gl.Uint {
	var buffer gl.Uint
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

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*len(vertices)), gl.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return buffer
}

func RenderFrustumMesh(cam *Camera, meshBuffer gl.Uint) {
	gl.BindBuffer(gl.ARRAY_BUFFER, meshBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.DrawArrays(gl.LINES, 0, 24)
}

func CreateMouseMesh(cam *Camera) gl.Uint {
	var buffer gl.Uint
	sizeFloat := int(unsafe.Sizeof([1]float32{}))

	startPos := cam.MousePos
	endPos := startPos.Add(cam.MouseDir.MulScalar(1000.0))
	vertices := [...]float32{
		float32(startPos.X), float32(startPos.Y), float32(startPos.Z),
		float32(endPos.X), float32(endPos.Y), float32(endPos.Z),
	}

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*len(vertices)), gl.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return buffer
}

func RenderMouseMesh(meshBuffer gl.Uint) {
	gl.BindBuffer(gl.ARRAY_BUFFER, meshBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.DrawArrays(gl.LINES, 0, 2)
}

func CreateGridMesh() gl.Uint {
	var buffer gl.Uint
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
		vertices = append(vertices, -g, i, -1.0)
		vertices = append(vertices, g, i, -1.0)

		vertices = append(vertices, i, -g, -1.0)
		vertices = append(vertices, i, g, -1.0)

		// YZ
		vertices = append(vertices, 0.0, -g, i)
		vertices = append(vertices, 0.0, g, i)

		vertices = append(vertices, 0.0, i, -g)
		vertices = append(vertices, 0.0, i, g)
	}

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(sizeFloat*len(vertices)), gl.Pointer(&vertices[0]), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return buffer
}

func RenderGridMesh(meshBuffer gl.Uint) {
	gl.BindBuffer(gl.ARRAY_BUFFER, meshBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)
	gl.DrawArrays(gl.LINES, 0, 3072)
}
