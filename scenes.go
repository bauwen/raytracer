package main

import (
	"math/rand"

	"./vec3"
)

func createSampleScene(lookFrom, lookAt *vec3.Vec3, fov *float64) (HitableList, []Light) {
	*lookFrom = vec3.New(0, 0, 0.8)
	*lookAt = vec3.New(0, 0, -1)
	*fov = 60.0
	
	world := make(HitableList, 5)
	world[0] = Sphere{
		Center: vec3.New(0.0, 0.0, -1.0),
		Radius: 0.5,
		Material: Lambertian{vec3.New(0.1, 0.2, 0.5)},
	}
	world[1] = Sphere{
		Center: vec3.New(0.0, -100.5, -1),
		Radius: 100,
		Material: Lambertian{vec3.New(0.8, 0.8, 0.0)},
	}
	world[2] = Sphere{
		Center: vec3.New(1.0, 0.0, -1.0),
		Radius: 0.5,
		Material: Metal{vec3.New(0.8, 0.6, 0.2), 0.0},
	}
	world[3] = Sphere{
		Center: vec3.New(-1.0, 0.0, -1.0),
		Radius: 0.5,
		Material: Dielectric{1.5},
	}
	world[4] = Sphere{
		Center: vec3.New(-1.0, 0.0, -1.0),
		Radius: -0.45,
		Material: Dielectric{1.5},
	}

	return world, nil
}

func makeRectangle(p1, p2, p3 vec3.Vec3, m Material) (Triangle, Triangle) {
	t1 := Triangle{p1, p2, p3, m}
	p4 := vec3.Add(p1, vec3.Sub(p3, p2))
	t2 := Triangle{p3, p4, p1, m}
	return t1, t2
}

func createTriangleScene(lookFrom, lookAt *vec3.Vec3, fov *float64) (HitableList, []Light) {
	*lookFrom = vec3.New(0, 4.5, 14.0)
	*lookAt = vec3.New(0, 4.0, -1)
	*fov = 45.0

	w := 4.5

	world := make(HitableList, 14)
	world[0], world[1] = makeRectangle(
		vec3.New(-w, 0.0, w),
		vec3.New(w, 0.0, w),
		vec3.New(w, 0.0, -w),
		Lambertian{vec3.New(0.2, 0.6, 0.1)},
	)
	world[2], world[3] = makeRectangle(
		vec3.New(-w, 0.0, -w),
		vec3.New(w, 0.0, -w),
		vec3.New(w, w*2, -w),
		Metal{vec3.New(1.0, 1.0, 1.0), 0.0},//Lambertian{vec3.New(0.8, 0.6, 0.1)},
	)
	world[4], world[5] = makeRectangle(
		vec3.New(-w, 0.0, w),
		vec3.New(-w, 0.0, -w),
		vec3.New(-w, w*2, -w),
		Lambertian{vec3.New(0.6, 0.1, 0.1)},
	)
	world[6], world[7] = makeRectangle(
		vec3.New(w, w*2, -w),
		vec3.New(w, 0.0, -w),
		vec3.New(w, 0.0, w),
		Lambertian{vec3.New(0.1, 0.1, 0.6)},
	)
	world[8], world[9] = makeRectangle(
		vec3.New(-w, w*2, -w),
		vec3.New(w, w*2, -w),
		vec3.New(w, w*2, w),
		Lambertian{vec3.New(0.6, 0.6, 0.1)},
	)
	world[10] = Sphere{
		Center: vec3.New(-2.0, 0.8, 1.0),
		Radius: 0.8,
		Material: Lambertian{vec3.New(0.1, 0.2, 0.5)},
	}
	world[11] = Sphere{
		Center: vec3.New(1.5, 1.6, -1.0),
		Radius: 1.6,
		Material: Metal{vec3.New(0.8, 0.6, 0.2), 0.0},
	}
	world[12] = Sphere{
		Center: vec3.New(0.0, 0.5, 2.5),
		Radius: 0.5,
		Material: Dielectric{1.5},
	}
	world[13] = Sphere{
		Center: vec3.New(0.0, 0.5, 2.5),
		Radius: -0.45,
		Material: Dielectric{1.5},
	}

	var lights []Light
	lights = append(lights, Light{
		P: vec3.New(-3.0, 0.9, 1.0),
		Intensity: vec3.New(1.0, 1.0, 1.0),
		Color: vec3.New(1.0, 1.0, 1.0),
	})

	/*lights = append(lights, Light{
		P: vec3.New(3.0, 3.5, -1.0),
		Intensity: vec3.New(1.0, 1.0, 1.0),
		Color: vec3.New(1.0, 1.0, 1.0),
	})*/

	lights = append(lights, Light{
		P: vec3.New(0.0, 6.8, 3.0),
		Intensity: vec3.New(1.0, 1.0, 1.0),
		Color: vec3.New(1.0, 1.0, 1.0),
	})

	return world, lights
}

