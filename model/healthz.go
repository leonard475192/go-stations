package model

// A HealthzResponse expresses health check message.
// if you change thie, you shold run `go vet`
type HealthzResponse struct {
	Message string `json:"message"`
}
