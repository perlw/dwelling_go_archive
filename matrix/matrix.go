package matrix

import "math"
import "fmt"

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

func NewTestMatrix() *Matrix {
	return &Matrix{
		values: [...]float32{
			0, 1, 2, 3,
			4, 5, 6, 7,
			8, 9, 10, 11,
			12, 13, 14, 15,
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

func (m *Matrix) Translate(x, y, z float32) {
	transMatrix := NewIdentityMatrix()

	transMatrix.values[3] = x
	transMatrix.values[7] = y
	transMatrix.values[11] = z

	m.Multiply(transMatrix)
}

func (m *Matrix) Multiply(in *Matrix) {
	values := [16]float32{}

	for y := 0; y < 4; y++ {
		j := y * 4
		row := [...]float32{m.values[j], m.values[j+1], m.values[j+2], m.values[j+3]}

		for x := 0; x < 4; x++ {
			i := (y * 4) + x

			col := [...]float32{in.values[x], in.values[4+x], in.values[8+x], in.values[12+x]}

			values[i] = (row[0] * col[0]) + (row[1] * col[1]) + (row[2] * col[2]) + (row[3] * col[3])
		}

	}

	fmt.Println(values)
}
