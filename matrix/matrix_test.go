package matrix

import "testing"

func TestIdentity(t *testing.T) {
	matrix := NewIdentityMatrix()
	testValues := [...]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	if matrix.values != testValues {
		t.Errorf("Expected %v, got %v", testValues, matrix.values)
	}
}

func TestPerspective(t *testing.T) {
	matrix := NewPerspectiveMatrix(75, 640/480, 1.0, 1000.0)
	testValues := [...]float32{
		1.3032254, 0, 0, 0,
		0, 1.3032254, 0, 0,
		0, 0, -1.002002, 2.002002,
		0, 0, -1, 0,
	}

	if matrix.values != testValues {
		t.Errorf("Expected %v, got %v", testValues, matrix.values)
	}
}

func TestMultiplyMatrix(t *testing.T) {
	tests := []struct {
		matrixA  [16]float32
		matrixB  [16]float32
		expected [16]float32
	}{
		{
			matrixA: [...]float32{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
			matrixB: [...]float32{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
			expected: [...]float32{
				56, 62, 68, 74,
				152, 174, 196, 218,
				248, 286, 324, 362,
				344, 398, 452, 506,
			},
		},
		{
			matrixA: [...]float32{
				0.674876, 0, 0.737931, 0,
				0, 1, 0, 0,
				-0.737931, 0, 0.674876, 0,
				0, 0, 0, 1,
			},
			matrixB: [...]float32{
				1, 0, 0, 2,
				0, 1, 0, 4,
				0, 0, 1, 6,
				0, 0, 0, 1,
			},
			expected: [...]float32{
				0.674876, 0, 0.737931, 5.777338,
				0, 1, 0, 4,
				-0.737931, 0, 0.674876, 2.5733938,
				0, 0, 0, 1,
			},
		},
		{
			matrixA: [...]float32{
				1, 0, 0, 2,
				0, 1, 0, 4,
				0, 0, 1, 6,
				0, 0, 0, 1,
			},
			matrixB: [...]float32{
				0.674876, 0, 0.737931, 0,
				0, 1, 0, 0,
				-0.737931, 0, 0.674876, 0,
				0, 0, 0, 1,
			},
			expected: [...]float32{
				0.674876, 0, 0.737931, 2,
				0, 1, 0, 4,
				-0.737931, 0, 0.674876, 6,
				0, 0, 0, 1,
			},
		},
	}

	for _, test := range tests {
		matrixA := &Matrix{values: test.matrixA}
		matrixB := &Matrix{values: test.matrixB}

		result := MultiplyMatrix(matrixA, matrixB)
		if result.values != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result.values)
		}
	}
}

func TestRotationMatrix(t *testing.T) {
	tests := []struct {
		x        float32
		y        float32
		z        float32
		expected [16]float32
	}{
		{
			x: 45,
			expected: [...]float32{
				1, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 1,
			},
		},
		{
			y: 45,
			expected: [...]float32{
				1, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 1,
			},
		},
		{
			z: 45,
			expected: [...]float32{
				1, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 1,
			},
		},
		{
			x: 45,
			y: 45,
			z: 45,
			expected: [...]float32{
				1, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 1,
			},
		},
	}

	for _, test := range tests {
		matrix := NewIdentityMatrix()
		matrix.RotateX(test.x)
		matrix.RotateY(test.y)
		matrix.RotateZ(test.z)

		if matrix.values != test.expected {
			t.Errorf("expected %v, got %v", test.expected, matrix.values)
		}
	}
}
