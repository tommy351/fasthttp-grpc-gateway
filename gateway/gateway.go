package gateway

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/golang/protobuf/proto"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ResponseMarshaler func(*fasthttp.RequestCtx, *Metadata, interface{})
type ResponseStreamMarshaler func(*fasthttp.RequestCtx, *Metadata, StreamRecvFunc)
type StreamRecvFunc func() (proto.Message, error)
type ResponseErrorHandler func(*fasthttp.RequestCtx, *Metadata, error)

var (
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

type Gateway struct {
	*fasthttprouter.Router

	ResponseMarshaler       ResponseMarshaler
	ResponseStreamMarshaler ResponseStreamMarshaler
	ResponseErrorHandler    ResponseErrorHandler
}

func NewGateway() *Gateway {
	router := fasthttprouter.New()

	return &Gateway{
		Router:                  router,
		ResponseMarshaler:       PrintJSON,
		ResponseStreamMarshaler: PrintJSONStream,
		ResponseErrorHandler:    DefaultResponseErrorHandler,
	}
}

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
