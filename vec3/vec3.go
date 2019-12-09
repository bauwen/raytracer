package vec3

import (
    "fmt"
    "math"
)

type Vec3 struct {
    X, Y, Z float64
}

func (v Vec3) String() string {
    return fmt.Sprintf("vec3(%v, %v, %v)", v.X, v.Y, v.Z)
}

func New(x, y, z float64) Vec3 {
	return Vec3{x, y, z}
}

func Add(v1, v2 Vec3) Vec3 {
    return Vec3{
        v1.X + v2.X,
        v1.Y + v2.Y,
        v1.Z + v2.Z,
    }
}

func Sub(v1, v2 Vec3) Vec3 {
    return Vec3{
        v1.X - v2.X,
        v1.Y - v2.Y,
        v1.Z - v2.Z,
    }
}

func Mul(v1, v2 Vec3) Vec3 {
    return Vec3{
        v1.X * v2.X,
        v1.Y * v2.Y,
        v1.Z * v2.Z,
    }
}

func Cross(v1, v2 Vec3) Vec3 {
    return Vec3{
        v1.Y*v2.Z - v1.Z*v2.Y,
        v1.Z*v2.X - v1.X*v2.Z,
        v1.X*v2.Y - v1.Y*v2.X,
    }
}

func Dot(v1, v2 Vec3) float64 {
    return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func LenSq(v Vec3) float64 {
    return Dot(v, v)
}

func Len(v Vec3) float64 {
    return math.Sqrt(LenSq(v))
}

func Norm(v Vec3) Vec3 {
	return Scale(v, 1.0/Len(v))
}

func Scale(v Vec3, s float64) Vec3 {
	return Vec3{
		v.X * s,
		v.Y * s,
		v.Z * s,
	}
}

func Translate(v Vec3, d float64) Vec3 {
	return Vec3{
		v.X + d,
		v.Y + d,
		v.Z + d,
	}
}
