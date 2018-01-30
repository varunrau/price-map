package main

import (
	"fmt"
	"log"

	"encoding/json"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"io/ioutil"
)

type CensusTracts struct {
	Meta struct{}  `json:"meta"`
	Data []GeoJSON `json:"data"`
}

type GeoJSON []interface{}

func (j GeoJSON) GetLatitude() string {
	return j[19].(string)
}

func (j GeoJSON) GetLongitude() string {
	return j[20].(string)
}

func main() {

	raw, err := ioutil.ReadFile("./census_tracts/sf.json")
	if err != nil {
		log.Fatalf("failed to read census data %s", err)
	}

	var censusData CensusTracts
	err = json.Unmarshal(raw, &censusData)
	if err != nil {
		log.Fatalf("failed to unmarshal census data %s", err)
	}

	c, err := maps.NewClient(maps.WithAPIKey(""))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	for i := range censusData.Data {
		data := censusData.Data[i+1]
		for j := range censusData.Data[2:] {
			toData := censusData.Data[j+2]
			fmt.Println(fmt.Sprintf("Lat: %s, Long: %s \nLat: %s, Long: %s \n\n", data.GetLatitude(), data.GetLongitude(), toData.GetLatitude(), toData.GetLongitude()))
			getRoute(c, data.GetLatitude(), data.GetLongitude(), toData.GetLatitude(), toData.GetLongitude())
		}
	}
}

func getRoute(c *maps.Client, latitude string, longitude string, lat2 string, long2 string) maps.Route {
	fmt.Println(fmt.Sprintf("%s,%s", latitude, longitude))
	fmt.Println(fmt.Sprintf("%s,%s", lat2, long2))
	r := &maps.DirectionsRequest{
		Origin:      fmt.Sprintf("%s,%s", latitude[1:], longitude),
		Destination: fmt.Sprintf("%s,%s", lat2[1:], long2),
		Mode:        maps.TravelModeTransit,
	}
	routes, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	cheapestRoute := routes[0]
	for i := range routes[1:] {
		route := routes[i]
		if cheapestRoute.Fare.Value > route.Fare.Value {
			cheapestRoute = route
		}
	}
	if cheapestRoute.Fare == nil {
		fmt.Println("no fare info")
	}
	fmt.Println("Cheapest Route:", cheapestRoute.Fare.Text)
	return cheapestRoute
}
