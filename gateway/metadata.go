package gateway

import "google.golang.org/grpc/metadata"

type Metadata struct {
	Header  metadata.MD
	Trailer metadata.MD
}
