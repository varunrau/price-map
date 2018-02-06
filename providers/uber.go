package providers

import (
	"log"
	"net/url"
)

type UberProvider interface {
	TractToTractProvider
}

type fareResponse struct {
	Prices []struct {
		DisplayName  string `json:"display_name"`
		HighEstimate int    `json:"high_estiamte"`
		LowEstimate  int    `json:"low_estimate"`
	} `json:"prices"`
}

type uberProvider struct {
	src    GeoJSON
	dst    GeoJSON
	apiKey string
	client *http.Client
	resp   *fareResponse
}

func NewTractToTractUberProvider(src, dst GeoJSON, uberAPIKey string) TractToTractProvider {
	return &uberProvider{
		src:    src,
		dst:    dst,
		apiKey: uberAPIKey,
		resp:   nil,
	}
}

func (p *uberProvider) GetType() string {
	return ProviderTypeUber
}

func (p *uberProvider) GetMetadata() interface{} {

	v := url.Values{}
	v.Set("start_latitude", p.src.GetLatitude())
	v.Set("start_longitude", p.src.GetLongitude())
	v.Set("end_latitude", p.dst.GetLatitude())
	v.Set("end_longitude", p.dst.GetLongitude())
	queryStr := v.Encode()

	resp, err := http.Get(fmt.Sprintf("https://api.uber.com/v1.2/estimates/price/%s", queryStr))
	if err != nil {
		log.Println("error from uber api %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed to read uber response body %s", err)
	}
	r := &fareResponse{}
	err := json.Unmarshal(body, r)
	if err != nil {
		log.Println("failed to unmarshal json %s", err)
	}
	p.resp = r
	return r, nil
}
