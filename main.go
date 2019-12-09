package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"

	"./vec3"
)

// Simple raytracer implementation, based on the tutorial found at
// http://www.realtimerendering.com/raytracing/Ray%20Tracing%20in%20a%20Weekend.pdf

// comment out a setup scene (may take several seconds to execute)
// arguments: width and height of picture + number of samples for antialiasing
// lower values = faster execution, higher values = more polished result
func main() {
	//setup1(400, 200, 100)
	setup2(500, 500, 50)
	//setup3(300, 300, 50)
	//setup4(200, 200, 5, "elephant.stl")
	//setup4(200, 200, 5, "tyranitar.stl")
}

const MAXFLOAT = 999999.99

var BACKGROUND_COLOR vec3.Vec3

// point light
type Light struct {
	P         vec3.Vec3
	Intensity vec3.Vec3
	Color     vec3.Vec3
}

// the actual ray tracing happens here
func pixel(ray Ray, lights []Light, world HitableList, depth int) vec3.Vec3 {
	record := HitRecord{}
	if world.Hit(ray, 0.001, MAXFLOAT, &record) {
		// comment out to see normal map
		//x, y, z := record.Normal.X, record.Normal.Y, record.Normal.Z
		//return vec3.Scale(vec3.New(x+1, y+1, z+1), 0.5)

		attenuation := vec3.New(0.0, 0.0, 0.0)
		var current vec3.Vec3
		rayOut := Ray{}
		if depth < 10 && record.Material.Scatter(ray, record, &attenuation, &rayOut) {
			//return vec3.Mul(pixel(rayOut, lights, world, depth+1), attenuation)
			current = vec3.Mul(pixel(rayOut, lights, world, depth+1), attenuation)
		}

		for _, light := range lights {
			if _, ok := record.Material.(Dielectric); ok {
				continue
			}
			shadowRay := Ray{record.P, vec3.Sub(light.P, record.P)}
			rec := HitRecord{}
			if world.Hit(shadowRay, 0.001, MAXFLOAT, &rec) {
				if vec3.LenSq(vec3.Sub(light.P, shadowRay.A)) > vec3.LenSq(vec3.Sub(rec.P, shadowRay.A)) {
					continue
				}
			}
			var albedo vec3.Vec3
			switch m := record.Material.(type) {
			case Lambertian:
				albedo = m.Albedo
			case Metal:
				albedo = m.Albedo
			}
			d := vec3.Dot(record.Normal, shadowRay.Direction()) / vec3.LenSq(vec3.Sub(light.P, shadowRay.A))
			if d < 0 {
				d = 0
			}

			res := vec3.Scale(vec3.Mul(vec3.Mul(albedo, light.Intensity), light.Color), d)
			current = vec3.Add(current, res)
			if current.X > 1 {
				current.X = 1.0
			}
			if current.Y > 1 {
				current.Y = 1.0
			}
			if current.Z > 1 {
				current.Z = 1.0
			}
		}
		return current
		//return vec3.New(0.0, 0.0, 0.0)
	}

	// background
	//unitDirection := vec3.Norm(ray.Direction())
	//t := 0.5 * (unitDirection.Y + 1.0)
	//from := vec3.New(0.0, 0.0, 0.0)
	//to := vec3.New(0, 0, 0)
	//from := vec3.New(1.0, 1.0, 1.0)
	//to := vec3.New(0.5, 0.7, 1.0)
	return BACKGROUND_COLOR //vec3.Add(vec3.Scale(from, 1.0-t), vec3.Scale(to, t))
}

func setup1(nx, ny, ns int) {
	BACKGROUND_COLOR = vec3.New(0.6, 0.8, 1.0)

	var lookFrom, lookAt vec3.Vec3
	var vfov float64
	world, lights := createSampleScene(&lookFrom, &lookAt, &vfov)

	setupExecute(nx, ny, ns, lookFrom, lookAt, vfov, world, lights)
}

func setup2(nx, ny, ns int) {
	BACKGROUND_COLOR = vec3.New(0.6, 0.8, 1.0)

	var lookFrom, lookAt vec3.Vec3
	var vfov float64
	world, lights := createAwesomeScene(&lookFrom, &lookAt, &vfov)

	setupExecute(nx, ny, ns, lookFrom, lookAt, vfov, world, lights)
}

func setup3(nx, ny, ns int) {
	BACKGROUND_COLOR = vec3.New(0.0, 0.0, 0.0)

	var lookFrom, lookAt vec3.Vec3
	var vfov float64
	world, lights := createTriangleScene(&lookFrom, &lookAt, &vfov)

	setupExecute(nx, ny, ns, lookFrom, lookAt, vfov, world, lights)
}

func setup4(nx, ny, ns int, filename string) {
	list, err := loadBinarySTLModel(filename, vec3.New(0.8, 0.1, 0.6), 1.0, vec3.New(0, 2, 0))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d triangles\n", len(list))

	BACKGROUND_COLOR = vec3.New(0.0, 0.0, 0.0)

	var lookFrom, lookAt vec3.Vec3
	var vfov float64
	world, lights := createModelScene(&lookFrom, &lookAt, &vfov)
	index := len(world)
	for _, tr := range list {
		world = append(world, tr)
		index += 1
	}

	setupExecute(nx, ny, ns, lookFrom, lookAt, vfov, world, lights)
}

func setupExecute(nx, ny, ns int, lookFrom, lookAt vec3.Vec3, vfov float64, world HitableList, lights []Light) {
	pixels := image.NewRGBA(image.Rect(0, 0, nx, ny))

	upVector := vec3.New(0, 1, 0)
	aspect := float64(nx) / float64(ny)
	camera := NewPinholeCamera(lookFrom, lookAt, upVector, vfov, aspect)

	mutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for j := ny - 1; j >= 0; j-- {
		wg.Add(1)
		go raytracer(pixels, j, nx, ny, ns, mutex, wg, camera, world, lights)
	}
	wg.Wait()

	f, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, pixels); err != nil {
		panic(err)
	}
}

func raytracer(pixels *image.RGBA, j, nx, ny, ns int, mutex *sync.Mutex, wg *sync.WaitGroup, camera PinholeCamera, world HitableList, lights []Light) {
	cs := make([]color.RGBA, nx)
	for i := 0; i < nx; i++ {
		// antialiasing (average of `ns` samples per pixel)
		col := vec3.New(0, 0, 0)
		for s := 0; s < ns; s++ {
			u := (float64(i) + rand.Float64()) / float64(nx)
			v := (float64(j) + rand.Float64()) / float64(ny)

			ray := camera.GetRay(u, v)
			col = vec3.Add(col, pixel(ray, lights, world, 0))
		}
		col = vec3.Scale(col, 1.0/float64(ns))
		col = vec3.New(math.Sqrt(col.X), math.Sqrt(col.Y), math.Sqrt(col.Z))

		c := color.RGBA{
			uint8(255 * col.X),
			uint8(255 * col.Y),
			uint8(255 * col.Z),
			255,
		}
		cs[i] = c
	}
	mutex.Lock()
	for i := 0; i < nx; i++ {
		pixels.SetRGBA(i, ny-1-j, cs[i])
	}
	mutex.Unlock()
	wg.Done()
}
