package vtile

import (
	"encoding/csv"
	"fmt"

	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"vector-tile/2.1"
)

// Tree a struct holder for tree information.
type Line struct {
	LAT    float64
	LONG   float64
	gid    string
	coords string
	id13   TileID
}

func Load() map[int]Line {
	content, err := ioutil.ReadFile("./tree_lines.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(strings.NewReader(string(content[:])))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// TreeID,qLegalStatus,qSpecies,qAddress,SiteOrder,qSiteInfo,PlantType,qCaretaker,qCareAssistant,PlantDate,DBH,PlotSize,PermitNotes,XCoord,YCoord,Latitude,Longitude,Location
	trees := map[int]Line{}
	for i, record := range records[1:] {
		lat, _ := strconv.ParseFloat(record[2], 64)
		lng, _ := strconv.ParseFloat(record[3], 64)
		gid := record[0]
		coords := record[1]
		trees[i] = Line{LAT: lat, LONG: lng, gid: gid, coords: coords}
	}
	return trees
}

func Make_Line_Geometry(line Line, pos []int32, levelOfDetail float64) ([]uint32, []int32) {
	coord := line.coords

	tileid := Get_XY(line.LAT, line.LONG, levelOfDetail)

	bound := TileXYToBounds(tileid)

	coords := Make_Coords(coord, bound)

	geometry, pos := Make_Line(coords, pos)

	return geometry, pos
}

func Make_Sweep(data map[int]Line, size int) (map[int]Line, map[TileID][]int) {
	var lineval Line
	var dummy TileID
	for k := range data {
		lineval = data[k]
		dummy = Get_XY(lineval.LAT, lineval.LONG, float64(size))
		lineval.id13 = dummy
		data[k] = lineval
	}
	sweepmap := map[TileID][]int{}
	for k, v := range data {
		sweepmap[v.id13] = append(sweepmap[v.id13], k)
	}

	return data, sweepmap
}

func Creating_Folder(sweepmap map[TileID][]int) {
	//var p string
	//var total []string
	for k := range sweepmap {
		zval := strconv.Itoa(int(k.z))
		xval := strconv.Itoa(int(k.x))
		//p = fmt.Sprintf("/tiles/%s/%s", zval, xval)
		///total = append(total, p)
		os.MkdirAll("tiles/"+zval+"/"+xval, os.ModePerm)
	}
	//fmt.Print(strings.Join(total, ","))
	//exec.Command("here", strings.Join(total, ",")).Run()
}

func createTileWithLine(xyz TileID, total []uint32) ([]byte, error) {
	tile := &vector_tile.Tile{}
	var layerVersion = vector_tile.Default_Tile_Layer_Version
	layerName := "lines"
	featureType := vector_tile.Tile_LINESTRING
	var extent = vector_tile.Default_Tile_Layer_Extent
	//var bound []Bounds

	tile.Layers = []*vector_tile.Tile_Layer{
		{
			Version: &layerVersion,
			Name:    &layerName,
			Extent:  &extent,
			Features: []*vector_tile.Tile_Feature{
				{
					Tags:     []uint32{},
					Type:     &featureType,
					Geometry: total,
				},
			},
		},
	}

	return proto.Marshal(tile)
}

func Make_Size_Sweep(sweepmap map[TileID][]int, data map[int]Line) {
	c := make(chan string)

	for k, v := range sweepmap {
		go func(k TileID, v []int, data map[int]Line, c chan<- string) {
			zval := strconv.Itoa(int(k.z))
			xval := strconv.Itoa(int(k.x))
			yval := strconv.Itoa(int(k.y))
			filename := fmt.Sprintf("tiles/%s/%s/%s", zval, xval, yval)
			//c := make(chan []uint32)
			pos := Pos()
			total := []uint32{}
			var geometry []uint32
			for _, val := range v {
				geometry, pos = Make_Line_Geometry(data[val], pos, float64(13))
				total = append(total, geometry...)
			}
			pbfdata, _ := createTileWithLine(k, total)
			ioutil.WriteFile(filename, []byte(pbfdata), 0644)
			c <- filename
		}(k, v, data, c)

	}

	counter := 0
	for counter < len(sweepmap) {
		select {
		case msg1 := <-c:
			fmt.Print("\n" + msg1)

		}
		counter += 1
	}
}
