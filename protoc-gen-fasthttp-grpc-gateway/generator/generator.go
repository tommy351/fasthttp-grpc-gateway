package generator

import (
	"bytes"
	"go/format"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var (
	baseImports = []string{
		"encoding/json",
		"net/url",
		"golang.org/x/net/context",
		"google.golang.org/grpc",
		"github.com/valyala/fasthttp",
		"github.com/golang/protobuf/proto",
		"github.com/tommy351/fasthttp-grpc-gateway/gateway",
	}
)

type Generator struct {
	files          map[string]*File
	messages       map[string]*Message
	enums          map[string]*Enum
	fileToGenerate []string
}

func New() *Generator {
	return &Generator{
		files:    map[string]*File{},
		messages: map[string]*Message{},
		enums:    map[string]*Enum{},
	}
}

func (g *Generator) Load(req *plugin.CodeGeneratorRequest) error {
	g.fileToGenerate = req.FileToGenerate

	// Load files
	for _, file := range req.ProtoFile {
		f := NewFile(file)
		f.GoPkg = &GoPackage{
			Path: g.goPackagePath(f.FileDescriptorProto),
			Name: defaultGoPackageName(f.FileDescriptorProto),
		}
		g.files[f.GetName()] = f

		// Load messages and enums
		f.LoadMessages(f.MessageType, []string{})
		f.LoadEnums(f.EnumType, []string{})

		for _, msg := range f.Messages {
			g.messages[msg.FQMN] = msg
		}

		for _, enum := range f.Enums {
			g.enums[enum.FQMN] = enum
		}
	}

	// Load services
	for _, file := range g.files {
		if err := file.LoadServices(g, file.Service); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) Generate() ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File

	for _, target := range g.fileToGenerate {
		file, ok := g.files[target]

		if !ok {
			continue
		}

		if len(file.Services) == 0 {
			continue
		}

		content, err := g.generateFile(file)

		if err != nil {
			return nil, err
		}

		formatted, err := format.Source(content)

		if err != nil {
			return nil, err
		}

		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(file.GetOutputFileName()),
			Content: proto.String(string(formatted)),
		})
	}

	return files, nil
}

func (g *Generator) generateFile(file *File) ([]byte, error) {
	tmpl := &gatewayTemplate{
		Generator: g,
		File:      file,
	}

	// Resolve imports
	imports := baseImports

	for _, svc := range file.Services {
		for _, m := range svc.Methods {
			if pkg := m.RequestType.File.GoPkg.Path; pkg != file.GoPkg.Path {
				imports = append(imports, pkg)
			}

			if pkg := m.ResponseType.File.GoPkg.Path; pkg != file.GoPkg.Path {
				imports = append(imports, pkg)
			}
		}
	}

	tmpl.Imports = NewGoPackageList(imports)

	w := bytes.NewBuffer(nil)

	if err := gatewayTmpl.Execute(w, tmpl); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (g *Generator) LookupMessage(name string) *Message {
	msg, ok := g.messages[name]

	if ok {
		return msg
	}

	return nil
}

func (g *Generator) LookupEnum(name string) *Enum {
	enum, ok := g.enums[name]

	if ok {
		return enum
	}

	return nil
}

func (g *Generator) goPackagePath(f *descriptor.FileDescriptorProto) string {
	name := f.GetName()
	gopkg := f.Options.GetGoPackage()
	idx := strings.LastIndex(gopkg, "/")

	if idx >= 0 {
		return gopkg
	}

	return path.Dir(name)
}

func defaultGoPackageName(f *descriptor.FileDescriptorProto) string {
	name := packageIdentityName(f)
	return strings.Replace(name, ".", "_", -1)
}

func packageIdentityName(f *descriptor.FileDescriptorProto) string {
	if f.Options != nil && f.Options.GoPackage != nil {
		gopkg := f.Options.GetGoPackage()
		idx := strings.LastIndex(gopkg, "/")
		if idx < 0 {
			return gopkg
		}

		return gopkg[idx+1:]
	}

	if f.Package == nil {
		base := filepath.Base(f.GetName())
		ext := filepath.Ext(base)
		return strings.TrimSuffix(base, ext)
	}
	return f.GetPackage()
}
