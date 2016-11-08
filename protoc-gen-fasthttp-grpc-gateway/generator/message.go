package generator

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

type Message struct {
	*descriptor.DescriptorProto

	FQMN   string
	File   *File
	Fields []*Field
}

func NewMessage(msg *descriptor.DescriptorProto) *Message {
	m := &Message{
		DescriptorProto: msg,
	}

	for _, field := range m.Field {
		m.Fields = append(m.Fields, NewField(field))
	}

	return m
}

func (m *Message) LookupField(name string) *Field {
	for _, field := range m.Fields {
		if field.GetName() == name {
			return field
		}
	}

	return nil
}
