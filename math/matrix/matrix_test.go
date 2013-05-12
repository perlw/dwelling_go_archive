package matrix

import "testing"

func TestIdentity(t *testing.T) {
	matrix := NewIdentityMatrix()
	testValues := [...]float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	if matrix.Values != testValues {
		t.Errorf("Expected %v, got %v", testValues, matrix.Values)
	}
}

func TestPerspective(t *testing.T) {
	matrix := NewPerspectiveMatrix(75, 640/480, 1.0, 1000.0)
	testValues := [...]float64{
		1.3032253728412058, 0, 0, 0,
		0, 1.3032253728412058, 0, 0,
		0, 0, -1.002002002002002, -2.002002002002002,
		0, 0, -1, 0,
	}

	if matrix.Values != testValues {
		t.Errorf("Expected %v, got %v", testValues, matrix.Values)
	}
}

func TestMultiplyMatrix(t *testing.T) {
	tests := []struct {
		matrixA  [16]float64
		matrixB  [16]float64
		expected [16]float64
	}{
		{
			matrixA: [...]float64{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
			matrixB: [...]float64{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
			expected: [...]float64{
				56, 62, 68, 74,
				152, 174, 196, 218,
				248, 286, 324, 362,
				344, 398, 452, 506,
			},
		},
		{
			matrixA: [...]float64{
				0.674876, 0, 0.737931, 0,
				0, 1, 0, 0,
				-0.737931, 0, 0.674876, 0,
				0, 0, 0, 1,
			},
			matrixB: [...]float64{
				1, 0, 0, 2,
				0, 1, 0, 4,
				0, 0, 1, 6,
				0, 0, 0, 1,
			},
			expected: [...]float64{
				0.674876, 0, 0.737931, 5.777338,
				0, 1, 0, 4,
				-0.737931, 0, 0.674876, 2.5733939999999995,
				0, 0, 0, 1,
			},
		},
		{
			matrixA: [...]float64{
				1, 0, 0, 2,
				0, 1, 0, 4,
				0, 0, 1, 6,
				0, 0, 0, 1,
			},
			matrixB: [...]float64{
				0.674876, 0, 0.737931, 0,
				0, 1, 0, 0,
				-0.737931, 0, 0.674876, 0,
				0, 0, 0, 1,
			},
			expected: [...]float64{
				0.674876, 0, 0.737931, 2,
				0, 1, 0, 4,
				-0.737931, 0, 0.674876, 6,
				0, 0, 0, 1,
			},
		},
	}

	for _, test := range tests {
		matrixA := &Matrix{Values: test.matrixA}
		matrixB := &Matrix{Values: test.matrixB}

		result := MultiplyMatrix(matrixA, matrixB)
		if result.Values != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result.Values)
		}
	}
}

func TestRotationMatrix(t *testing.T) {
	tests := []struct {
		x        float64
		y        float64
		z        float64
		expected [16]float64
	}{
		{
			x: 45,
			expected: [...]float64{
				1, 0, 0, 0,
				0, 0.7071067811865476, 0.7071067811865475, 0,
				0, -0.7071067811865475, 0.7071067811865476, 0,
				0, 0, 0, 1,
			},
		},
		{
			y: 45,
			expected: [...]float64{
				0.7071067811865476, 0, -0.7071067811865475, 0,
				0, 1, 0, 0,
				0.7071067811865475, 0, 0.7071067811865476, 0,
				0, 0, 0, 1,
			},
		},
		{
			z: 45,
			expected: [...]float64{
				0.7071067811865476, 0.7071067811865475, 0, 0,
				-0.7071067811865475, 0.7071067811865476, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
		},
		{
			x: 45,
			y: 45,
			z: 45,
			expected: [...]float64{
				0.5000000000000001, 0.5, -0.7071067811865475, 0,
				-0.14644660940672627, 0.8535533905932737, 0.5, 0,
				0.8535533905932737, -0.14644660940672627, 0.5000000000000001, 0,
				0, 0, 0, 1,
			},
		},
	}

	for _, test := range tests {
		matrix := NewIdentityMatrix()
		matrix.RotateX(test.x)
		matrix.RotateY(test.y)
		matrix.RotateZ(test.z)

		if matrix.Values != test.expected {
			t.Errorf("expected %v, got %v", test.expected, matrix.Values)
		}
	}
}
