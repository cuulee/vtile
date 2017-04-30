package vtile

import (
	"encoding/json"
	"fmt"
	"math"
)

type ResponseCoords2 struct {
	Coords [][]float64 `json:"coords"`
}

// Point represents a point in space.
type Point struct {
	X float64
	Y float64
}

func Get_coords_json2(stringcoords string) [][]float64 {
	stringcoords = fmt.Sprintf(`{"coords":%s}`, stringcoords)
	res := ResponseCoords2{}
	json.Unmarshal([]byte(stringcoords), &res)

	return res.Coords
}

// Distance finds the length of the hypotenuse between two points.
// Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
func Distance(p1 Point, p2 Point) float64 {
	first := math.Pow(float64(p2.X-p1.X), 2)
	second := math.Pow(float64(p2.Y-p1.Y), 2)
	return math.Sqrt(first + second)
}

func single_point(row []float64, bound Bounds) []int32 {
	deltax := (bound.e - bound.w)
	deltay := (bound.n - bound.s)

	factorx := (row[0] - bound.w) / deltax
	factory := (bound.n - row[1]) / deltay

	xval := int32(factorx * 4096)
	yval := int32(factory * 4096)

	//here1 := uint32((row[0] - bound.w) / (bound.e - bound.w))
	//here2 := uint32((bound.n-row[1])/(bound.n-bound.s)) * 4096
	return []int32{xval, yval}
}

func Make_Coords(coord string, bound Bounds, sizevalue float64) [][]int32 {
	coords := Get_coords_json2(coord)
	var newlist [][]int32
	var pt, oldpt Point
	count := 0
	//var oldi []float64
	oldpt = Point{coords[0][0], coords[0][1]}
	for _, i := range coords {
		pt = Point{i[0], i[1]}
		if count == 0 {
			newlist = append(newlist, single_point(i, bound))
			count = 1
		} else {
			if (Distance(oldpt, pt) > sizevalue) || (len(coords)-1 == count) {
				newlist = append(newlist, single_point(i, bound))
				oldpt = pt
			}
			count += 1
		}

	} //fmt.Print(newlist)

	// if  reduced to one fuck it
	if (len(newlist) == 1) || (len(newlist) == 0) {
		var newslice [][]int32
		//fmt.Print(len(newslice))
		return newslice
	}

	return newlist
}
