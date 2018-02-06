package providers

import (
	"log"
)

type TractToTractTransitProvider interface {
	TractToTractProvider
}

type mapsResponse struct {
	routes        []maps.Routes `json:"routes"`
	cheapestRoute maps.Route    `json:"cheapest_route"`
}

type transitProvider struct {
	src    GeoJSON
	dst    GeoJSON
	client *maps.Client
	resp   *mapsResponse
}

func NewTractToTractTransitProvider(src, dst GeoJSON, client *maps.Client) TractToTractTransitProvider {
	return &transitProvider{
		src:    src,
		dst:    dst,
		client: client,
		resp:   nil,
	}
}

func (p *transitProvider) GetType() string {
	return ProviderTypeTransit
}

func (p *transitProvider) GetMetadata() interface{} {
	if p.resp != nil {
		return p.resp
	}
	resp, err := p.fetchData()
	if err != nil {
		return ""
	}
	p.resp = resp
	return p.resp
}

func (p *transitProvider) fetchData() (*mapsResponse, error) {
	origin := fmt.Sprintf("%s,%s", p.src.GetLatitude(), p.src.GetLongitude())
	dest := fmt.Sprintf("%s,%s", p.dst.GetLatitude(), p.dst.GetLongitude())
	req := &maps.DirectionRequest{
		Origin:      origin,
		Destination: dest,
		Mode:        maps.TravelModeTransit,
	}
	routes, _, err := p.client.Directions(context.Background(), req)
	if err != nil {
		log.Printf("google maps err %s", err)
	}
	resp := &mapsResponse{
		routes: routes,
	}
	if len(routes) == 0 {
		log.Println("no routes found")
		return nil, nil
	}
	cheapestRoute := routes[0]
	for i := range routes[1:] {
		route := routes[i]
		if cheapestRoute.Fare.Value > route.Fare.Value {
			cheapestRoute = route
		}
	}
	if cheapestRoute.Fare == nil {
		log.Println("no fares found")
	}
	p.resp.cheapestRoute = cheapestRoute
	return p.resp, nil
}
