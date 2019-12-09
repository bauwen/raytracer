package main

import (
	"math"
	"math/rand"

	"./vec3"
)

type Material interface {
    Scatter(rayIn Ray, record HitRecord, attenuation *vec3.Vec3, rayOut *Ray) bool
}

func randomUnitInSphere() vec3.Vec3 {
	p := vec3.New(0.0, 0.0, 0.0)
	for {
		p.X = rand.Float64()*2 - 1
		p.Y = rand.Float64()*2 - 1
		p.Z = rand.Float64()*2 - 1
		if vec3.LenSq(p) < 1.0 {
			break
		}
	}
	return p
}

func reflect(v, n vec3.Vec3) vec3.Vec3 {
	scalar := -vec3.Dot(v, n) * 2
	return vec3.Add(v, vec3.Scale(n, scalar))
}

func refract(v, n vec3.Vec3, niOverNt float64, refracted *vec3.Vec3) bool {
	uv := vec3.Norm(v)
	dt := vec3.Dot(uv, n)
	discriminant := 1.0 - niOverNt*niOverNt*(1 - dt*dt)
	if discriminant > 0 {
		lhs := vec3.Scale(vec3.Sub(uv, vec3.Scale(n, dt)), niOverNt)
		rhs := vec3.Scale(n, math.Sqrt(discriminant))
		*refracted = vec3.Sub(lhs, rhs)
		return true
	}
	return false
}

func schlick(cosine, refractiveIndex float64) float64 {
	r0 := (1-refractiveIndex) / (1+refractiveIndex)
	return r0*r0 + (1 - r0*r0)*math.Pow(1 - cosine, 5)
}

// lambertian
type Lambertian struct {
	Albedo vec3.Vec3
}

func (l Lambertian) Scatter(rayIn Ray, record HitRecord, attenuation *vec3.Vec3, rayOut *Ray) bool {
	target := vec3.Add(vec3.Add(record.P, record.Normal), randomUnitInSphere())
	rayOut.A = record.P
	rayOut.B = vec3.Sub(target, record.P)
	attenuation.X = l.Albedo.X
	attenuation.Y = l.Albedo.Y
	attenuation.Z = l.Albedo.Z
	return true
}


// metal
type Metal struct {
	Albedo vec3.Vec3
	Fuzz   float64
}

func (m Metal) Scatter(rayIn Ray, record HitRecord, attenuation *vec3.Vec3, rayOut *Ray) bool {
	reflected := reflect(vec3.Norm(rayIn.Direction()), record.Normal)
	rayOut.A = record.P
	if m.Fuzz > 0 {
		rayOut.B = vec3.Add(reflected, vec3.Scale(randomUnitInSphere(), m.Fuzz))
	} else {
		rayOut.B = reflected
	}
	attenuation.X = m.Albedo.X
	attenuation.Y = m.Albedo.Y
	attenuation.Z = m.Albedo.Z
	return vec3.Dot(rayOut.B, record.Normal) > 0
}


// dielectric
type Dielectric struct {
	RefractiveIndex float64  // typically air = 1, glass = 1.3-1.7, diamond = 2.4
}

func (d Dielectric) Scatter(rayIn Ray, record HitRecord, attenuation *vec3.Vec3, rayOut *Ray) bool {
	*attenuation = vec3.New(1.0, 1.0, 1.0)

	outwardNormal := vec3.New(0.0, 0.0, 0.0)
	var niOverNt float64
	var cosine float64
	if vec3.Dot(rayIn.Direction(), record.Normal) > 0 {
		outwardNormal = vec3.Scale(record.Normal, -1)
		niOverNt = d.RefractiveIndex
		cosine = d.RefractiveIndex * vec3.Dot(rayIn.Direction(), record.Normal) / vec3.Len(rayIn.Direction())
	} else {
		outwardNormal = record.Normal
		niOverNt = 1.0 / d.RefractiveIndex
		cosine = -vec3.Dot(rayIn.Direction(), record.Normal) / vec3.Len(rayIn.Direction())
	}
	
	refracted := vec3.New(0.0, 0.0, 0.0)
	rayOut.A = record.P
	reflectProb := 1.0
	if refract(rayIn.Direction(), outwardNormal, niOverNt, &refracted) {
		reflectProb = schlick(cosine, d.RefractiveIndex)
	}
	if rand.Float64() < reflectProb {
		rayOut.B = reflect(vec3.Norm(rayIn.Direction()), record.Normal)
	} else {
		rayOut.B = refracted
	}
	return true
}