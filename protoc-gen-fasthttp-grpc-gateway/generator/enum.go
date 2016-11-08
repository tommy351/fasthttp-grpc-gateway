package generator

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

type Enum struct {
	*descriptor.EnumDescriptorProto
	FQMN string
}

func NewEnum(enum *descriptor.EnumDescriptorProto) *Enum {
	return &Enum{
		EnumDescriptorProto: enum,
	}
}
