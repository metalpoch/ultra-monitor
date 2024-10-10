package osm

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
)

const URL string = "https://nominatim.openstreetmap.org/"

type response struct {
	BoundingBox []string `json:"boundingbox"`
	Address     struct {
		State         string `json:"state"`
		County        string `json:"county"`
		Municipality  string `json:"municipality"`
		Suburb        string `json:"suburb"`
		Neighbourhood string `json:"neighbourhood"`
	} `json:"address"`
}

type location struct {
	State         string    `json:"state"`
	County        string    `json:"county"`
	Municipality  string    `json:"municipality"`
	Suburb        string    `json:"suburb"`
	Neighbourhood string    `json:"neighbourhood"`
	BoundingBox   []float64 `json:"boundingbox"`
}

func getData(url string) (location, error) {
	var result location
	res, err := http.Get(url)
	if err != nil {
		log.Println("error to find coord/osm_id:", err.Error())
		return result, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error to read body response:", err.Error())
		return result, err
	}

	var resp response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Println("error to unmarshal response:", err.Error())
		return result, err
	}

	result = location{
		State:         resp.Address.State,
		County:        resp.Address.County,
		Municipality:  resp.Address.Municipality,
		Suburb:        resp.Address.Suburb,
		Neighbourhood: resp.Address.Neighbourhood,
	}

	for _, s := range resp.BoundingBox {
		f, _ := strconv.ParseFloat(s, 64)
		result.BoundingBox = append(result.BoundingBox, f)
	}
	return result, nil
}

func LocationByCoord(lat, lon float64) (location, error) {
	return getData(fmt.Sprintf("%s/reverse?lat=%f&lon=%f&format=json", URL, lat, lon))
}

func LocationById(id uint) (location, error) {
	return getData(fmt.Sprintf("%s/lookup?osm_ids=%d&format=json", URL, id))
}
