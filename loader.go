package main

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"io/ioutil"
	
	"math"
	
	"./vec3"
)

func loadBinarySTLModel(filename string, color vec3.Vec3, scale float64, translation vec3.Vec3) ([]Triangle, error) {
    contents, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
	}
	fmt.Printf("Loaded 3D model `%s`\n", filename)

    var list []Triangle

    r := bytes.NewReader(contents[80:])
    var amount uint32
    binary.Read(r, binary.LittleEndian, &amount)

    for i := 0; i < int(amount); i++ {
        for j := 0; j < 3*4; j++ {
            r.ReadByte()  // skip normal vector
        }
        vertices := make([]float32, 3*3)
        for j := 0; j < 3*3; j++ {
            binary.Read(r, binary.LittleEndian, &vertices[j])
        }
        list = append(list, Triangle{
            vec3.New(float64(vertices[0]), float64(vertices[1]), float64(vertices[2])),
            vec3.New(float64(vertices[3]), float64(vertices[4]), float64(vertices[5])),
            vec3.New(float64(vertices[6]), float64(vertices[7]), float64(vertices[8])),
            Lambertian{color},
		})
        for j := 0; j < 2; j++ {
            r.ReadByte()  // skip attribute byte count
        }
	}
	
	// normalize vertices
	N := float64(amount * 3)
	avgX := 0.0
	avgY := 0.0
	avgZ := 0.0
	for _, tr := range list {
		avgX += (tr.Vertex1.X + tr.Vertex2.X + tr.Vertex3.X) / N
		avgY += (tr.Vertex1.Y + tr.Vertex2.Y + tr.Vertex3.Y) / N
		avgZ += (tr.Vertex1.Z + tr.Vertex2.Z + tr.Vertex3.Z) / N
	}
	varX := 0.0
	varY := 0.0
	varZ := 0.0
	sq := func (x float64) float64 { return x*x }
	for _, tr := range list {
		varX += (sq(tr.Vertex1.X - avgX) + sq(tr.Vertex2.X - avgX) + sq(tr.Vertex3.X - avgX)) / N
		varY += (sq(tr.Vertex1.Y - avgY) + sq(tr.Vertex2.Y - avgY) + sq(tr.Vertex3.Y - avgY)) / N
		varZ += (sq(tr.Vertex1.Z - avgZ) + sq(tr.Vertex2.Z - avgZ) + sq(tr.Vertex3.Z - avgZ)) / N
	}
	varX = math.Sqrt(varX)
	varY = math.Sqrt(varY)
	varZ = math.Sqrt(varZ)
	for i := 0; i < len(list); i++ {
		tr := &list[i]

		tx := tr.Vertex1.X
		ty := tr.Vertex1.Y
		tz := tr.Vertex1.Z

		tr.Vertex1.X = (tr.Vertex3.X - avgX) / varX * scale + translation.X
		tr.Vertex1.Y = (tr.Vertex3.Y - avgY) / varY * scale + translation.Y
		tr.Vertex1.Z = (tr.Vertex3.Z - avgZ) / varZ * -scale + translation.Z

		tr.Vertex2.X = (tr.Vertex2.X - avgX) / varX * scale + translation.X
		tr.Vertex2.Y = (tr.Vertex2.Y - avgY) / varY * scale + translation.Y
		tr.Vertex2.Z = (tr.Vertex2.Z - avgZ) / varZ * -scale + translation.Z

		tr.Vertex3.X = (tx - avgX) / varX * scale + translation.X
		tr.Vertex3.Y = (ty - avgY) / varY * scale + translation.Y
		tr.Vertex3.Z = (tz - avgZ) / varZ * -scale + translation.Z
	}
	/*
	for _, tr := range list {
		fmt.Println(tr)
	}
	fmt.Println(avgX, avgY, avgZ)
	fmt.Println(varX, varY, varZ)
    fmt.Println("list size:", len(list))
	*/
    return list, nil
}

