package main

import (
	"math"
	"./vec3"
)

type Hitable interface {
    Hit(ray Ray, tMin, tMax float64, record *HitRecord) bool
}

type HitRecord struct {
	T float64
	P vec3.Vec3
	Normal vec3.Vec3
	Material Material
	HitLight bool
	Light Light
}


// list of hitables
type HitableList []Hitable

func (l HitableList) Hit(ray Ray, tMin, tMax float64, record *HitRecord) bool {
	tempRecord := HitRecord{}
	hitAnything := false
	closestSoFar := tMax
	for _, hitable := range l {
		if hitable.Hit(ray, tMin, closestSoFar, &tempRecord) {
			hitAnything = true
			closestSoFar = tempRecord.T
			record.T = tempRecord.T
			record.P = tempRecord.P
			record.Normal = tempRecord.Normal
			record.Material = tempRecord.Material
			record.HitLight = tempRecord.HitLight
			record.Light = tempRecord.Light
		}
	}
	return hitAnything
}


// sphere
type Sphere struct {
    Center vec3.Vec3
	Radius float64
	Material Material
}

func (s Sphere) Hit(ray Ray, tMin, tMax float64, record *HitRecord) bool {
    oc := vec3.Sub(ray.Origin(), s.Center)
	a := vec3.Dot(ray.Direction(), ray.Direction())
	b := 2.0 * vec3.Dot(oc, ray.Direction())
	c := vec3.Dot(oc, oc) - s.Radius*s.Radius
	discriminant := b*b - 4*a*c
	if discriminant > 0 {
		t := (-b - math.Sqrt(discriminant)) / (2*a)
		if tMin < t && t < tMax {
			record.T = t
			record.P = ray.PointAtParameter(t)
			record.Normal = vec3.Scale(vec3.Sub(record.P, s.Center), 1.0/s.Radius)
			record.Material = s.Material
			return true
		}
		t = (-b + math.Sqrt(discriminant)) / (2*a)
		if tMin < t && t < tMax {
			record.T = t
			record.P = ray.PointAtParameter(t)
			record.Normal = vec3.Scale(vec3.Sub(record.P, s.Center), 1.0/s.Radius)
			record.Material = s.Material
			return true
		}
	}
	return false
}


// plane (infinite, no borders)
type Plane struct {
	Point vec3.Vec3
	Normal vec3.Vec3
	Material Material
}

func (p Plane) Hit(ray Ray, tMin, tMax float64, record *HitRecord) bool {
	n := vec3.Norm(p.Normal)
	denom := vec3.Dot(n, ray.Direction())
	if math.Abs(denom) < 0.0001 {
		return false
	}
	nom := vec3.Dot(n, ray.Origin()) - vec3.Dot(n, p.Point)
	t := -nom / denom
	if tMin < t && t < tMax {
		record.T = t
		record.P = ray.PointAtParameter(t)
		record.Normal = n
		record.Material = p.Material
		return true
	}
	return false
}


// triangle
type Triangle struct {
	Vertex1 vec3.Vec3
	Vertex2 vec3.Vec3
	Vertex3 vec3.Vec3
	Material Material
}

// https://en.wikipedia.org/wiki/Möller–Trumbore_intersection_algorithm
func (tr Triangle) Hit(ray Ray, tMin, tMax float64, record *HitRecord) bool {
	edge1 := vec3.Sub(tr.Vertex2, tr.Vertex1)
	edge2 := vec3.Sub(tr.Vertex3, tr.Vertex1)
	h := vec3.Cross(ray.Direction(), edge2)
	a := vec3.Dot(edge1, h)
	if -0.0001 < a && a < 0.0001 {
		return false
	}

	f := 1.0 / a
	s := vec3.Sub(ray.Origin(), tr.Vertex1)
	u := f * vec3.Dot(s, h)
	if u < 0.0 || 1.0 < u {
		return false
	}

	q := vec3.Cross(s, edge1)
	v := f * vec3.Dot(ray.Direction(), q)
	if v < 0.0 || 1.0 < v + u {
		return false
	}

	t := f * vec3.Dot(edge2, q)
	if tMin < t && t < tMax {
		record.T = t
		record.P = ray.PointAtParameter(t)
		record.Normal = vec3.Norm(vec3.Cross(edge1, edge2))
		record.Material = tr.Material
		return true
	}
	
	return false
}