package generator

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	options "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"
)

type NameGetter interface {
	GetName() string
}

type File struct {
	*descriptor.FileDescriptorProto

	GoPkg    *GoPackage
	Messages []*Message
	Enums    []*Enum
	Services []*Service
}

func NewFile(file *descriptor.FileDescriptorProto) *File {
	f := &File{
		FileDescriptorProto: file,
	}

	return f
}

func (f *File) LoadMessages(msgs []*descriptor.DescriptorProto, outers []string) {
	for _, msg := range msgs {
		m := NewMessage(msg)
		m.FQMN = f.getFQMN(msg, outers)
		m.File = f
		f.Messages = append(f.Messages, m)
		nextOuters := append(outers, m.GetName())

		f.LoadMessages(m.NestedType, nextOuters)
		f.LoadEnums(m.EnumType, nextOuters)
	}
}

func (f *File) LoadEnums(enums []*descriptor.EnumDescriptorProto, outers []string) {
	for _, enum := range enums {
		e := NewEnum(enum)
		e.FQMN = f.getFQMN(enum, outers)
		f.Enums = append(f.Enums, e)
	}
}

func (f *File) LoadServices(g *Generator, services []*descriptor.ServiceDescriptorProto) error {
	for _, svc := range services {
		s := NewService(svc)
		f.Services = append(f.Services, s)

		for _, m := range s.Methods {
			m.RequestType = g.LookupMessage(m.GetInputType())
			m.ResponseType = g.LookupMessage(m.GetOutputType())

			opts, err := extractAPIOptions(m.MethodDescriptorProto)

			if err != nil {
				return err
			}

			if opts == nil {
				continue
			}

			b, err := f.newBinding(g, m, opts, 0)

			if err != nil {
				return err
			}

			m.Bindings = append(m.Bindings, b)

			for i, additional := range opts.GetAdditionalBindings() {
				if len(additional.AdditionalBindings) > 0 {
					return fmt.Errorf("additional_binding in additional_binding not allowed: %s.%s", s.GetName(), m.GetName())
				}

				b, err := f.newBinding(g, m, additional, i+1)

				if err != nil {
					return err
				}

				m.Bindings = append(m.Bindings, b)
			}
		}
	}

	return nil
}

func (f *File) newBinding(g *Generator, m *Method, opts *options.HttpRule, idx int) (*Binding, error) {
	var httpMethod HTTPMethod
	var path string
	var bodyType *Message

	switch {
	case opts.GetGet() != "":
		httpMethod = GET
		path = opts.GetGet()

	case opts.GetPost() != "":
		httpMethod = POST
		path = opts.GetPost()

	case opts.GetPut() != "":
		httpMethod = PUT
		path = opts.GetPut()

	case opts.GetPatch() != "":
		httpMethod = PATCH
		path = opts.GetPut()

	case opts.GetDelete() != "":
		httpMethod = DELETE
		path = opts.GetDelete()

	default:
		return nil, fmt.Errorf("none of pattern specified")
	}

	if (httpMethod == GET || httpMethod == DELETE) && opts.Body != "" {
		return nil, fmt.Errorf("request body is not allowed for %s method: %s", httpMethod, m.GetName())
	}

	p, err := NewPath(path)

	if err != nil {
		return nil, err
	}

	if opts.Body == "*" {
		bodyType = m.RequestType
	} else if opts.Body != "" {
		field := m.RequestType.LookupField(opts.Body)

		if field != nil {
			bodyType = g.LookupMessage(field.GetTypeName())
		}

		if bodyType == nil {
			return nil, fmt.Errorf("unable to lookup message type for request body field %s", opts.Body)
		}
	}

	b := &Binding{
		Method:     m,
		HTTPMethod: httpMethod,
		Path:       p,
		Index:      idx,
		Body:       opts.Body,
		BodyType:   bodyType,
	}

	return b, nil
}

func (f *File) getFQMN(namer NameGetter, outers []string) string {
	arr := []string{""}

	if pkg := f.GetPackage(); pkg != "" {
		arr = append(arr, pkg)
	}

	arr = append(arr, outers...)
	arr = append(arr, namer.GetName())

	return strings.Join(arr, ".")
}

func (f *File) GetOutputFileName() string {
	name := f.GetName()
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	return fmt.Sprintf("%s.pb.fgw.go", base)
}

func (f *File) GetPackageName() string {
	return path.Base(path.Dir(f.GetName()))
}

func extractAPIOptions(m *descriptor.MethodDescriptorProto) (*options.HttpRule, error) {
	if m.Options == nil {
		return nil, nil
	}

	if !proto.HasExtension(m.Options, options.E_Http) {
		return nil, nil
	}

	ext, err := proto.GetExtension(m.Options, options.E_Http)

	if err != nil {
		return nil, err
	}

	opts, ok := ext.(*options.HttpRule)

	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}

	return opts, nil
}
