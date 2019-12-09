package main

import (
    "./vec3"
)

type Ray struct {
    A, B vec3.Vec3
}

func (r Ray) Origin() vec3.Vec3 {
    return r.A
}

func (r Ray) Direction() vec3.Vec3 {
    return r.B
}

func (r Ray) PointAtParameter(t float64) vec3.Vec3 {
    return vec3.Add(r.A, vec3.Scale(r.B, t))
}

