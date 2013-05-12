package matrix

import (
	gl "github.com/chsc/gogl/gl33"
)

func (m *Matrix) ToGL() [16]gl.Float {
	var glMatrix [16]gl.Float

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			glMatrix[(x*4)+y] = gl.Float(m.Values[(y*4)+x])
		}
	}

	return glMatrix
}
