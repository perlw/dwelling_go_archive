package camera

import (
	"dwelling/math/matrix"
	gl "github.com/chsc/gogl/gl33"
	"math"
	"unsafe"
)

type Camera struct {
	X, Y, Z             float64
	Rx, Ry, Rz          float64
	CullX, CullY, CullZ float64
	Fx, Fy, Fz          float64
	Frx, Fry, Frz       float64

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
	view.RotateX(-cam.Rx)
	view.RotateY(-cam.Ry)
	view.RotateZ(-cam.Rz)
	view.Translate(-cam.X, -cam.Y, -cam.Z)

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

	cam.Fx = cam.X
	cam.Fy = cam.Y
	cam.Fz = cam.Z
	cam.Frx = cam.Rx
	cam.Fry = cam.Ry
	cam.Frz = cam.Rz
}

func (cam *Camera) CubeInView(origo [3]float64, size float64) int {
	corners := [8][3]float64{
		{origo[0], origo[1], origo[2]},
		{origo[0] + size, origo[1], origo[2]},
		{origo[0] + size, origo[1], origo[2] + size},
		{origo[0], origo[1], origo[2] + size},
		{origo[0], origo[1] + size, origo[2] + size},
		{origo[0] + size, origo[1] + size, origo[2] + size},
		{origo[0] + size, origo[1] + size, origo[2]},
		{origo[0], origo[1] + size, origo[2]},
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

func (plane *Plane) ClassifyPoint(vector [3]float64) float64 {
	return (plane.A * vector[0]) + (plane.B * vector[1]) + (plane.C * vector[2]) + plane.D
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
