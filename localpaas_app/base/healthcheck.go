package base

type HealthcheckType string

const (
	HealthcheckTypeREST HealthcheckType = "REST"
	HealthcheckTypeGRPC HealthcheckType = "gRPC"
)

var (
	AllHealthcheckTypes = []HealthcheckType{HealthcheckTypeREST, HealthcheckTypeGRPC}
)
