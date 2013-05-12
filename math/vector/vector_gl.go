package vector

import (
	gl "github.com/chsc/gogl/gl33"
)

func (v Vector2i) ToGL() [2]gl.Int {
	return [2]gl.Int{gl.Int(v.X), gl.Int(v.Y)}
}

func (v Vector3i) ToGL() [3]gl.Int {
	return [3]gl.Int{gl.Int(v.X), gl.Int(v.Y), gl.Int(v.Z)}
}

func (v Vector2f) ToGL() [2]gl.Float {
	return [2]gl.Float{gl.Float(v.X), gl.Float(v.Y)}
}

func (v Vector3f) ToGL() [3]gl.Float {
	return [3]gl.Float{gl.Float(v.X), gl.Float(v.Y), gl.Float(v.Z)}
}

func (v Vector4f) ToGL() [4]gl.Float {
	return [4]gl.Float{gl.Float(v.X), gl.Float(v.Y), gl.Float(v.Z), gl.Float(v.W)}
}
