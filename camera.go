package main

import (
	"math"
	"math/rand"

    "./vec3"
)

type Camera interface {
	GetRay(u, v float64) Ray
}

func randomInUnitDisk() vec3.Vec3 {
	p := vec3.New(0.0, 0.0, 0.0)
	for {
		p.X = rand.Float64()*2 - 1
		p.Y = rand.Float64()*2 - 1
		p.Z = 0
		if vec3.LenSq(p) < 1.0 {
			break
		}
	}
	return p
}


// pinhole camera
type PinholeCamera struct {
    LowerLeftCorner vec3.Vec3
    Horizontal vec3.Vec3
    Vertical vec3.Vec3
    Origin vec3.Vec3
}

/*
func NewPinholeCamera(vfov, aspect float64) Camera {
	theta := vfov * math.Pi / 180
	halfHeight := math.Tan(theta/2)
	halfWidth := aspect * halfHeight
    return PinholeCamera{
        LowerLeftCorner: vec3.New(-halfWidth, -halfHeight, -1.0),
        Horizontal: vec3.New(2*halfWidth, 0.0, 0.0),
        Vertical: vec3.New(0.0, 2*halfHeight, 0.0),
        Origin: vec3.New(0.0, 0.0, 0.0),
    }
}
*/
func NewPinholeCamera(lookFrom, lookAt, vup vec3.Vec3, vfov, aspect float64) PinholeCamera {
	theta := vfov * math.Pi / 180
	halfHeight := math.Tan(theta/2)
	halfWidth := aspect * halfHeight

	w := vec3.Norm(vec3.Sub(lookFrom, lookAt))
	u := vec3.Norm(vec3.Cross(vup, w))
	v := vec3.Cross(w, u)

	llc := lookFrom
	llc = vec3.Sub(llc, vec3.Scale(u, halfWidth))
	llc = vec3.Sub(llc, vec3.Scale(v, halfHeight))
	llc = vec3.Sub(llc, w)

    return PinholeCamera{
        LowerLeftCorner: llc,
        Horizontal: vec3.Scale(u, 2*halfWidth),
        Vertical: vec3.Scale(v, 2*halfHeight),
        Origin: lookFrom,
	}
}

func (c PinholeCamera) GetRay(u, v float64) Ray {
	dx := vec3.Scale(c.Horizontal, u)
	dy := vec3.Scale(c.Vertical, v)
	direction := c.LowerLeftCorner
	direction = vec3.Add(direction, dx)
	direction = vec3.Add(direction, dy)
	direction = vec3.Sub(direction, c.Origin)
	return Ray{c.Origin, direction}
}


// lens camera (depth of field, "defocus blur")
type LensCamera struct {
    LowerLeftCorner vec3.Vec3
    Horizontal vec3.Vec3
    Vertical vec3.Vec3
	Origin vec3.Vec3
	U, V, W vec3.Vec3
	LensRadius float64
}

func NewLensCamera(lookFrom, lookAt, vup vec3.Vec3, vfov, aspect, aperture, focusDist float64) LensCamera {
	theta := vfov * math.Pi / 180
	halfHeight := math.Tan(theta/2)
	halfWidth := aspect * halfHeight

	w := vec3.Norm(vec3.Sub(lookFrom, lookAt))
	u := vec3.Norm(vec3.Cross(vup, w))
	v := vec3.Cross(w, u)

	llc := lookFrom
	llc = vec3.Sub(llc, vec3.Scale(u, halfWidth * focusDist))
	llc = vec3.Sub(llc, vec3.Scale(v, halfHeight * focusDist))
	llc = vec3.Sub(llc, vec3.Scale(w, focusDist))

    return LensCamera{
        LowerLeftCorner: llc,
        Horizontal: vec3.Scale(u, 2*halfWidth * focusDist),
        Vertical: vec3.Scale(v, 2*halfHeight * focusDist),
		Origin: lookFrom,
		U: u,
		V: v,
		W: w,
		LensRadius: aperture / 2,
	}
}

func (c LensCamera) GetRay(s, t float64) Ray {
	rd := vec3.Scale(randomInUnitDisk(), c.LensRadius)
	offset := vec3.Add(vec3.Scale(c.U, rd.X), vec3.Scale(c.V, rd.Y))

	dx := vec3.Scale(c.Horizontal, s)
	dy := vec3.Scale(c.Vertical, t)
	direction := c.LowerLeftCorner
	direction = vec3.Add(direction, dx)
	direction = vec3.Add(direction, dy)
	direction = vec3.Sub(direction, c.Origin)
	return Ray{vec3.Add(c.Origin, offset), vec3.Sub(direction, offset)}
}
