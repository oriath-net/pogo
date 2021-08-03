package dat

import (
	"bytes"
	"embed"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
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

func (p *DataParser) Parse(r io.Reader, fileName string) ([]interface{}, error) {
	var err error

	fileName = path.Base(fileName)

	df, dfExists := p.formats[fileName]
	if !dfExists {
		fileBaseName := strings.TrimSuffix(fileName, path.Ext(fileName))
		data, err := fs.ReadFile(p.formatSource, fileBaseName+".xml")
		if err != nil {
			return nil, fmt.Errorf("unable to load format definition for %s: %w", fileName, err)
		}

		df, err = p.typeFromXML(data)
		if err != nil {
			return nil, fmt.Errorf("unable to parse format definition for %s: %w", fileName, err)
		}

		df.width = widthForFilename(fileName)
	}

	if p.debug {
		fmt.Printf("Data file width: %s\n", df.width.String())
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
	df.Size() // make sure offsets are calculated
	fmt.Println("Fields are defined as:")
	for _, f := range df.Fields {
		fmt.Printf("-> %-20s @ %-3d (%s)\n", f.Name, f.Offset, f.Type)
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
	switch typ {
	case TypeUint8, TypeUint16, TypeUint32, TypeUint64,
		TypeInt32, TypeInt64,
		TypeFloat32, TypeFloat64,
		TypeBool:
		return binary.Read(bytes.NewReader(rowdat), binary.LittleEndian, tgt.Addr().Interface())

	case TypeShortID:
		var val uint64

		if ds.rowFormat.width.is64Bit() {
			val = binary.LittleEndian.Uint64(rowdat)
			if val == NullID64 {
				return nil
			}
		} else {
			val = uint64(binary.LittleEndian.Uint32(rowdat))
			if uint32(val) == NullID32 {
				return nil
			}
		}
		tmp := reflect.New(tgt.Type().Elem())
		reflect.Indirect(tmp).SetUint(val)
		tgt.Set(tmp)
		return nil

	case TypeLongID:
		val := binary.LittleEndian.Uint64(rowdat)

		if ds.rowFormat.width.is64Bit() {
			hival := binary.LittleEndian.Uint64(rowdat[8:])
			if hival == NullID64 && val == NullID64 {
				return nil
			}
			if hival != 0 {
				return fmt.Errorf("unexpected value in high half of longid (%016x %016x)", val, hival)
			}
		} else {
			if val == NullID64 {
				return nil
			}
		}
		tmp := reflect.New(tgt.Type().Elem())
		reflect.Indirect(tmp).SetUint(val)
		tgt.Set(tmp)
		return nil

	case TypeListUint8, TypeListUint16,
		TypeListUint32, TypeListUint64,
		TypeListInt32, TypeListInt64,
		TypeListFloat32, TypeListFloat64,
		TypeListBool:
		return ds.readScalarArray(tgt, rowdat, dyndat)

	case TypeString:
		return ds.readString(tgt, rowdat, dyndat)

	case TypeListShortID:
		if ds.rowFormat.width.is64Bit() {
			return ds.readScalarArray(tgt, rowdat, dyndat)
		} else {
			r32 := reflect.New(reflect.TypeOf([]uint32{})).Elem()
			err := ds.readScalarArray(r32, rowdat, dyndat)
			if err != nil {
				return err
			}

			arr32 := r32.Interface().([]uint32)
			n := len(arr32)
			arr64 := make([]uint64, n)
			for i := 0; i < n; i++ {
				arr64[i] = uint64(arr32[i])
			}

			tgt.Set(reflect.ValueOf(arr64))
			return nil
		}

	case TypeListLongID:
		if ds.rowFormat.width.is64Bit() {
			r128 := reflect.New(reflect.TypeOf([][2]uint64{})).Elem()
			err := ds.readScalarArray(r128, rowdat, dyndat)
			if err != nil {
				return err
			}

			arr128 := r128.Interface().([][2]uint64)
			n := len(arr128)
			arr64 := make([]uint64, n)
			for i := 0; i < n; i++ {
				if arr128[i][1] != 0 {
					return fmt.Errorf("unexpected value in high half of longid in array (%016x %016x)", arr128[i][0], arr128[i][1])
				}
				arr64[i] = arr128[i][0]
			}

			tgt.Set(reflect.ValueOf(arr64))
			return nil
		} else {
			return ds.readScalarArray(tgt, rowdat, dyndat)
		}

	case TypeListString:
		return ds.readStringArray(tgt, rowdat, dyndat)

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

func (ds *dataState) readScalarArray(tgt reflect.Value, rowdat []byte, dyndat []byte) error {
	var offset, count int64
	if ds.rowFormat.width.is64Bit() {
		count = int64(binary.LittleEndian.Uint64(rowdat[0:]))
		offset = int64(binary.LittleEndian.Uint64(rowdat[8:]))
	} else {
		count = int64(binary.LittleEndian.Uint32(rowdat[0:]))
		offset = int64(binary.LittleEndian.Uint32(rowdat[4:]))
	}

	if offset < 0 {
		return fmt.Errorf("array offset is negative (%x)", offset)
	}
	if offset > int64(len(dyndat)) {
		return fmt.Errorf("array offset too large (%x)", offset)
	}
	if count < 0 {
		return fmt.Errorf("array count is negative (%x)", count)
	}
	if count > 65535 {
		return fmt.Errorf("array count too large (%x)", count)
	}

	arr := reflect.MakeSlice(tgt.Type(), int(count), int(count))
	err := binary.Read(bytes.NewReader(dyndat[offset:]), binary.LittleEndian, arr.Interface())
	tgt.Set(arr)
	return err
}

func (ds *dataState) readString(tgt reflect.Value, rowdat []byte, dyndat []byte) error {
	var offset int64
	if ds.rowFormat.width.is64Bit() {
		offset = int64(binary.LittleEndian.Uint64(rowdat))
	} else {
		offset = int64(binary.LittleEndian.Uint32(rowdat))
	}

	if offset < 0 {
		return fmt.Errorf("string offset is negative (%x)", offset)
	}
	if offset > int64(len(dyndat)) {
		return fmt.Errorf("string offset too large (%x)", offset)
	}

	str, err := ds.readStringFrom(bytes.NewReader(dyndat[offset:]))
	if err != nil {
		return err
	}
	tgt.SetString(str)
	return nil
}

func (ds *dataState) readStringArray(tgt reflect.Value, rowdat []byte, dyndat []byte) error {
	var offsets reflect.Value
	if ds.rowFormat.width.is64Bit() {
		offsets = reflect.New(reflect.TypeOf([]int64{})).Elem()
	} else {
		offsets = reflect.New(reflect.TypeOf([]int32{})).Elem()
	}
	err := ds.readScalarArray(offsets, rowdat, dyndat)
	if err != nil {
		return fmt.Errorf("unable to read string array offsets: %w", err)
	}

	count := offsets.Len()
	strs := make([]string, count)

	for i := 0; i < count; i++ {
		offset := offsets.Index(i).Int()
		str, err := ds.readStringFrom(bytes.NewReader(dyndat[offset:]))
		if err != nil {
			return err
		}
		strs[i] = str
	}

	tgt.Set(reflect.ValueOf(strs))
	return nil
}

func (ds *dataState) readStringFrom(rr *bytes.Reader) (string, error) {
	if ds.rowFormat.width.isUTF32() {
		str := make([]rune, 0, 32)
		for {
			var ch rune
			err := binary.Read(rr, binary.LittleEndian, &ch)
			if err != nil {
				return "", fmt.Errorf("failed reading string: %w", err)
			}
			if ch == 0 {
				return string(str), nil
			}
			str = append(str, ch)
		}
	} else {
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
}
