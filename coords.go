package vtile

import (
	"encoding/json"
	"fmt"
)

type ResponseCoords2 struct {
	Coords [][]float64 `json:"coords"`
}

func Get_coords_json2(stringcoords string) [][]float64 {
	stringcoords = fmt.Sprintf(`{"coords":%s}`, stringcoords)
	res := ResponseCoords2{}
	json.Unmarshal([]byte(stringcoords), &res)

	return res.Coords
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

func Make_Coords(coord string, bound Bounds) [][]int32 {
	coords := Get_coords_json2(coord)
	var newlist [][]int32
	for _, i := range coords {
		newlist = append(newlist, single_point(i, bound))
	}
	return newlist
}
