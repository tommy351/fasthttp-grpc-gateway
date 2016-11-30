package gateway

import "google.golang.org/grpc/metadata"

// Metadata stores headers and trailers of a request.
type Metadata struct {
	Header  metadata.MD
	Trailer metadata.MD
}
