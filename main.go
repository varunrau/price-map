package main

import (
	"fmt"
	"log"

	"encoding/json"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	GoogleMaps GoogleMapsConfig `yaml:"google-maps"`
}

type GoogleMapsConfig struct {
	APIKey string `yaml:"api"`
}

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

	rawYaml, err := ioutil.ReadFile("./secrets.yaml")
	if err != nil {
		log.Fatalf("failed to read secrets.yaml %s", err)
	}
	config := Config{}
	yaml.Unmarshal(rawYaml, &config)

	var censusData CensusTracts
	err = json.Unmarshal(raw, &censusData)
	if err != nil {
		log.Fatalf("failed to unmarshal census data %s", err)
	}

	c, err := maps.NewClient(maps.WithAPIKey(config.GoogleMaps.APIKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	for i := range censusData.Data {
		data := censusData.Data[i+1]
		for j := range censusData.Data[2:] {
			toData := censusData.Data[j+2]
		}
	}
}
