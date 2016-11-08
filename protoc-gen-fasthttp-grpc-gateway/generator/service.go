package generator

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

type Service struct {
	*descriptor.ServiceDescriptorProto

	Methods []*Method
}

func NewService(svc *descriptor.ServiceDescriptorProto) *Service {
	s := &Service{
		ServiceDescriptorProto: svc,
		Methods:                []*Method{},
	}

	for _, m := range s.Method {
		s.Methods = append(s.Methods, NewMethod(m))
	}

	return s
}
