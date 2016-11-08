package gateway

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"

	"github.com/valyala/fasthttp"
)

var (
	strTransferEncoding       = []byte("Transfer-Encoding")
	strChunked                = []byte("chunked")
	contentTypeFormURLEncoded = []byte("application/x-www-form-urlencoded")
	contentTypeJSON           = []byte("application/json")

	jsonArrayStart     = []byte("[")
	jsonArrayEnd       = []byte("]")
	jsonArrayDelimiter = []byte(",")
)

func IsBodyURLEncoded(ctx *fasthttp.RequestCtx) bool {
	contentType := ctx.Request.Header.ContentType()
	return bytes.HasPrefix(contentType, contentTypeFormURLEncoded)
}

func IsBodyJSON(ctx *fasthttp.RequestCtx) bool {
	contentType := ctx.Request.Header.ContentType()
	return bytes.HasPrefix(contentType, contentTypeJSON)
}

func PrintJSON(ctx *fasthttp.RequestCtx, body interface{}) {
	encoder := json.NewEncoder(ctx)

	if err := encoder.Encode(body); err != nil {
		panic(err)
	}

	ctx.Response.Header.SetContentTypeBytes(contentTypeJSON)
}

func PrintJSONStream(ctx *fasthttp.RequestCtx, recv StreamRecvFunc) {
	ctx.Response.Header.SetBytesKV(strTransferEncoding, strChunked)
	ctx.Response.Header.SetContentTypeBytes(contentTypeJSON)
	ctx.Response.SetBodyStreamWriter(func(w *bufio.Writer) {
		encoder := json.NewEncoder(w)

		for {
			res, err := recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				panic(err)
			}

			if err := encoder.Encode(res); err != nil {
				panic(err)
			}
		}
	})
}

func PrintJSONStreamArray(ctx *fasthttp.RequestCtx, recv StreamRecvFunc) {
	ctx.Response.Header.SetBytesKV(strTransferEncoding, strChunked)
	ctx.Response.Header.SetContentTypeBytes(contentTypeJSON)
	ctx.Response.SetBodyStreamWriter(func(w *bufio.Writer) {
		encoder := json.NewEncoder(w)
		started := false

		w.Write(jsonArrayStart)

		for {
			res, err := recv()

			if err == io.EOF {
				w.Write(jsonArrayEnd)
				break
			}

			if !started {
				started = true
			} else {
				w.Write(jsonArrayDelimiter)
			}

			if err != nil {
				panic(err)
			}

			if err := encoder.Encode(res); err != nil {
				panic(err)
			}
		}
	})
}
