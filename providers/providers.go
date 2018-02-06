package providers

const (
	ProviderTypeTransit = "transit"
	ProviderTypeUber    = "uber"
)

type TractToTractProvider interface {
	GetType() string
	GetMetadata() interface{}
}
