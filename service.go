package shogoa

import "net/http"

// Service is the data structure supporting goa services.
// It provides methods for configuring a service and running it.
// At the basic level a service consists of a set of controllers, each implementing a given
// resource actions. goagen generates global functions - one per resource - that make it
// possible to mount the corresponding controller onto a service. A service contains the
// middleware, not found handler, encoders and muxes shared by all its controllers.
type Service struct {
	// Name of service used for logging, tracing etc.
	Name string

	// TODO
	// Mux ServeMux

	// Server is the service HTTP server.
	Server *http.Server

	// TODO
	// Decoder *HTTPDecoder

	// TODO
	// Encoder *HTTPEncoder
}

func New(name string) *Service {
	// TODO: implement me
	return &Service{
		Name: name,
	}
}
