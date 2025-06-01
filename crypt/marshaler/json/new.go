package json

import (
	"encoding/json"
	"reflect"

	gtproto "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/neo532/gokit/crypt/marshaler"
)

func init() {
	marshaler.RegisterMarshaler(NewJson())
}

type opt func(cc *Json)

func WithMarshalEmitUnpopulated(b bool) opt {
	return func(cc *Json) {
		cc.marshalOptions.EmitUnpopulated = b
	}
}

// Codec is a Codec implementation with json.
type Json struct {
	marshalOptions   protojson.MarshalOptions
	unmarshalOptions protojson.UnmarshalOptions
}

func NewJson(opts ...opt) (cc *Json) {
	cc = &Json{
		marshalOptions:   MarshalOptions,
		unmarshalOptions: UnmarshalOptions,
	}
	for _, o := range opts {
		o(cc)
	}
	return
}

func (cc *Json) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return cc.marshalOptions.Marshal(m)
	case gtproto.Message:
		mv := gtproto.MessageV2(m)
		return cc.marshalOptions.Marshal(mv)
	default:
		return json.Marshal(m)
	}
}

func (cc *Json) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return cc.unmarshalOptions.Unmarshal(data, m)
	case gtproto.Message:
		mv := gtproto.MessageV2(m)
		return cc.unmarshalOptions.Unmarshal(data, mv)
	default:

		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return cc.unmarshalOptions.Unmarshal(data, m)
		}
		return json.Unmarshal(data, m)
	}
}

func (cc *Json) Name() string {
	return "json"
}