func createModelScene(lookFrom, lookAt *vec3.Vec3, fov *float64) (HitableList, []Light) {
	*lookFrom = vec3.New(-10.0, 3.5, -4.0)
	*lookAt = vec3.New(0, 3.0, 0.0)
	*fov = 45.0

	w := 4.5

	world := make(HitableList, 10, 1000)
	world[0], world[1] = makeRectangle(
		vec3.New(-w, 0.0, w),
		vec3.New(w, 0.0, w),
		vec3.New(w, 0.0, -w),
		Lambertian{vec3.New(0.2, 0.6, 0.1)},
	)
	world[2], world[3] = makeRectangle(
		vec3.New(-w, 0.0, -w),
		vec3.New(w, 0.0, -w),
		vec3.New(w, w*2, -w),
		Lambertian{vec3.New(0.6, 0.1, 0.1)},//Metal{vec3.New(1.0, 1.0, 1.0), 0.0},//Lambertian{vec3.New(0.8, 0.6, 0.1)},
	)
	world[4], world[5] = makeRectangle(
		
		vec3.New(w, w*2, w),
		vec3.New(w, 0.0, w),
		vec3.New(-w, 0.0, w),
		Lambertian{vec3.New(0.1, 0.1, 0.6)},
	)
	world[6], world[7] = makeRectangle(
		vec3.New(w, w*2, -w),
		vec3.New(w, 0.0, -w),
		vec3.New(w, 0.0, w),
		Metal{vec3.New(1.0, 1.0, 1.0), 0.0},
	)
	world[8], world[9] = makeRectangle(
		vec3.New(-w, w*2, -w),
		vec3.New(w, w*2, -w),
		vec3.New(w, w*2, w),
		Lambertian{vec3.New(0.6, 0.6, 0.1)},
	)//*/

	var lights []Light
	lights = append(lights, Light{
		P: vec3.New(-3.0, 0.9, 1.0),
		Intensity: vec3.New(1.0, 1.0, 1.0),
		Color: vec3.New(1.0, 1.0, 1.0),
	})

	lights = append(lights, Light{
		P: vec3.New(0.0, 2.5, -3.0),
		Intensity: vec3.New(1.0, 1.0, 1.0),
		Color: vec3.New(1.0, 1.0, 1.0),
	})

	lights = append(lights, Light{
		P: vec3.New(0.0, 3.8, 3.0),
		Intensity: vec3.New(1.0, 1.0, 1.0),
		Color: vec3.New(1.0, 1.0, 1.0),
	})

	return world, lights
}

func createAwesomeScene(lookFrom, lookAt *vec3.Vec3, fov *float64) (HitableList, []Light) {
	*lookFrom = vec3.New(6.0, 1.7, 3.0)
	*lookAt = vec3.New(0, 0, -1)
	*fov = 90.0

	n := 488
	world := make([]Hitable, n)
	world[0] = Sphere{
		Center: vec3.New(0, -1000, 0),
		Radius: 1000,
		Material: Lambertian{vec3.New(0.5, 0.5, 0.5)},
	}

	i := 1
	for a := -11.0; a < 11.0; a++ {
		for b := -11.0; b < 11.0; b++ {
			chooseMat := rand.Float64()
			center := vec3.New(a + 0.9*rand.Float64(), 0.2, b + 0.9*rand.Float64())
			if vec3.Len(vec3.Sub(center, vec3.New(4, 0.2, 0))) > 0.9 {
				if chooseMat < 0.8 { // diffuse
					r := rand.Float64()*rand.Float64()
					g := rand.Float64()*rand.Float64()
					b := rand.Float64()*rand.Float64()
					world[i] = Sphere{
						Center: center,
						Radius: 0.2,
						Material: Lambertian{vec3.New(r, g, b)},
					}
					i += 1
				} else if chooseMat < 0.95 { // metal
					r := 0.5*(1 + rand.Float64())
					g := 0.5*(1 + rand.Float64())
					b := 0.5*(1 + rand.Float64())
					fuzz := 0.5*rand.Float64()
					world[i] = Sphere{
						Center: center,
						Radius: 0.2,
						Material: Metal{vec3.New(r, g, b), fuzz},
					}
					i += 1
				} else { // glass
					world[i] = Sphere{
						Center: center,
						Radius: 0.2,
						Material: Dielectric{1.5},
					}
					i += 1
				}
			}
		}
	}

	world[i+0] = Sphere{vec3.New(0, 1, 0), 1.0, Dielectric{1.5}}
	world[i+1] = Sphere{vec3.New(-4, 1, 0), 1.0, Lambertian{vec3.New(0.4, 0.2, 0.1)}}
	world[i+2] = Sphere{vec3.New(4, 1, 0), 1.0, Metal{vec3.New(0.7, 0.6, 0.5), 0.0}}

	return world[:i+3], nil
}
