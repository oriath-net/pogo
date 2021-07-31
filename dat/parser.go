package dat

import (
	"bytes"
	"embed"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
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
	formatSource fs.FS
	formats      map[string]DataFormat
	debug        bool
	version      string
}

type dataState struct {
	parser    *DataParser
	rowType   reflect.Type
	rowFormat *DataFormat
	rowSize   int
	rowCount  int
	rowData   []byte
	dynData   []byte
}

//go:embed formats/xml/*.xml
var rawEmbeddedFormats embed.FS

var embeddedFormats fs.FS

func init() {
	// we embedded formats/xml/*.xml, but we actually just want *.xml
	subfs, err := fs.Sub(rawEmbeddedFormats, "formats/xml")
	if err != nil {
		panic(err)
	}
	embeddedFormats = subfs
}

func InitParser(version string) *DataParser {
	return &DataParser{
		formatSource: embeddedFormats,
		formats:      make(map[string]DataFormat),
		version:      version,
	}
}

func (dp *DataParser) EnableDebug() {
	dp.debug = true
}

func (dp *DataParser) SetFormatDir(path string) {
	dp.formatSource = os.DirFS(path)
}

func (p *DataParser) Parse(r io.Reader, formatName string) ([]interface{}, error) {
	var err error

	df, dfExists := p.formats[formatName]
	if !dfExists {
		data, err := fs.ReadFile(p.formatSource, formatName+".xml")
		if err != nil {
			return nil, err // FIXME wrap
		}

		df, err = p.typeFromXML(data)
		if err != nil {
			return nil, err // FIXME wrap
		}
	}

	if p.debug {
		p.dumpDataFormat(df)
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
		return nil, fmt.Errorf("row size (%d bytes) is larger than data file", rowSize)
	}

	ds := &dataState{
		parser:    p,
		rowType:   df.Type(),
		rowSize:   rowSize,
		rowCount:  int(rowCount),
		rowFormat: &df,
		rowData:   dat[4:dynOffset],
		dynData:   dat[dynOffset:],
	}

	err = ds.checkBoundary(dat)
	if err != nil {
		return nil, err
	}

	rows := make([]interface{}, ds.rowCount)
	for i := range rows {
		rows[i], err = ds.readRow(i)
		if err != nil {
			return nil, err
		}
	}

	return rows, err
}

func (p *DataParser) dumpDataFormat(df DataFormat) {
	log.Println("Fields are defined as:")
	for _, f := range df.Fields {
		log.Printf("-> %-20s @ %-3d (%s)\n", f.Name, f.Offset, f.Type)
	}
}

func (ds *dataState) checkBoundary(data []byte) error {
	var expectBBBB uint64
	err := binary.Read(bytes.NewReader(ds.dynData), binary.LittleEndian, &expectBBBB)
	if err != nil {
		return err
	}
	if expectBBBB != NotTheBs {
		if ds.rowCount == 0 {
			return fmt.Errorf(
				"format specification inconsistent with data file (zero rows?)",
			)
		}
		boundary := strings.Index(string(data), "\xbb\xbb\xbb\xbb\xbb\xbb\xbb\xbb")
		actualRowSize := (boundary - 4) / ds.rowCount
		return fmt.Errorf(
			"format specification inconsistent with data file (spec defines %d bytes/row, file has %d bytes/row)",
			ds.rowSize, actualRowSize,
		)
	}
	return nil
}

func (ds *dataState) readRow(id int) (interface{}, error) {
	r := reflect.New(ds.rowType).Elem()
	r.FieldByName("PogoRowID").SetInt(int64(id))

	for i, field := range ds.rowFormat.Fields {
		err := ds.readField(
			r.Field(i+1),
			field.Type,
			ds.rowData[id*ds.rowSize+field.Offset:],
			ds.dynData,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading field %s of row %d: %w", field.Name, i, err)
		}
	}

	return r.Interface(), nil
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
		var nullval uint64
		switch typ {
		case TypeNullableInt32:
			x := NullID32
			nullval = uint64(int32(x))
		case TypeNullableInt64:
			nullval = NullID64
		}
		if uint64(reflect.Indirect(tmp).Int()) != nullval {
			tgt.Set(tmp)
		}
		return nil

	case TypeString:
		return ds.readString(tgt, rr, dyndat)

	case TypeListUint8, TypeListUint16,
		TypeListUint32, TypeListUint64,
		TypeListInt32, TypeListInt64,
		TypeListFloat32, TypeListFloat64,
		TypeListNullableInt32, TypeListNullableInt64,
		TypeListBool:
		return ds.readScalarArray(tgt, rr, dyndat)

	case TypeListString:
		return ds.readStringArray(tgt, rr, dyndat)

	default:
		panic(fmt.Errorf("type '%s' not handled in readField", typ))
	}
}

func (ds *dataState) dynReader(rr *bytes.Reader, dyndat []byte) (*bytes.Reader, error) {
	var off int32
	err := binary.Read(rr, binary.LittleEndian, &off)
	if err != nil {
		return nil, err
	}
	if off < 0 || int(off) > len(dyndat) {
		return nil, fmt.Errorf("invalid offset to dynamic data (%08x)", uint32(off))
	}
	return bytes.NewReader(dyndat[off:]), nil
}

func (ds *dataState) readScalarArray(tgt reflect.Value, rr *bytes.Reader, dyndat []byte) error {
	var count int32
	err := binary.Read(rr, binary.LittleEndian, &count)
	if err != nil {
		return err
	}
	dr, err := ds.dynReader(rr, dyndat)
	if err != nil {
		return err
	}
	if count < 0 {
		return fmt.Errorf("array length was negative")
	}
	arr := reflect.MakeSlice(tgt.Type(), int(count), int(count))
	err = binary.Read(dr, binary.LittleEndian, arr.Interface())
	tgt.Set(arr)
	return err
}

func (ds *dataState) readString(tgt reflect.Value, rr *bytes.Reader, dyndat []byte) error {
	dr, err := ds.dynReader(rr, dyndat)
	if err != nil {
		return err
	}
	str, err := ds.readStringFrom(dr)
	if err != nil {
		return err
	}
	tgt.SetString(str)
	return nil
}

func (ds *dataState) readStringArray(tgt reflect.Value, rr *bytes.Reader, dyndat []byte) error {
	var count int32
	err := binary.Read(rr, binary.LittleEndian, &count)
	if err != nil {
		return err
	}
	dr, err := ds.dynReader(rr, dyndat)
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
		str, err := ds.readStringFrom(bytes.NewReader(dyndat[off:]))
		if err != nil {
			return err
		}
		strs[i] = str
	}
	tgt.Set(reflect.ValueOf(strs))
	return nil
}

func (ds *dataState) readStringFrom(rr *bytes.Reader) (string, error) {
	str := make([]uint16, 0, 32)
	for {
		var ch uint16
		err := binary.Read(rr, binary.LittleEndian, &ch)
		if err != nil {
			return "", fmt.Errorf("failed reading string: %w", err)
		}
		if ch == 0 {
			return string(utf16.Decode(str)), nil
		}
		str = append(str, ch)
	}
}
