package dat

import (
	"reflect"
)

type FieldType string

const (
	TypeUint8             FieldType = "u8"
	TypeUint16                      = "u16"
	TypeUint32                      = "u32"
	TypeUint64                      = "u64"
	TypeInt32                       = "i32"
	TypeInt64                       = "i64"
	TypeFloat32                     = "f32"
	TypeFloat64                     = "f64"
	TypeBool                        = "bool"
	TypeString                      = "string"
	TypeListUint8                   = "u8[]"
	TypeListUint16                  = "u16[]"
	TypeListUint32                  = "u32[]"
	TypeListUint64                  = "u64[]"
	TypeListInt32                   = "i32[]"
	TypeListInt64                   = "i64[]"
	TypeListFloat32                 = "f32[]"
	TypeListFloat64                 = "f64[]"
	TypeListBool                    = "bool[]"
	TypeListString                  = "string[]"
	TypeNullableInt32               = "r32"
	TypeNullableInt64               = "r64"
	TypeListNullableInt32           = "r32[]"
	TypeListNullableInt64           = "r64[]"
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
		TypeListBool, TypeListString,
		TypeNullableInt32, TypeNullableInt64,
		TypeListNullableInt32, TypeListNullableInt64:
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
	case TypeInt32, TypeUint32, TypeNullableInt32:
		return 4
	case TypeInt64, TypeUint64, TypeNullableInt64:
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
		TypeListString,
		TypeListNullableInt32, TypeListNullableInt64:
		return 8
	default:
		panic("invalid FieldType")
	}
}

func (ft FieldType) reflectType() reflect.Type {
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
		return reflect.TypeOf([]float32{})
	case TypeListFloat64:
		return reflect.TypeOf([]float64{})
	case TypeListString:
		return reflect.TypeOf([]string{})
	case TypeNullableInt32:
		return reflect.TypeOf((*int32)(nil))
	case TypeNullableInt64:
		return reflect.TypeOf((*int64)(nil))
	case TypeListNullableInt32:
		// FIXME: Implement these as lists of nullable values? These rarely
		// (never?) actually contain null values, but it'd be nice to handle
		// properly
		return reflect.TypeOf([]int32{})
	case TypeListNullableInt64:
		return reflect.TypeOf([]int64{})
	default:
		panic("invalid FieldType")
	}
}

type DataField struct {
	Name   string
	Type   FieldType
	Offset int
}

type DataFormat struct {
	Name          string
	Fields        []DataField
	generatedType *reflect.Type
	size          int
	builtOffsets  bool
}

func (df *DataFormat) buildOffsets() {
	n := 0
	for i := range df.Fields {
		df.Fields[i].Offset = n
		n += df.Fields[i].Type.Size()
	}
	df.size = n
}

func (df *DataFormat) Size() int {
	if !df.builtOffsets {
		df.buildOffsets()
	}
	return df.size
}

func (df *DataFormat) buildType() {
	reflectFields := make([]reflect.StructField, 1+len(df.Fields))
	reflectFields[0] = reflect.StructField{
		Name: "PogoRowID",
		Type: reflect.TypeOf(int(0)),
		Tag:  `json:"_key"`,
	}

	for i := range df.Fields {
		dff := &df.Fields[i]
		reflectFields[i+1] = reflect.StructField{
			Name: dff.Name,
			Type: dff.Type.reflectType(),
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
