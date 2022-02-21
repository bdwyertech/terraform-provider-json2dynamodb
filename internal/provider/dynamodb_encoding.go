package provider

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	smithyjson "github.com/aws/smithy-go/encoding/json"
)

func SerializeAttributeMap(v map[string]types.AttributeValue) (jsonBytes []byte, err error) {
	value := smithyjson.NewEncoder()
	object := value.Object()

	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		om := object.Key(key)
		if vv := v[key]; vv == nil {
			continue
		}
		if err = serializeDocumentAttributeValue(v[key], om); err != nil {
			return
		}
	}
	object.Close()
	jsonBytes = value.Bytes()
	return
}

func serializeDocumentAttributeValue(v types.AttributeValue, value smithyjson.Value) error {
	object := value.Object()
	defer object.Close()

	switch uv := v.(type) {
	case *types.AttributeValueMemberB:
		av := object.Key("B")
		av.Base64EncodeBytes(uv.Value)

	case *types.AttributeValueMemberBOOL:
		av := object.Key("BOOL")
		av.Boolean(uv.Value)

	case *types.AttributeValueMemberBS:
		av := object.Key("BS")
		if err := serializeDocumentBinarySetAttributeValue(uv.Value, av); err != nil {
			return err
		}

	case *types.AttributeValueMemberL:
		av := object.Key("L")
		if err := serializeDocumentListAttributeValue(uv.Value, av); err != nil {
			return err
		}

	case *types.AttributeValueMemberM:
		av := object.Key("M")
		if err := serializeDocumentMapAttributeValue(uv.Value, av); err != nil {
			return err
		}

	case *types.AttributeValueMemberN:
		av := object.Key("N")
		av.String(uv.Value)

	case *types.AttributeValueMemberNS:
		av := object.Key("NS")
		if err := serializeDocumentNumberSetAttributeValue(uv.Value, av); err != nil {
			return err
		}

	case *types.AttributeValueMemberNULL:
		av := object.Key("NULL")
		av.Boolean(uv.Value)

	case *types.AttributeValueMemberS:
		av := object.Key("S")
		av.String(uv.Value)

	case *types.AttributeValueMemberSS:
		av := object.Key("SS")
		if err := serializeDocumentStringSetAttributeValue(uv.Value, av); err != nil {
			return err
		}

	default:
		return fmt.Errorf("attempted to serialize unknown member type %T for union %T", uv, v)

	}
	return nil
}

func serializeDocumentBinarySetAttributeValue(v [][]byte, value smithyjson.Value) error {
	array := value.Array()
	defer array.Close()

	for i := range v {
		av := array.Value()
		if vv := v[i]; vv == nil {
			continue
		}
		av.Base64EncodeBytes(v[i])
	}
	return nil
}

func serializeDocumentListAttributeValue(v []types.AttributeValue, value smithyjson.Value) error {
	array := value.Array()
	defer array.Close()

	for i := range v {
		av := array.Value()
		if vv := v[i]; vv == nil {
			continue
		}
		if err := serializeDocumentAttributeValue(v[i], av); err != nil {
			return err
		}
	}
	return nil
}

func serializeDocumentStringSetAttributeValue(v []string, value smithyjson.Value) error {
	array := value.Array()
	defer array.Close()

	for i := range v {
		av := array.Value()
		av.String(v[i])
	}
	return nil
}

func serializeDocumentNumberSetAttributeValue(v []string, value smithyjson.Value) error {
	array := value.Array()
	defer array.Close()

	for i := range v {
		av := array.Value()
		av.String(v[i])
	}
	return nil
}

func serializeDocumentMapAttributeValue(v map[string]types.AttributeValue, value smithyjson.Value) error {
	object := value.Object()
	defer object.Close()

	for key := range v {
		om := object.Key(key)
		if vv := v[key]; vv == nil {
			continue
		}
		if err := serializeDocumentAttributeValue(v[key], om); err != nil {
			return err
		}
	}
	return nil
}
