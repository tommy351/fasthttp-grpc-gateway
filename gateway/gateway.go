package gateway

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/golang/protobuf/proto"
	"github.com/valyala/fasthttp"
)

type ResponseMarshaler func(*fasthttp.RequestCtx, interface{})
type ResponseStreamMarshaler func(*fasthttp.RequestCtx, StreamRecvFunc)
type StreamRecvFunc func() (proto.Message, error)

type Gateway struct {
	*fasthttprouter.Router

	ResponseMarshaler       ResponseMarshaler
	ResponseStreamMarshaler ResponseStreamMarshaler
}

func NewGateway() *Gateway {
	return &Gateway{
		Router:                  fasthttprouter.New(),
		ResponseMarshaler:       PrintJSON,
		ResponseStreamMarshaler: PrintJSONStream,
	}
}
