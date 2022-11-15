package spot

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"github.com/SebastiaanKlippert/go-soda"
	"fmt"
)


var dist = .3
const radius = 6371 // Earth's mean radius in kilometers
const token = "Z5hvqVBd8QNF1YqEXNHs5eRtw"
const p_inv = "https://data.lacity.org/resource/s49e-q6j2.json?"
const p_status = "https://data.lacity.org/resource/e7h6-4a3e.json?"

type geo struct {
	lat float64
	long float64
}

type spot struct {
	id string ""
	status string ""
	blockface string ""
	loc geo
}

type block struct {
	blocknum uint16
	st string
}

func get_lat_long(lat string, long string) geo {

	lat_conv, err := strconv.ParseFloat(lat,64)
	if err != nil {
		panic(err)
	}

	long_conv, err := strconv.ParseFloat(long,64)
	if err != nil {
		panic(err)
	}

	return geo{lat: lat_conv, long: long_conv}

}

func degrees2radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func (s *spot) close_to_loc(ug geo) bool{
	
	degreesLat := degrees2radians(s.loc.lat - ug.lat)
 	degreesLong := degrees2radians(s.loc.long - ug.long)
 	a := (math.Sin(degreesLat/2)*math.Sin(degreesLat/2) +
 		math.Cos(degrees2radians(ug.lat))*
 			math.Cos(degrees2radians(s.loc.lat))*math.Sin(degreesLong/2)*
 			math.Sin(degreesLong/2))
 	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
 	d := radius * c
	if d <= dist {
		return true
	} else {
		return false
	}
}

func get_status() []map[string]interface{} {
	sodareq := soda.NewGetRequest(p_status, token)
	return getjson(sodareq)
}

func get_loc(id string) []map[string]interface{}{
	sodareq := soda.NewGetRequest(p_inv, token)
	sodareq.Filters["spaceid"] = id
	return getjson(sodareq)
}

func getjson(sodareq *soda.GetRequest) []map[string]interface{}{
	resp, err := sodareq.Get()

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	results := make([]map[string]interface{}, 0)
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		panic(err)
	}
	return results
}

func sendresults(b block, openspots uint16, w http.ResponseWriter) {
	w.Write([]byte(fmt.Sprintln(b," -> ", openspots, " spots available")))
}

func FindSpots(lat string, long string, w *http.ResponseWriter) {
	origin := get_lat_long(lat,long)

	ch := make(chan uint8,20)
	bl_ch := make(chan string,1000)

	for _, space := range get_status() {
		s := new(spot)
		s.id = space["spaceid"].(string)
		s.status = space["occupancystate"].(string)

		go func(s spot, ch chan uint8, bl_ch chan string, w http.ResponseWriter) {

			results := get_loc(s.id)
			for _, result := range results {
				s.blockface = result["blockface"].(string)
				lat := result["latlng"].(map[string]interface {})["latitude"].(string)
				long := result["latlng"].(map[string]interface {})["longitude"].(string)
				s.loc = get_lat_long(lat,long)
			}
			if s.close_to_loc(origin) {
				if s.status == "VACANT" {
					bl_ch <- s.blockface
				}
			}
			
			<- ch
		}(*s, ch, bl_ch, *w)
		ch <- 1
	}
	
	blks := make(map[block]uint16)
	breakout := false
	for {

		

		select {
			case blk := <-bl_ch:
				blk_parts := strings.Split(blk, " ")
				var b block
				num, err := strconv.ParseInt(blk_parts[0],10, 16)
				if err != nil {
					continue
				} else {
					b.blocknum = uint16(num)
				}

				b.st = strings.Join(blk_parts[1:]," ")

				if _, ok := blks[b]; ok {
					blks[b]++
				} else {
					blks[b] = 1
				}	
			default:
				breakout = true
			}		
		if breakout {
			break
		}
	}
	for i, j := range blks {
		sendresults(i, j, *w)
	}

}
