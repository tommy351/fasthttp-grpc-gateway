package generator

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

type Method struct {
	*descriptor.MethodDescriptorProto

	RequestType  *Message
	ResponseType *Message
	Bindings     []*Binding
}

func NewMethod(method *descriptor.MethodDescriptorProto) *Method {
	m := &Method{
		MethodDescriptorProto: method,
	}

	return m
}
