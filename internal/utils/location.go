package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Address struct {
	Tourism       string `json:"tourism"`
	HouseNumber   string `json:"house_number"`
	Road          string `json:"road"`
	Suburb        string `json:"suburb"`
	City          string `json:"city"`
	County        string `json:"county"`
	StateDistrict string `json:"state_district"`
	State         string `json:"state"`
	ISO31662Lvl4  string `json:"ISO3166-2-lvl4"`
	Postcode      string `json:"postcode"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
}

type Location struct {
	PlaceID     int64    `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmID       int64    `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	Addresstype string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Address     Address  `json:"address"`
	Boundingbox []string `json:"boundingbox"`
}

func GetLocation(lat string, lon string) (Location, error) {
	// url := "https://nominatim.openstreetmap.org/reverse?lat=22.3459&lon=87.3266&format=json"
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json", lat, lon)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return Location{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return Location{}, err
	}
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return Location{}, err
	}
	var location Location
	err = json.Unmarshal(responseBody, &location)
	if err != nil {
		fmt.Println(err)
		return Location{}, err
	}

	fmt.Printf("%+v\n", location)
	return location, nil
}
