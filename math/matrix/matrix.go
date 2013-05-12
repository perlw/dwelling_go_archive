package matrix

import (
	"dwelling/math/vector"
	"errors"
	"math"
	"strconv"
)

type Matrix struct {
	Values [16]float64
}

func NewIdentityMatrix() *Matrix {
	return &Matrix{
		Values: [...]float64{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}
}

func NewPerspectiveMatrix(fov, ratio, nearZ, farZ float64) *Matrix {
	fovRadii := float64((fov / 2.0) * (math.Pi / 180.0))
	f := 1.0 / math.Tan(fovRadii)
	zDiff := farZ - nearZ

	matrix := &Matrix{
		Values: [...]float64{
			f / ratio, 0, 0, 0,
			0, f, 0, 0,
			0, 0, -(farZ + nearZ) / zDiff, -(2.0 * farZ * nearZ) / zDiff,
			0, 0, -1, 0,
		},
	}

	return matrix
}

func MultiplyMatrix(matrixA, matrixB *Matrix) *Matrix {
	values := [16]float64{}

	for y := 0; y < 4; y++ {
		j := y * 4
		row := [...]float64{matrixA.Values[j], matrixA.Values[j+1], matrixA.Values[j+2], matrixA.Values[j+3]}

		for x := 0; x < 4; x++ {
			i := (y * 4) + x

			col := [...]float64{matrixB.Values[x], matrixB.Values[4+x], matrixB.Values[8+x], matrixB.Values[12+x]}

			values[i] = (row[0] * col[0]) + (row[1] * col[1]) + (row[2] * col[2]) + (row[3] * col[3])
		}

	}

	return &Matrix{Values: values}
}

func MultiplyVector3f(m *Matrix, v vector.Vector3f) vector.Vector3f {
	result := vector.Vector3f{}

	for t := 0; t < 4; t++ {
		if row, err := m.RowToVector3f(t); err == nil {
			result = result.Add(row.Mul(v))
		}
	}

	return result
}

func MultiplyVector4f(m *Matrix, v vector.Vector4f) vector.Vector4f {
	result := vector.Vector4f{}

	for t := 0; t < 4; t++ {
		if row, err := m.RowToVector4f(t); err == nil {
			result = result.Add(row.Mul(v))
		}
	}

	return result
}

func InvertMatrix(matrix *Matrix) (*Matrix, error) {
	values := [16]float64{}
	m := matrix.Values

	values[0] = m[5]*m[10]*m[15] -
		m[5]*m[11]*m[14] -
		m[9]*m[6]*m[15] +
		m[9]*m[7]*m[14] +
		m[13]*m[6]*m[11] -
		m[13]*m[7]*m[10]

	values[4] = -m[4]*m[10]*m[15] +
		m[4]*m[11]*m[14] +
		m[8]*m[6]*m[15] -
		m[8]*m[7]*m[14] -
		m[12]*m[6]*m[11] +
		m[12]*m[7]*m[10]

	values[8] = m[4]*m[9]*m[15] -
		m[4]*m[11]*m[13] -
		m[8]*m[5]*m[15] +
		m[8]*m[7]*m[13] +
		m[12]*m[5]*m[11] -
		m[12]*m[7]*m[9]

	values[12] = -m[4]*m[9]*m[14] +
		m[4]*m[10]*m[13] +
		m[8]*m[5]*m[14] -
		m[8]*m[6]*m[13] -
		m[12]*m[5]*m[10] +
		m[12]*m[6]*m[9]

	values[1] = -m[1]*m[10]*m[15] +
		m[1]*m[11]*m[14] +
		m[9]*m[2]*m[15] -
		m[9]*m[3]*m[14] -
		m[13]*m[2]*m[11] +
		m[13]*m[3]*m[10]

	values[5] = m[0]*m[10]*m[15] -
		m[0]*m[11]*m[14] -
		m[8]*m[2]*m[15] +
		m[8]*m[3]*m[14] +
		m[12]*m[2]*m[11] -
		m[12]*m[3]*m[10]

	values[9] = -m[0]*m[9]*m[15] +
		m[0]*m[11]*m[13] +
		m[8]*m[1]*m[15] -
		m[8]*m[3]*m[13] -
		m[12]*m[1]*m[11] +
		m[12]*m[3]*m[9]

	values[13] = m[0]*m[9]*m[14] -
		m[0]*m[10]*m[13] -
		m[8]*m[1]*m[14] +
		m[8]*m[2]*m[13] +
		m[12]*m[1]*m[10] -
		m[12]*m[2]*m[9]

	values[2] = m[1]*m[6]*m[15] -
		m[1]*m[7]*m[14] -
		m[5]*m[2]*m[15] +
		m[5]*m[3]*m[14] +
		m[13]*m[2]*m[7] -
		m[13]*m[3]*m[6]

	values[6] = -m[0]*m[6]*m[15] +
		m[0]*m[7]*m[14] +
		m[4]*m[2]*m[15] -
		m[4]*m[3]*m[14] -
		m[12]*m[2]*m[7] +
		m[12]*m[3]*m[6]

	values[10] = m[0]*m[5]*m[15] -
		m[0]*m[7]*m[13] -
		m[4]*m[1]*m[15] +
		m[4]*m[3]*m[13] +
		m[12]*m[1]*m[7] -
		m[12]*m[3]*m[5]

	values[14] = -m[0]*m[5]*m[14] +
		m[0]*m[6]*m[13] +
		m[4]*m[1]*m[14] -
		m[4]*m[2]*m[13] -
		m[12]*m[1]*m[6] +
		m[12]*m[2]*m[5]

	values[3] = -m[1]*m[6]*m[11] +
		m[1]*m[7]*m[10] +
		m[5]*m[2]*m[11] -
		m[5]*m[3]*m[10] -
		m[9]*m[2]*m[7] +
		m[9]*m[3]*m[6]

	values[7] = m[0]*m[6]*m[11] -
		m[0]*m[7]*m[10] -
		m[4]*m[2]*m[11] +
		m[4]*m[3]*m[10] +
		m[8]*m[2]*m[7] -
		m[8]*m[3]*m[6]

	values[11] = -m[0]*m[5]*m[11] +
		m[0]*m[7]*m[9] +
		m[4]*m[1]*m[11] -
		m[4]*m[3]*m[9] -
		m[8]*m[1]*m[7] +
		m[8]*m[3]*m[5]

	values[15] = m[0]*m[5]*m[10] -
		m[0]*m[6]*m[9] -
		m[4]*m[1]*m[10] +
		m[4]*m[2]*m[9] +
		m[8]*m[1]*m[6] -
		m[8]*m[2]*m[5]

	det := m[0]*values[0] + m[1]*values[4] + m[2]*values[8] + m[3]*values[12]
	if det == 0 {
		return NewIdentityMatrix(), errors.New("matrix: not possible to invert")
	}

	det = 1.0 / det
	for t := range values {
		values[t] = values[t] * det
	}

	return &Matrix{Values: values}, nil
}

func (m *Matrix) Translate(x, y, z float64) {
	transMatrix := NewIdentityMatrix()

	transMatrix.Values[3] = x
	transMatrix.Values[7] = y
	transMatrix.Values[11] = z

	m.Values = MultiplyMatrix(m, transMatrix).Values
}

func (m *Matrix) TranslateVector(v vector.Vector3f) {
	transMatrix := NewIdentityMatrix()

	transMatrix.Values[3] = v.X
	transMatrix.Values[7] = v.Y
	transMatrix.Values[11] = v.Z

	m.Values = MultiplyMatrix(m, transMatrix).Values
}

func (m *Matrix) RotateX(rot float64) {
	radii := rot * (math.Pi / 180.0)
	cos := math.Cos(radii)
	sin := math.Sin(radii)

	rotMatrix := &Matrix{
		Values: [...]float64{
			1, 0, 0, 0,
			0, cos, sin, 0,
			0, -sin, cos, 0,
			0, 0, 0, 1,
		},
	}

	m.Values = MultiplyMatrix(m, rotMatrix).Values
}

func (m *Matrix) RotateY(rot float64) {
	radii := rot * (math.Pi / 180.0)
	cos := math.Cos(radii)
	sin := math.Sin(radii)

	rotMatrix := &Matrix{
		Values: [...]float64{
			cos, 0, -sin, 0,
			0, 1, 0, 0,
			sin, 0, cos, 0,
			0, 0, 0, 1,
		},
	}

	m.Values = MultiplyMatrix(m, rotMatrix).Values
}

func (m *Matrix) RotateZ(rot float64) {
	radii := rot * (math.Pi / 180.0)
	cos := math.Cos(radii)
	sin := math.Sin(radii)

	rotMatrix := &Matrix{
		Values: [...]float64{
			cos, sin, 0, 0,
			-sin, cos, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}

	m.Values = MultiplyMatrix(m, rotMatrix).Values
}

func (m *Matrix) RowToVector3f(row int) (vector.Vector3f, error) {
	if row < 0 || row > 3 {
		return vector.Vector3f{}, errors.New("matrix: row " + strconv.Itoa(row) + " is out of bounds")
	}

	rowIndex := row * 4
	return vector.Vector3f{m.Values[rowIndex+0], m.Values[rowIndex+1], m.Values[rowIndex+2]}, nil
}

func (m *Matrix) RowToVector4f(row int) (vector.Vector4f, error) {
	if row < 0 || row > 3 {
		return vector.Vector4f{}, errors.New("matrix: row " + strconv.Itoa(row) + " is out of bounds")
	}

	rowIndex := row * 4
	return vector.Vector4f{m.Values[rowIndex+0], m.Values[rowIndex+1], m.Values[rowIndex+2], m.Values[rowIndex+3]}, nil
}
