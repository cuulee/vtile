package polygon

import (
	"encoding/json"
	"fmt"
)

func Pos() []int32 {
	return []int32{0, 0}
}

func get_coords_json(stringcoords string) [][]int32 {
	stringcoords = fmt.Sprintf(`{"coords":%s}`, stringcoords)
	res := ResponseCoords{}
	json.Unmarshal([]byte(stringcoords), &res)

	return res.Coords
}

func moverow(row []int32, geometry []uint32) []uint32 {
	geometry = append(geometry, moveTo(1))
	geometry = append(geometry, uint32(paramEnc(row[0])))
	geometry = append(geometry, uint32(paramEnc(row[1])))
	return geometry

}
func linerow(row []int32, geometry []uint32) []uint32 {
	geometry = append(geometry, uint32(paramEnc(row[0])))
	geometry = append(geometry, uint32(paramEnc(row[1])))
	return geometry

}

type ResponseCoords struct {
	Coords [][]int32 `json:"coords"`
}

func cmdEnc(id uint32, count uint32) uint32 {
	return (id & 0x7) | (count << 3)
}

func moveTo(count uint32) uint32 {
	return cmdEnc(1, count)
}

func lineTo(count uint32) uint32 {
	return cmdEnc(2, count)
}

func closePath(count uint32) uint32 {
	return cmdEnc(7, count)
}

func paramEnc(value int32) int32 {
	return (value << 1) ^ (value >> 31)
}

func Make_Polygon(coords [][]int32, position []int32) ([]uint32, []int32) {
	var count uint32
	count = 0
	var geometry []uint32
	var oldrow []int32
	//total := map[uint32][]int32{}
	//var linetocount uint32
	linetocount := uint32(len(coords) - 2)
	coords = coords[:len(coords)-1]

	for _, row := range coords {
		if count == 0 {
			geometry = moverow([]int32{row[0] - position[0], row[1] - position[1]}, geometry)
			geometry = append(geometry, lineTo(linetocount))

			count = 1
		} else {
			geometry = linerow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
		}
		oldrow = row
	}

	geometry = append(geometry, closePath(1))

	return geometry, oldrow
}

// makes a polygon containing a hole
func Make_Polygon_Hole(coords [][]int32, hole [][]int32, position []int32) ([]uint32, []int32) {
	var count uint32
	count = 0
	var geometry []uint32
	var oldrow []int32
	//total := map[uint32][]int32{}
	//var linetocount uint32
	linetocount := uint32(len(coords) - 2)
	coords = coords[:len(coords)-1]

	for _, row := range coords {
		if count == 0 {
			geometry = moverow([]int32{row[0] - position[0], row[1] - position[1]}, geometry)
			geometry = append(geometry, lineTo(linetocount))

			count = 1
		} else {
			geometry = linerow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
		}
		oldrow = row
	}

	geometry = append(geometry, closePath(1))
	linetocount = uint32(len(hole) - 2)
	hole = hole[:len(hole)-1]
	count = 0
	for _, row := range hole {
		if count == 0 {
			//fmt.Print([]int32{oldrow[0] - row[0], oldrow[1] - row[1]}, row, "pos1\n")
			geometry = moverow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
			geometry = append(geometry, lineTo(linetocount))

			count = 1
		} else {
			geometry = linerow([]int32{row[0] - oldrow[0], row[1] - oldrow[1]}, geometry)
		}
		oldrow = row
	}

	geometry = append(geometry, closePath(1))

	return geometry, oldrow
}
