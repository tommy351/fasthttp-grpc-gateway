package gateway

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"

	"github.com/valyala/fasthttp"
)

const (
	DefaultContextTimeout = 0 * time.Second
	MetadataHeaderPrefix  = "Grpc-Metadata-"
	MetadataTrailerPrefix = "Grpc-Trailer-"
)

var (
	headerGRPCTimeout   = []byte("Grpc-Timeout")
	headerAuthorization = []byte("Authorization")

	strMetadataHeaderPrefix  = []byte(MetadataHeaderPrefix)
	strMetadataTrailerPrefix = []byte(MetadataTrailerPrefix)
	lenMetadataHeaderPrefix  = len(strMetadataHeaderPrefix)
	lenMetadataTrailerPrefix = len(strMetadataTrailerPrefix)
)

func AnnotateContext(ctx context.Context, r *fasthttp.RequestCtx) (context.Context, error) {
	var pairs []string
	timeout := DefaultContextTimeout
	req := &r.Request

	if tm := req.Header.PeekBytes(headerGRPCTimeout); len(tm) > 0 {
		var err error
		timeout, err = decodeDuration(tm)

		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "invalid grpc-timeout: %s", tm)
		}
	}

	req.Header.VisitAll(func(k, v []byte) {
		if bytes.Equal(k, headerAuthorization) {
			pairs = append(pairs, "authorization", string(v))
			return
		}

		if bytes.HasPrefix(k, strMetadataHeaderPrefix) {
			pairs = append(pairs, string(k[lenMetadataHeaderPrefix:]), string(v))
			return
		}
	})

	if timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, timeout)
	}

	if len(pairs) == 0 {
		return ctx, nil
	}

	return metadata.NewContext(ctx, metadata.Pairs(pairs...)), nil
}

func decodeDuration(b []byte) (time.Duration, error) {
	size := len(b)

	if size < 2 {
		return 0, fmt.Errorf("duration string is too short: %s", b)
	}

	d, ok := parseDurationUnit(b[size-1])

	if !ok {
		return 0, fmt.Errorf("invalid duration unit: %s", b)
	}

	t, err := strconv.ParseInt(string(b[:size-1]), 10, 64)

	if err != nil {
		return 0, err
	}

	return d * time.Duration(t), nil
}

func parseDurationUnit(b byte) (time.Duration, bool) {
	switch b {
	case 'H':
		return time.Hour, true
	case 'M':
		return time.Minute, true
	case 'S':
		return time.Second, true
	case 'm':
		return time.Millisecond, true
	case 'u':
		return time.Microsecond, true
	case 'n':
		return time.Nanosecond, true
	}

	return 0, false
}
