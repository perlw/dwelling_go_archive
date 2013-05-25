package vector

import "math"

type Vector2i struct {
	X, Y int
}

type Vector3i struct {
	X, Y, Z int
}

type Vector2f struct {
	X, Y float64
}

type Vector3f struct {
	X, Y, Z float64
}

type Vector4f struct {
	X, Y, Z, W float64
}

func DotProduct(v1, v2 Vector3f) float64 {
	return (v1.X * v2.X) + (v1.Y * v2.Y) + (v1.Z * v2.Z)
}

func Vector3fTo4f(vec Vector3f, w float64) Vector4f {
	return Vector4f{vec.X, vec.Y, vec.Z, w}
}

func Vector4fTo3f(vec Vector4f) Vector3f {
	return Vector3f{vec.X, vec.Y, vec.Z}
}

// Vector2i
func (v Vector2i) Add(v2 Vector2i) Vector2i {
	return Vector2i{v.X + v2.X, v.Y + v2.Y}
}

func (v Vector2i) Sub(v2 Vector2i) Vector2i {
	return Vector2i{v.X - v2.X, v.Y - v2.Y}
}

func (v Vector2i) Mul(v2 Vector2i) Vector2i {
	return Vector2i{v.X * v2.X, v.Y * v2.Y}
}

func (v Vector2i) Div(v2 Vector2i) Vector2i {
	return Vector2i{v.X / v2.X, v.Y / v2.Y}
}

func (v Vector2i) AddScalar(s int) Vector2i {
	return Vector2i{v.X + s, v.Y + s}
}

func (v Vector2i) SubScalar(s int) Vector2i {
	return Vector2i{v.X - s, v.Y - s}
}

func (v Vector2i) MulScalar(s int) Vector2i {
	return Vector2i{v.X * s, v.Y * s}
}

func (v Vector2i) DivScalar(s int) Vector2i {
	return Vector2i{v.X / s, v.Y / s}
}

// Vector3i
func (v Vector3i) Add(v2 Vector3i) Vector3i {
	return Vector3i{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v Vector3i) Sub(v2 Vector3i) Vector3i {
	return Vector3i{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v Vector3i) Mul(v2 Vector3i) Vector3i {
	return Vector3i{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z}
}

func (v Vector3i) Div(v2 Vector3i) Vector3i {
	return Vector3i{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z}
}

func (v Vector3i) AddScalar(s int) Vector3i {
	return Vector3i{v.X + s, v.Y + s, v.Z + s}
}

func (v Vector3i) SubScalar(s int) Vector3i {
	return Vector3i{v.X - s, v.Y - s, v.Z - s}
}

func (v Vector3i) MulScalar(s int) Vector3i {
	return Vector3i{v.X * s, v.Y * s, v.Z * s}
}

func (v Vector3i) DivScalar(s int) Vector3i {
	return Vector3i{v.X / s, v.Y / s, v.Z / s}
}

// Vector2f
func (v Vector2f) Add(v2 Vector2f) Vector2f {
	return Vector2f{v.X + v2.X, v.Y + v2.Y}
}

func (v Vector2f) Sub(v2 Vector2f) Vector2f {
	return Vector2f{v.X - v2.X, v.Y - v2.Y}
}

func (v Vector2f) Mul(v2 Vector2f) Vector2f {
	return Vector2f{v.X * v2.X, v.Y * v2.Y}
}

func (v Vector2f) Div(v2 Vector2f) Vector2f {
	return Vector2f{v.X / v2.X, v.Y / v2.Y}
}

func (v Vector2f) AddScalar(s float64) Vector2f {
	return Vector2f{v.X + s, v.Y + s}
}

func (v Vector2f) SubScalar(s float64) Vector2f {
	return Vector2f{v.X - s, v.Y - s}
}

func (v Vector2f) MulScalar(s float64) Vector2f {
	return Vector2f{v.X * s, v.Y * s}
}

func (v Vector2f) DivScalar(s float64) Vector2f {
	return Vector2f{v.X / s, v.Y / s}
}

func (v Vector2f) Length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

func (v Vector2f) Normalize() Vector2f {
	return v.DivScalar(v.Length())
}

// Vector3f
func (v Vector3f) Add(v2 Vector3f) Vector3f {
	return Vector3f{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v Vector3f) Sub(v2 Vector3f) Vector3f {
	return Vector3f{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v Vector3f) Mul(v2 Vector3f) Vector3f {
	return Vector3f{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z}
}

func (v Vector3f) Div(v2 Vector3f) Vector3f {
	return Vector3f{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z}
}

func (v Vector3f) AddScalar(s float64) Vector3f {
	return Vector3f{v.X + s, v.Y + s, v.Z + s}
}

func (v Vector3f) SubScalar(s float64) Vector3f {
	return Vector3f{v.X - s, v.Y - s, v.Z - s}
}

func (v Vector3f) MulScalar(s float64) Vector3f {
	return Vector3f{v.X * s, v.Y * s, v.Z * s}
}

func (v Vector3f) DivScalar(s float64) Vector3f {
	return Vector3f{v.X / s, v.Y / s, v.Z / s}
}

func (v Vector3f) Length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z))
}

func (v Vector3f) Normalize() Vector3f {
	return v.DivScalar(v.Length())
}

// Vector4f
func (v Vector4f) Add(v2 Vector4f) Vector4f {
	return Vector4f{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z, v.W + v2.W}
}

func (v Vector4f) Sub(v2 Vector4f) Vector4f {
	return Vector4f{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z, v.W - v2.W}
}

func (v Vector4f) Mul(v2 Vector4f) Vector4f {
	return Vector4f{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z, v.W * v2.W}
}

func (v Vector4f) Div(v2 Vector4f) Vector4f {
	return Vector4f{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z, v.W / v2.W}
}

func (v Vector4f) AddScalar(s float64) Vector4f {
	return Vector4f{v.X + s, v.Y + s, v.Z + s, v.W + s}
}

func (v Vector4f) SubScalar(s float64) Vector4f {
	return Vector4f{v.X - s, v.Y - s, v.Z - s, v.W - s}
}

func (v Vector4f) MulScalar(s float64) Vector4f {
	return Vector4f{v.X * s, v.Y * s, v.Z * s, v.Z * s}
}

func (v Vector4f) DivScalar(s float64) Vector4f {
	return Vector4f{v.X / s, v.Y / s, v.Z / s, v.W / s}
}

func (v Vector4f) Length() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y) + (v.Z * v.Z) + (v.W * v.W))
}

func (v Vector4f) Normalize() Vector4f {
	return v.DivScalar(v.Length())
}
