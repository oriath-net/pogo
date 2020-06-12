package dat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"unicode/utf16"
)

const (
	NotTheBs uint64 = 0xbbbb_bbbb_bbbb_bbbb
	NullID64 int64  = -0x101010101010102 // unsigned: fefefefe_fefefefe
	NullID32 int32  = -0x1010102         // unsigned: fefefefe
)

type DataParser struct {
	formats map[string]DataFormat
	debug   bool
}

type dataState struct {
	parser     *DataParser
	lastOffset int
	rowID      int
	currField  string
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

func (dp *DataParser) EnableDebug() {
	dp.debug = true
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

	var rowCount int32
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

	ds := &dataState{
		parser:     p,
		lastOffset: 8,
	}

	for i := range rows {
		elem := reflect.New(rowType).Elem()

		elem.FieldByName("PogoRowID").SetInt(int64(i))

		rowOffset := rowSize * i
		for j, field := range df.Fields {
			if p.debug {
				ds.rowID = i
				ds.currField = field.Name
			}
			fieldSize := field.Type.Size()
			err := ds.readField(elem.Field(j+1), field.Type, rowDat[rowOffset:], dynDat)
			if err != nil {
				return nil, fmt.Errorf("error reading row %d at offset %x: %w", i, rowOffset, err)
			}
			rowOffset += fieldSize
		}

		rows[i] = elem.Interface()
	}

	if p.debug {
		if ds.lastOffset == len(dynDat) {
			log.Printf("All dynamic data consumed")
		} else {
			log.Printf("%d bytes of unused dynamic data starting at %06x", len(dynDat)-ds.lastOffset, ds.lastOffset)
		}
	}

	return rows, err
}

func (ds *dataState) readField(tgt reflect.Value, typ FieldType, rowdat []byte, dyndat []byte) error {
	rr := bytes.NewReader(rowdat)
	switch typ {
	case TypeUint8, TypeUint16, TypeUint32, TypeUint64,
		TypeInt32, TypeInt64,
		TypeFloat32, TypeFloat64,
		TypeBool:
		return binary.Read(rr, binary.LittleEndian, tgt.Addr().Interface())

	case TypeNullableInt32, TypeNullableInt64:
		tmp := reflect.New(tgt.Type().Elem())
		err := binary.Read(rr, binary.LittleEndian, tmp.Interface())
		if err != nil {
			return err
		}
		var nullval int64
		switch typ {
		case TypeNullableInt32:
			nullval = int64(NullID32)
		case TypeNullableInt64:
			nullval = int64(NullID64)
		}
		if reflect.Indirect(tmp).Int() != nullval {
			tgt.Set(tmp)
		}
		return nil

	case TypeString:
		return ds.readString(tgt, rr, dyndat)

	case TypeListUint8, TypeListUint16,
		TypeListUint32, TypeListUint64,
		TypeListInt32, TypeListInt64,
		TypeListFloat32, TypeListFloat64,
		TypeListBool:
		return ds.readScalarArray(tgt, rr, dyndat)

	case TypeListString:
		return ds.readStringArray(tgt, rr, dyndat)

	default:
		panic(fmt.Errorf("type '%s' not handled in readField", typ))
	}
}

func (ds *dataState) dynReader(rr *bytes.Reader, dyndat []byte) (*bytes.Reader, int, error) {
	var off int32
	err := binary.Read(rr, binary.LittleEndian, &off)
	if err != nil {
		return nil, 0, err
	}
	return bytes.NewReader(dyndat[off:]), int(off), nil
}

func (ds *dataState) readScalarArray(tgt reflect.Value, rr *bytes.Reader, dyndat []byte) error {
	var count int32
	err := binary.Read(rr, binary.LittleEndian, &count)
	if err != nil {
		return err
	}
	dr, off, err := ds.dynReader(rr, dyndat)
	if err != nil {
		return err
	}
	arr := reflect.MakeSlice(tgt.Type(), int(count), int(count))
	err = binary.Read(dr, binary.LittleEndian, arr.Interface())
	tgt.Set(arr)
	if ds.parser.debug {
		ds.debugOffset(int(off), int(count)*int(tgt.Type().Elem().Size()), "before")
	}
	return err
}

func (ds *dataState) readString(tgt reflect.Value, rr *bytes.Reader, dyndat []byte) error {
	dr, off, err := ds.dynReader(rr, dyndat)
	if err != nil {
		return err
	}
	str, count, err := ds.readStringFrom(dr)
	if err != nil {
		return err
	}
	tgt.SetString(str)

	if ds.parser.debug {
		ds.debugOffset(off, count, "before")
	}
	return nil
}

func (ds *dataState) readStringArray(tgt reflect.Value, rr *bytes.Reader, dyndat []byte) error {
	var count int32
	err := binary.Read(rr, binary.LittleEndian, &count)
	if err != nil {
		return err
	}
	dr, offtab_off, err := ds.dynReader(rr, dyndat)
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
		str, str_count, err := ds.readStringFrom(bytes.NewReader(dyndat[off:]))
		if err != nil {
			return err
		}
		if ds.parser.debug {
			ds.debugOffset(int(off), str_count, "before strings of")
		}
		strs[i] = str
	}
	tgt.Set(reflect.ValueOf(strs))
	if ds.parser.debug {
		ds.debugOffset(offtab_off, 4*int(count), "before offset table of")
	}
	return nil
}

func (ds *dataState) readStringFrom(rr *bytes.Reader) (string, int, error) {
	str := make([]uint16, 0, 32)
	for {
		var ch uint16
		err := binary.Read(rr, binary.LittleEndian, &ch)
		if err != nil {
			return "", 0, fmt.Errorf("failed reading string: %w", err)
		}
		if ch == 0 {
			return string(utf16.Decode(str)), 2*len(str) + 4, nil
		}
		str = append(str, ch)
	}
}

func (ds *dataState) debugOffset(off int, increment int, context string) {
	if off < ds.lastOffset {
		log.Printf("data offset MOVED BACKWARDS from %x to %x %s field %s of row %d", ds.lastOffset, off, context, ds.currField, ds.rowID)
	} else if off > ds.lastOffset {
		log.Printf("data offset skipped %d bytes from %x to %x %s field %s of row %d", off-ds.lastOffset, ds.lastOffset, off, context, ds.currField, ds.rowID)
	}
	ds.lastOffset = off + increment
}
