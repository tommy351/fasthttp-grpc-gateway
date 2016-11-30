package gateway

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/golang/protobuf/proto"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ResponseMarshaler receives metadata and response and writes them to the request context.
type ResponseMarshaler func(*fasthttp.RequestCtx, *Metadata, interface{})

// ResponseStreamMarshaler receives metadata and response of stream and writes them to the request context.
type ResponseStreamMarshaler func(*fasthttp.RequestCtx, *Metadata, StreamRecvFunc)

// StreamRecvFunc returns protobuf messages.
type StreamRecvFunc func() (proto.Message, error)

// MarshalErrorHandler handles errors of converting JSON requests to protobuf.
type MarshalErrorHandler func(*fasthttp.RequestCtx, error)

// ResponseErrorHandler handles errors of sending requests to gRPC services.
type ResponseErrorHandler func(*fasthttp.RequestCtx, *Metadata, error)

var (
	// GRPCErrorCodes is a map of gRPC error codes and their corresponding HTTP status codes.
	GRPCErrorCodes = map[codes.Code]int{
		codes.OK:                 fasthttp.StatusOK,
		codes.Canceled:           fasthttp.StatusRequestTimeout,
		codes.Unknown:            fasthttp.StatusInternalServerError,
		codes.InvalidArgument:    fasthttp.StatusBadRequest,
		codes.DeadlineExceeded:   fasthttp.StatusRequestTimeout,
		codes.NotFound:           fasthttp.StatusNotFound,
		codes.AlreadyExists:      fasthttp.StatusConflict,
		codes.PermissionDenied:   fasthttp.StatusForbidden,
		codes.Unauthenticated:    fasthttp.StatusUnauthorized,
		codes.ResourceExhausted:  fasthttp.StatusForbidden,
		codes.FailedPrecondition: fasthttp.StatusPreconditionFailed,
		codes.Aborted:            fasthttp.StatusConflict,
		codes.OutOfRange:         fasthttp.StatusBadRequest,
		codes.Unimplemented:      fasthttp.StatusNotImplemented,
		codes.Internal:           fasthttp.StatusInternalServerError,
		codes.Unavailable:        fasthttp.StatusServiceUnavailable,
		codes.DataLoss:           fasthttp.StatusInternalServerError,
	}
)

// Gateway extends fasthttprouter.Router.
type Gateway struct {
	*fasthttprouter.Router

	// Default to PrintJSON.
	ResponseMarshaler ResponseMarshaler

	// Default to PrintJSONStream. You can try PrintJSONStreamArray instead.
	ResponseStreamMarshaler ResponseStreamMarshaler

	// Default to DefaultMarshalErrorHandler.
	MarshalErrorHandler MarshalErrorHandler

	// Default to DefaultResponseErrorHandler.
	ResponseErrorHandler ResponseErrorHandler
}

// NewGateway creates a new gateway with default handlers and router.
func NewGateway() *Gateway {
	router := fasthttprouter.New()

	return &Gateway{
		Router:                  router,
		ResponseMarshaler:       PrintJSON,
		ResponseStreamMarshaler: PrintJSONStream,
		MarshalErrorHandler:     DefaultMarshalErrorHandler,
		ResponseErrorHandler:    DefaultResponseErrorHandler,
	}
}

// DefaultMarshalErrorHandler returns HTTP status code 400 and "bad request" in response body.
func DefaultMarshalErrorHandler(r *fasthttp.RequestCtx, err error) {
	r.SetStatusCode(fasthttp.StatusBadRequest)
	r.WriteString("Bad request")
}

// DefaultResponseErrorHandler returns corresponding HTTP status code of gRPC error and
// error message as the response body.
func DefaultResponseErrorHandler(r *fasthttp.RequestCtx, meta *Metadata, err error) {
	code := grpc.Code(err)
	desc := grpc.ErrorDesc(err)

	if status, ok := GRPCErrorCodes[code]; ok {
		r.SetStatusCode(status)
	} else {
		r.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	r.WriteString(desc)
}
