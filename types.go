package pogo

import (
	"reflect"
)

type FieldType string

const (
	TypeUint8       FieldType = "uint8"
	TypeUint16                = "uint16"
	TypeUint32                = "uint32"
	TypeUint64                = "uint64"
	TypeInt32                 = "int32"
	TypeInt64                 = "int64"
	TypeFloat32               = "float32"
	TypeFloat64               = "float64"
	TypeBool                  = "bool"
	TypeString                = "string"
	TypeListUint8             = "[]uint8"
	TypeListUint16            = "[]uint16"
	TypeListUint32            = "[]uint32"
	TypeListUint64            = "[]uint64"
	TypeListInt32             = "[]int32"
	TypeListInt64             = "[]int64"
	TypeListFloat32           = "[]float32"
	TypeListFloat64           = "[]float64"
	TypeListBool              = "[]bool"
	TypeListString            = "[]string"
)

func (ft FieldType) Valid() bool {
	switch ft {
	case TypeUint8, TypeUint16,
		TypeUint32, TypeUint64,
		TypeInt32, TypeInt64,
		TypeFloat32, TypeFloat64,
		TypeBool, TypeString,
		TypeListUint8, TypeListUint16,
		TypeListUint32, TypeListUint64,
		TypeListInt32, TypeListInt64,
		TypeListFloat32, TypeListFloat64,
		TypeListBool, TypeListString:
		return true
	default:
		return false
	}
}

func (ft FieldType) Size() int {
	switch ft {
	case TypeUint8:
		return 1
	case TypeUint16:
		return 2
	case TypeUint32:
		return 4
	case TypeUint64:
		return 8
	case TypeInt32:
		return 4
	case TypeInt64:
		return 8
	case TypeFloat32:
		return 4
	case TypeFloat64:
		return 8
	case TypeBool:
		return 1
	case TypeString:
		return 4
	case TypeListUint8, TypeListUint16,
		TypeListUint32, TypeListUint64,
		TypeListInt32, TypeListInt64,
		TypeListFloat32, TypeListFloat64,
		TypeListString:
		return 8
	default:
		panic("invalid FieldType")
	}
}

func (ft FieldType) ReflectType() reflect.Type {
	switch ft {
	case TypeUint8:
		return reflect.TypeOf(uint8(0))
	case TypeUint16:
		return reflect.TypeOf(uint16(0))
	case TypeUint32:
		return reflect.TypeOf(uint32(0))
	case TypeUint64:
		return reflect.TypeOf(uint64(0))
	case TypeInt32:
		return reflect.TypeOf(int32(0))
	case TypeInt64:
		return reflect.TypeOf(int64(0))
	case TypeFloat32:
		return reflect.TypeOf(float32(0))
	case TypeFloat64:
		return reflect.TypeOf(float64(0))
	case TypeBool:
		return reflect.TypeOf(bool(false))
	case TypeString:
		return reflect.TypeOf(string(""))
	case TypeListUint8:
		return reflect.TypeOf([]uint8{})
	case TypeListUint16:
		return reflect.TypeOf([]uint16{})
	case TypeListUint32:
		return reflect.TypeOf([]uint32{})
	case TypeListUint64:
		return reflect.TypeOf([]uint64{})
	case TypeListInt32:
		return reflect.TypeOf([]int32{})
	case TypeListInt64:
		return reflect.TypeOf([]int64{})
	case TypeListFloat32:
		return reflect.TypeOf([]float64{})
	case TypeListFloat64:
		return reflect.TypeOf([]float64{})
	case TypeListString:
		return reflect.TypeOf([]string{})
	default:
		panic("invalid FieldType")
	}
}

type DataField struct {
	Name string
	Type FieldType
}

type DataFormat struct {
	Name          string
	Fields        []DataField
	generatedType *reflect.Type
}

func (df *DataFormat) Size() int {
	sz := 0
	for i := range df.Fields {
		sz += df.Fields[i].Type.Size()
	}
	return sz
}

func (df *DataFormat) buildType() {
	reflectFields := make([]reflect.StructField, 1+len(df.Fields))
	reflectFields[0] = reflect.StructField{
		Name: "PogoRowID",
		Type: reflect.TypeOf(int(0)),
		Tag:  `json:"_key"`,
	}

	for i := range df.Fields {
		reflectFields[i+1] = reflect.StructField{
			Name: df.Fields[i].Name,
			Type: df.Fields[i].Type.ReflectType(),
		}
	}

	t := reflect.StructOf(reflectFields)
	df.generatedType = &t
}

func (df DataFormat) Type() reflect.Type {
	if df.generatedType == nil {
		df.buildType()
	}
	return *df.generatedType
}
