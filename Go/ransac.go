//Student Name: Yu-Chen Lee
//Student ID: 300240688

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Point3D struct {
	X float64
	Y float64
	Z float64
}

type Plane3D struct {
	A float64
	B float64
	C float64
	D float64
}

type Plane3DwSupport struct {
	Plane3D
	SupportSize int
}

// the read function uses scanner to get the slice of point3d
func ReadXYZ(filename string) []Point3D {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	var points []Point3D
	scanner := bufio.NewScanner(file)
	scanner.Scan()

	for scanner.Scan() {
		container := strings.Fields(scanner.Text())
		x, err := strconv.ParseFloat(container[0], 64)
		if err != nil {
			return nil
		}
		y, err := strconv.ParseFloat(container[1], 64)
		if err != nil {
			return nil
		}
		z, err := strconv.ParseFloat(container[2], 64)
		if err != nil {
			return nil
		}
		points = append(points, Point3D{x, y, z})
	}

	return points
}

// saves a slice of point into an XYZfile sing bufio.Writer
func SaveXYZ(filename string, points []Point3D) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, p := range points {
		fmt.Fprintf(w, "%f %f %f\n", p.X, p.Y, p.Z)
	}
	w.Flush()
	return nil
}

func (p1 *Point3D) GetDistence(p2 *Point3D) float64 {
	x_part := p1.X - p2.X
	y_part := p1.Y - p2.Y
	z_part := p1.Z - p2.Z

	return math.Sqrt(x_part*x_part + y_part*y_part + z_part*z_part)
}

// the plane formula that calc the plane with three points
func GetPlane(points []Point3D) Plane3D {
	a1 := points[1].X - points[0].X
	b1 := points[1].Y - points[0].Y
	c1 := points[1].Z - points[0].Z

	a2 := points[2].X - points[0].X
	b2 := points[2].Y - points[0].Y
	c2 := points[2].Z - points[0].Z

	a := b1*c2 - b2*c1
	b := a2*c1 - a1*c2
	c := a1*b2 - b1*a2
	d := (-a*points[0].X - b*points[0].Y - c*points[0].Z)

	return Plane3D{a, b, c, d}
}

func GetNumberOfIterations(confidence float64, percentOfPointsOnPlane float64) int {

	i := math.Log(1.0-confidence) / math.Log(1.0-math.Pow(percentOfPointsOnPlane, 3))

	return int(i)
}

// calc the number of supporting points
func GetSupport(plane Plane3D, points []Point3D, eps float64) Plane3DwSupport {
	var supportingPoints []Point3D
	for _, p := range points {
		if math.Abs(plane.A*p.X+plane.B*p.Y+plane.C*p.Z+plane.D) <= eps {
			supportingPoints = append(supportingPoints, p)
		}
	}
	supportSize := len(supportingPoints)
	return Plane3DwSupport{plane, supportSize}
}

// extracts the points that support the given plane if the distance is shorter than eps
func GetSupportingPoints(plane Plane3D, points []Point3D, eps float64) []Point3D {
	var supportingPoints []Point3D
	for _, p := range points {
		if math.Abs(plane.A*p.X+plane.B*p.Y+plane.C*p.Z+plane.D) <= eps {
			supportingPoints = append(supportingPoints, p)
		}
	}
	return supportingPoints
}

// create a new slice of points that  were removed
func RemovePlane(plane Plane3D, points []Point3D, eps float64) []Point3D {
	var restPoints []Point3D
	for _, p := range points {
		if math.Abs(plane.A*p.X+plane.B*p.Y+plane.C*p.Z+plane.D) > eps {
			restPoints = append(restPoints, p)
		}
	}
	return restPoints
}

// the start of piple,using channel to pass the values
func randomPointsGenerator(pc []Point3D) <-chan Point3D {
	out := make(chan Point3D)
	go func() {
		defer close(out)
		for {
			out <- pc[rand.Intn(len(pc))]
		}
	}()
	return out
}

// generate three points to calc the plane from randomPoints
func triPointsGenerator(in <-chan Point3D) <-chan [3]Point3D {
	out := make(chan [3]Point3D)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			var threePoints [3]Point3D
			threePoints[0] = <-in
			threePoints[1] = <-in
			threePoints[2] = <-in

			out <- threePoints
		}
	}()
	return out
}

func takeN(n int, in <-chan [3]Point3D) <-chan [3]Point3D {
	out := make(chan [3]Point3D)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			threePoints := <-in
			out <- threePoints
		}
	}()
	return out
}

func planeEstimator(in <-chan [3]Point3D) <-chan Plane3D {
	out := make(chan Plane3D)
	go func() {
		defer close(out)
		for threePoints := range in {
			plane := GetPlane(threePoints[:])
			out <- plane
		}
	}()
	return out
}

// find all the points that support the dom plane
func supportingPointFinder(in <-chan Plane3D, pc []Point3D, eps float64) <-chan Plane3DwSupport {
	out := make(chan Plane3DwSupport)
	go func() {
		defer close(out)
		for plane := range in {
			out <- GetSupport(plane, pc, eps)
		}
	}()
	return out
}

func fanIn(ins ...<-chan Plane3DwSupport) <-chan Plane3DwSupport {
	out := make(chan Plane3DwSupport)
	var wg sync.WaitGroup
	wg.Add(len(ins))
	for _, in := range ins {
		go func(in <-chan Plane3DwSupport) {
			defer wg.Done()
			for plane := range in {
				out <- plane
			}
		}(in)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// compare the plane to get the most dom plane
func domPlaneIdentifier(in <-chan Plane3DwSupport, bestSupport *Plane3DwSupport) {
	for support := range in {
		if support.SupportSize > bestSupport.SupportSize {
			*bestSupport = support
		}
	}
}

func main() {

	runtime.GOMAXPROCS(10)
	s := time.Now()

	args := os.Args[1:]
	filename := args[0]
	confidence, err := strconv.ParseFloat(args[1], 64)
	percentOfPointsOnPlane, err := strconv.ParseFloat(args[2], 64)
	eps, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		os.Exit(1)
	}
	pc := ReadXYZ(filename)

	for i := 1; i <= 3; i++ {

		var bestSupport Plane3DwSupport

		iterationa := GetNumberOfIterations(confidence, percentOfPointsOnPlane)

		thePc := randomPointsGenerator(pc)
		for iter := 0; iter < iterationa; iter++ {
			threePoints := triPointsGenerator(thePc)
			takeArray := takeN(3, threePoints)
			plane := planeEstimator(takeArray)
			supportPlane := supportingPointFinder(plane, pc, eps)
			fanIn := fanIn(supportPlane)
			go domPlaneIdentifier(fanIn, &bestSupport)

		}

		supportingPoints := GetSupportingPoints(bestSupport.Plane3D, pc, eps)

		resultFile := fmt.Sprintf("%s_p%d", filename, i)
		SaveXYZ(resultFile, supportingPoints)

		pc = RemovePlane(bestSupport.Plane3D, pc, eps)
	}
	removedFile := fmt.Sprintf("%s_p0", filename)
	SaveXYZ(removedFile, pc)

	duration := time.Since(s)
	fmt.Println(duration)
	fmt.Println("-----DONE-----")
}
