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
	ids    map[int]TileID
}

func Load(filename string) map[int]Line {
	content, err := ioutil.ReadFile("./" + filename)
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

func Make_Line_Geometry(line Line, pos []int32, levelOfDetail float64, sizevalue float64) ([]uint32, []int32) {
	coord := line.coords

	tileid := Get_XY(line.LAT, line.LONG, levelOfDetail)

	bound := TileXYToBounds(tileid)

	coords := Make_Coords(coord, bound, sizevalue)

	if len(coords) == 0 {
		geometry := []uint32{}
		return geometry, pos
	} else {
		geometry, pos := Make_Line(coords, pos)
		return geometry, pos
	}

	//return geometry, pos
}

func Make_Sweep(data map[int]Line, sizes []int) (map[int]Line, map[TileID][]int) {
	var lineval Line
	var dummy TileID
	var sizemap map[int]TileID
	for k := range data {
		lineval = data[k]
		sizemap = map[int]TileID{}
		for _, size := range sizes {
			dummy = Get_XY(lineval.LAT, lineval.LONG, float64(size))
			sizemap[size] = dummy
		}
		lineval.ids = sizemap
		data[k] = lineval
	}
	fmt.Print(data[1])

	sweepmap := map[TileID][]int{}
	for k, v := range data {
		for _, value := range v.ids {
			sweepmap[value] = append(sweepmap[value], k)
		}
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

func create_sizes(line Line, sizes []int) map[int]float64 {
	var lineid TileID
	var mybounds Bounds
	sizemap := map[int]float64{}
	for _, size := range sizes {
		lineid = Get_XY(line.LAT, line.LONG, float64(size))
		mybounds = TileXYToBounds(lineid)
		sizemap[size] = (mybounds.e - mybounds.w) / 4096
	}
	return sizemap
}

func Make_Size_Sweep(sweepmap map[TileID][]int, data map[int]Line, sizemap map[int]float64) {
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
				geometry, pos = Make_Line_Geometry(data[val], pos, float64(k.z), sizemap[int(k.z)])
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

func Make_Line_Tiles(data map[int]Line, sizes []int) {
	// Making initial sweep
	data, sweepmap := Make_Sweep(data, sizes)

	// making sizemap
	sizemap := create_sizes(data[1], sizes)
	fmt.Print(sizemap)

	// Creating Folders and shit
	Creating_Folder(sweepmap)

	// Creating all the files
	Make_Size_Sweep(sweepmap, data, sizemap)

}
