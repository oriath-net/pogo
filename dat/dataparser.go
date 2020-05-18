package dat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"unicode/utf16"
)

const (
	NotTheBs uint64 = 0xbbbb_bbbb_bbbb_bbbb
	NullID64 uint64 = 0xfefe_fefe_fefe_fefe
	NullID32 uint32 = 0xfefe_fefe
)

type DataParser struct {
	formats map[string]DataFormat
}

func InitParser() *DataParser {
	return &DataParser{
		formats: make(map[string]DataFormat),
	}
}

func (dp *DataParser) LoadFormats(path string) error {
	types, err := typesFromFile(path)
	if err != nil {
		return fmt.Errorf("unable to load types from %s: %w", path, err)
	}
	for _, t := range types {
		dp.formats[t.Name] = t
	}
	return nil
}

func (p *DataParser) ParseFile(dataPath string, dataFormat *string) ([]interface{}, error) {
	var f string
	if dataFormat != nil {
		f = *dataFormat
	} else {
		f = strings.TrimSuffix(path.Base(dataPath), ".dat")
	}

	r, err := os.Open(dataPath)
	if err != nil {
		return nil, err
	}

	return p.Parse(r, f)
}

func (p *DataParser) Parse(r io.Reader, formatName string) ([]interface{}, error) {
	df, dfExists := p.formats[formatName]
	if !dfExists {
		return nil, fmt.Errorf("data format '%s' not available", formatName)
	}

	dat, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var rowCount uint32
	err = binary.Read(bytes.NewReader(dat), binary.LittleEndian, &rowCount)
	if err != nil {
		return nil, err
	}

	rowSize := df.Size()
	dynOffset := int64(4 + int(rowCount)*rowSize)

	if int(dynOffset) > len(dat) {
		return nil, fmt.Errorf("row size is too large for data file")
	}

	rowDat := dat[4:dynOffset]
	dynDat := dat[dynOffset:]

	var expectBBBB uint64
	err = binary.Read(bytes.NewReader(dynDat), binary.LittleEndian, &expectBBBB)
	if err != nil {
		return nil, err
	}
	if expectBBBB != NotTheBs {
		return nil, fmt.Errorf("unable to find separator at %x - format specification may be incorrect", dynOffset)
	}

	rowType := df.Type()

	rows := make([]interface{}, rowCount)

	for i := range rows {
		elem := reflect.New(rowType).Elem()

		elem.FieldByName("PogoRowID").SetInt(int64(i))

		rowOffset := rowSize * i
		for j, field := range df.Fields {
			fieldSize := field.Type.Size()
			err := p.readField(elem.Field(j+1), field.Type, rowDat[rowOffset:], dynDat)
			if err != nil {
				return nil, fmt.Errorf("error reading row %d at offset %x: %w", i, rowOffset, err)
			}
			rowOffset += fieldSize
		}

		rows[i] = elem.Interface()
	}

	return rows, err
}

func (p *DataParser) readField(tgt reflect.Value, typ FieldType, rowdat []byte, dyndat []byte) error {
	rr := bytes.NewReader(rowdat)
	switch typ {
	case TypeUint8, TypeUint16, TypeUint32, TypeUint64,
		TypeInt32, TypeInt64,
		TypeFloat32, TypeFloat64,
		TypeBool:
		return p.readScalar(tgt, rr, typ.ReflectType())

	case TypeString:
		return p.readString(tgt, rr, dyndat)

	case TypeListUint8, TypeListUint16,
		TypeListUint32, TypeListUint64,
		TypeListInt32, TypeListInt64,
		TypeListFloat32, TypeListFloat64,
		TypeListBool:
		return p.readScalarArray(tgt, rr, dyndat, typ.ReflectType())

	case TypeListString:
		return p.readStringArray(tgt, rr, dyndat)

	default:
		panic(fmt.Errorf("type '%s' not handled in readField", typ))
	}
}

func (p *DataParser) dynReader(rr *bytes.Reader, dyndat []byte) (*bytes.Reader, error) {
	var off uint32
	err := binary.Read(rr, binary.LittleEndian, &off)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(dyndat[off:]), nil
}

func (p *DataParser) readScalar(tgt reflect.Value, rr *bytes.Reader, typ reflect.Type) error {
	val := reflect.New(typ)
	err := binary.Read(rr, binary.LittleEndian, val.Interface())
	tgt.Set(reflect.Indirect(val))
	return err
}

func (p *DataParser) readScalarArray(tgt reflect.Value, rr *bytes.Reader, dyndat []byte, typ reflect.Type) error {
	var count uint32
	err := binary.Read(rr, binary.LittleEndian, &count)
	if err != nil {
		return err
	}
	dr, err := p.dynReader(rr, dyndat)
	if err != nil {
		return err
	}
	arr := reflect.MakeSlice(typ, int(count), int(count))
	err = binary.Read(dr, binary.LittleEndian, arr.Interface())
	tgt.Set(arr)
	return err
}

func (p *DataParser) readString(tgt reflect.Value, rr *bytes.Reader, dyndat []byte) error {
	dr, err := p.dynReader(rr, dyndat)
	if err != nil {
		return err
	}
	str, err := p.readStringFrom(dr)
	if err != nil {
		return err
	}
	tgt.SetString(str)
	return nil
}

func (p *DataParser) readStringArray(tgt reflect.Value, rr *bytes.Reader, dyndat []byte) error {
	var count uint32
	err := binary.Read(rr, binary.LittleEndian, &count)
	if err != nil {
		return err
	}
	dr, err := p.dynReader(rr, dyndat)
	if err != nil {
		return err
	}
	strOffsets := make([]int32, count)
	err = binary.Read(dr, binary.LittleEndian, strOffsets)
	if err != nil {
		return err
	}
	strs := make([]string, count)
	for i, off := range strOffsets {
		str, err := p.readStringFrom(bytes.NewReader(dyndat[off:]))
		if err != nil {
			return err
		}
		strs[i] = str
	}
	tgt.Set(reflect.ValueOf(strs))
	return nil
}

func (p *DataParser) readStringFrom(rr *bytes.Reader) (string, error) {
	str := make([]uint16, 0, 32)
	for {
		var ch uint16
		err := binary.Read(rr, binary.LittleEndian, &ch)
		if err != nil {
			return "", fmt.Errorf("failed reading string: %w", err)
		}
		if ch == 0 {
			break
		}
		str = append(str, ch)
	}
	return string(utf16.Decode(str)), nil
}
