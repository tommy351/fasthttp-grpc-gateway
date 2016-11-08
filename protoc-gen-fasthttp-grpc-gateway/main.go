package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/tommy351/fasthttp-grpc-gateway/protoc-gen-fasthttp-grpc-gateway/generator"
)

func main() {
	req, err := parseReq(os.Stdin)

	if err != nil {
		emitError(err)
		return
	}

	g := generator.New()

	if err := g.Load(req); err != nil {
		emitError(err)
		return
	}

	files, err := g.Generate()

	if err != nil {
		emitError(err)
		return
	}

	emitFiles(files)
}

func parseReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	input, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var req plugin.CodeGeneratorRequest

	if err := proto.Unmarshal(input, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func emitFiles(out []*plugin.CodeGeneratorResponse_File) {
	emitResp(&plugin.CodeGeneratorResponse{
		File: out,
	})
}

func emitError(err error) {
	emitResp(&plugin.CodeGeneratorResponse{
		Error: proto.String(err.Error()),
	})
}

func emitResp(resp *plugin.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stdout.Write(buf); err != nil {
		log.Fatal(err)
	}
}
