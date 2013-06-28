package camera

import (
	"bedrock"
	"bedrock/math/matrix"
	"bedrock/math/vector"
	"math"
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

	debugData DebugData
}

type Plane struct {
	A, B, C, D float64
}

func (cam *Camera) Init() error {
	ratio := float64(bedrock.ScreenWidth) / float64(bedrock.ScreenHeight)

	cam.Pos = vector.Vector3f{X: -48.0, Y: 32.0, Z: -48.0}
	cam.Rot = vector.Vector3f{X: 0.0, Y: 135, Z: 0.0}
	cam.ProjectionMatrix = matrix.NewPerspectiveMatrix(53.13, ratio, 1.0, 1000.0)
	cam.FrustumPos = cam.Pos
	cam.FrustumRot = cam.Rot
	cam.CullPos = cam.Pos
	cam.UpdateViewMatrix()
	cam.UpdatePVMatrix()
	cam.UpdateFrustum()

	if err := cam.setUpDebugRenderer(); err != nil {
		return err
	}

	return nil
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
