package openstreetmap

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/goccy/go-json"
)

type OSM struct {
	URL string
}

type location struct {
	State         string   `json:"state"`
	County        string   `json:"county"`
	Municipality  string   `json:"municipality"`
	Suburb        string   `json:"suburb"`
	Neighbourhood string   `json:"neighbourhood"`
	BoundingBox   []string `json:"boundingbox"`
}

func (o OSM) LocationByCoord(lat, lon float64) (location, error) {
	var result location
	url := fmt.Sprintf("%s/?latitude=%f&longitude=%f", o.URL, lat, lon)

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

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("error to unmarshal response:", err.Error())
		return result, err
	}

	if result.State == "" {
		return result, errors.New("error location not found")
	}

	return result, nil
}
