package generator

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

type Field struct {
	*descriptor.FieldDescriptorProto
}

func NewField(field *descriptor.FieldDescriptorProto) *Field {
	return &Field{
		FieldDescriptorProto: field,
	}
}
