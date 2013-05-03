package matrix

import "math"

type Matrix struct {
	values [16]float32
}

func NewIdentityMatrix() *Matrix {
	return &Matrix{
		values: [...]float32{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}
}

func NewPerspectiveMatrix(fov, ratio, nearZ, farZ float32) *Matrix {
	fovRadii := float64((fov / 2.0) * (math.Pi / 180.0))
	f := float32(1.0 / math.Tan(fovRadii))
	zDiff := farZ - nearZ

	matrix := &Matrix{
		values: [...]float32{
			f / ratio, 0, 0, 0,
			0, f, 0, 0,
			0, 0, -(farZ + nearZ) / zDiff, (2.0 * farZ * nearZ) / zDiff,
			0, 0, -1.0, 0.0,
		},
	}

	return matrix
}

func MultiplyMatrix(matrixA, matrixB *Matrix) *Matrix {
	values := [16]float32{}

	for y := 0; y < 4; y++ {
		j := y * 4
		row := [...]float32{matrixA.values[j], matrixA.values[j+1], matrixA.values[j+2], matrixA.values[j+3]}

		for x := 0; x < 4; x++ {
			i := (y * 4) + x

			col := [...]float32{matrixB.values[x], matrixB.values[4+x], matrixB.values[8+x], matrixB.values[12+x]}

			values[i] = (row[0] * col[0]) + (row[1] * col[1]) + (row[2] * col[2]) + (row[3] * col[3])
		}

	}

	return &Matrix{values: values}
}

func (m *Matrix) Translate(x, y, z float32) {
	transMatrix := NewIdentityMatrix()

	transMatrix.values[3] = x
	transMatrix.values[7] = y
	transMatrix.values[11] = z

	m = MultiplyMatrix(m, transMatrix)
}

func (m *Matrix) RotateX(rot float32) {
	rotMatrix := &Matrix{
		values: [...]float32{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}

	m = MultiplyMatrix(m, rotMatrix)
}

func (m *Matrix) RotateY(rot float32) {
	rotMatrix := &Matrix{
		values: [...]float32{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}

	m = MultiplyMatrix(m, rotMatrix)
}

func (m *Matrix) RotateZ(rot float32) {
	rotMatrix := &Matrix{
		values: [...]float32{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}

	m = MultiplyMatrix(m, rotMatrix)
}
