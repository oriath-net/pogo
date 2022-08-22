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
	"path"
	"reflect"
	"strings"
	"unicode/utf16"

	"github.com/oriath-net/pogo/util"
)

const (
	notTheBs  string = "\xbb\xbb\xbb\xbb\xbb\xbb\xbb\xbb"
	nullRow64 uint64 = 0xfefe_fefe_fefe_fefe
	nullRow32 uint32 = 0xfefe_fefe
)

type DataParser struct {
	formatSource fs.FS
	formats      map[string]DataFormat
	debug        int
	strict       int
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

	// debug info
	curRow      int
	curField    string
	lastOffset  int
	seenOffsets map[int]bool
}

//go:embed formats/*.json
var rawEmbeddedFormats embed.FS

var embeddedFormats fs.FS

func init() {
	// we embedded formats/*.json, but we actually just want *.json
	subfs, err := fs.Sub(rawEmbeddedFormats, "formats")
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

func (dp *DataParser) SetDebug(level int) {
	dp.debug = level
}

func (dp *DataParser) SetStrict(level int) {
	dp.strict = level
}

func (dp *DataParser) SetFormatDir(path string) {
	dp.formatSource = os.DirFS(path)
}

func (p *DataParser) getFormat(filename string) (DataFormat, error) {
	df, ok := p.formats[filename]
	if ok {
		return df, nil
	}

	fileBaseName := strings.TrimSuffix(filename, path.Ext(filename))
	f, err := p.formatSource.Open(fileBaseName + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			return DataFormat{}, err
		} else {
			return DataFormat{}, fmt.Errorf("unable to load format definition for %s: %w", filename, err)
		}
	}
	defer f.Close()

	jfmt := JsonFormat{}
	err = util.ReadJson(f, &jfmt)
	if err != nil {
		return DataFormat{}, fmt.Errorf("unable to parse format definition for %s: %w", filename, err)
	}

	df, err = jfmt.BuildType(p)
	if err != nil {
		return DataFormat{}, fmt.Errorf("unable to parse format definition for %s: %w", filename, err)
	}

	df.width = widthForFilename(filename)

	p.formats[filename] = df

	return df, nil
}

func (p *DataParser) Parse(r io.Reader, filename string) ([]any, error) {
	var err error

	filename = path.Base(filename)
	df, err := p.getFormat(filename)
	if err != nil {
		return nil, err
	}

	if p.debug >= 2 {
		p.dumpDataFormat(df)
	}

	dat, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var rowCount32 int32
	err = binary.Read(bytes.NewReader(dat), binary.LittleEndian, &rowCount32)
	if err != nil {
		return nil, err
	}
	rowCount := int(rowCount32)

	rowSize := df.Size()
	dynOffset := 4 + rowCount*rowSize

	if int(dynOffset) > len(dat) {
		return nil, p.debugBoundary(dat, rowSize, rowCount)
	}

	ds := &dataState{
		parser:    p,
		rowType:   df.Type(),
		rowSize:   rowSize,
		rowCount:  rowCount,
		rowFormat: &df,
		rowData:   dat[4:dynOffset],
		dynData:   dat[dynOffset:],
	}

	if p.debug > 0 || p.strict > 0 {
		ds.lastOffset = 8 // boundary occupies 0-7
		ds.seenOffsets = make(map[int]bool)
	}

	if !bytes.Equal(ds.dynData[:8], []byte(notTheBs)) {
		return nil, p.debugBoundary(dat, rowSize, rowCount)
	}

	rows := make([]any, ds.rowCount)
	for i := range rows {
		ds.curRow = i
		if ds.parser.debug >= 2 {
			log.Println("")
			log.Printf("row %d:", i)
		}
		rows[i], err = ds.readRow(i)
		if err != nil {
			return nil, err
		}
	}

	if ds.lastOffset < len(ds.dynData) {
		if ds.parser.debug >= 1 {
			log.Printf("*** last dynamic offset was %x, leaving %d bytes unused", ds.lastOffset, len(ds.dynData)-ds.lastOffset)
		}
		if ds.parser.strict > 0 {
			return nil, fmt.Errorf("%d trailing bytes of dynamic data unused", len(ds.dynData)-ds.lastOffset)
		}
	} else if ds.parser.debug >= 2 {
		log.Printf("*** all dynamic data used")
	}

	return rows, err
}

func (p *DataParser) dumpDataFormat(df DataFormat) {
	size := df.Size() // make sure offsets are calculated
	log.Printf("Data file width: %s", df.width.String())
	log.Printf("Fixed fields are %d bytes:", size)
	for _, f := range df.Fields {
		log.Printf(" -> +%-3x | %-40s %s", f.Offset, f.Name, f.Type)
	}
}

func (p *DataParser) debugBoundary(data []byte, rowSize int, rowCount int) error {
	if rowCount == 0 {
		return fmt.Errorf(
			"format specification inconsistent with data file (zero rows?)",
		)
	}
	boundary := bytes.Index(data, []byte(notTheBs))
	actualRowSize := (boundary - 4) / rowCount
	return fmt.Errorf(
		"format specification inconsistent with data file (spec defines %d bytes/row, file has %d bytes/row)",
		rowSize, actualRowSize,
	)
}

func (ds *dataState) readRow(id int) (any, error) {
	r := reflect.New(ds.rowType).Elem()
	r.FieldByName("PogoRowID").SetInt(int64(id))

	for i, field := range ds.rowFormat.Fields {
		ds.curField = field.Name
		if ds.parser.debug >= 2 {
			log.Printf(" -> %12s %-40s | fix %x", field.Type, field.Name, 4+id*ds.rowSize+field.Offset)
		}
		err := ds.readField(
			r.Field(i+1),
			field.Type,
			ds.rowData[id*ds.rowSize+field.Offset:],
			ds.dynData,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading field %s of row %d: %w", field.Name, id, err)
		}
	}

	return r.Interface(), nil
}

func (ds *dataState) readField(tgt reflect.Value, typ FieldType, rowdat []byte, dyndat []byte) error {
	switch typ {
	case TypeBool:
		val := rowdat[0]
		if ds.parser.strict > 0 && val > 1 {
			return fmt.Errorf("non-boolean value %02x in bool field", val)
		}
		tgt.Set(reflect.ValueOf(val != 0))
		return nil

	case TypeUint8, TypeUint16, TypeUint32, TypeUint64,
		TypeInt32, TypeInt64,
		TypeFloat32, TypeFloat64:
		return binary.Read(bytes.NewReader(rowdat), binary.LittleEndian, tgt.Addr().Interface())

	case TypeShortID:
		var val uint64

		if ds.rowFormat.width.is64Bit() {
			val = binary.LittleEndian.Uint64(rowdat)
			if val == nullRow64 {
				return nil
			}
		} else {
			val = uint64(binary.LittleEndian.Uint32(rowdat))
			if uint32(val) == nullRow32 {
				return nil
			}
		}
		if ds.parser.strict > 0 && val > 0x10_0000 {
			return fmt.Errorf("improbably large shortid %x", val)
		}
		tgt.Set(reflect.ValueOf(&val))
		return nil

	case TypeLongID:
		val := binary.LittleEndian.Uint64(rowdat)

		if ds.rowFormat.width.is64Bit() {
			hival := binary.LittleEndian.Uint64(rowdat[8:])
			if hival == nullRow64 && val == nullRow64 {
				return nil
			}
			if hival != 0 {
				return fmt.Errorf("unexpected value in high half of longid (%016x %016x)", val, hival)
			}
		} else {
			if val == nullRow64 {
				return nil
			}
		}
		if ds.parser.strict > 0 && val > 0x10_0000 {
			return fmt.Errorf("improbably large longid %x", val)
		}
		tgt.Set(reflect.ValueOf(&val))
		return nil

	case TypeListUint8, TypeListUint16,
		TypeListUint32, TypeListUint64,
		TypeListInt32, TypeListInt64,
		TypeListFloat32, TypeListFloat64,
		TypeListBool:
		return ds.readArray(tgt, rowdat, dyndat)

	case TypeString:
		return ds.readString(tgt, rowdat, dyndat)

	case TypeListShortID:
		if ds.rowFormat.width.is64Bit() {
			return ds.readArray(tgt, rowdat, dyndat)
		} else {
			r32 := reflect.New(reflect.TypeOf([]uint32{})).Elem()
			err := ds.readArray(r32, rowdat, dyndat)
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
			err := ds.readArray(r128, rowdat, dyndat)
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
			return ds.readArray(tgt, rowdat, dyndat)
		}

	case TypeListString:
		return ds.readStringArray(tgt, rowdat, dyndat)

	case TypeListVoid:
		var count int64
		if ds.rowFormat.width.is64Bit() {
			count = int64(binary.LittleEndian.Uint64(rowdat[0:]))
		} else {
			count = int64(binary.LittleEndian.Uint32(rowdat[0:]))
		}
		if count > 0 && ds.parser.strict > 0 {
			return fmt.Errorf("non-empty void[]")
		}
		// gaze into the abyss
		voidArray := make([]any, count)
		tgt.Set(reflect.ValueOf(voidArray))
		return ds.usedDyndat("array", ds.lastOffset, 0, int(count))

	default:
		panic(fmt.Errorf("type '%s' not handled in readField", typ))
	}
}

func (ds *dataState) rawReadArray(tgt reflect.Value, rowdat []byte, dyndat []byte) (int, int, int, error) {
	var offset, count int64
	if ds.rowFormat.width.is64Bit() {
		count = int64(binary.LittleEndian.Uint64(rowdat[0:]))
		offset = int64(binary.LittleEndian.Uint64(rowdat[8:]))
	} else {
		count = int64(binary.LittleEndian.Uint32(rowdat[0:]))
		offset = int64(binary.LittleEndian.Uint32(rowdat[4:]))
	}

	if offset < 8 {
		return 0, 0, 0, fmt.Errorf("array offset too low (%x)", offset)
	}
	if offset > int64(len(dyndat)) {
		return 0, 0, 0, fmt.Errorf("array offset too large (%x)", offset)
	}
	if count < 0 {
		return 0, 0, 0, fmt.Errorf("array count negative (%x)", count)
	}
	if count > 65535 {
		return 0, 0, 0, fmt.Errorf("array count too large (%x)", count)
	}

	rdr := bytes.NewReader(dyndat[offset:])
	arr := reflect.MakeSlice(tgt.Type(), int(count), int(count))
	err := binary.Read(rdr, binary.LittleEndian, arr.Interface())
	if err != nil {
		return 0, 0, 0, err
	}
	tgt.Set(arr)

	readLength, _ := rdr.Seek(0, io.SeekCurrent)
	return int(offset), int(readLength), int(count), nil
}

func (ds *dataState) readArray(tgt reflect.Value, rowdat []byte, dyndat []byte) error {
	offset, length, count, err := ds.rawReadArray(tgt, rowdat, dyndat)
	if err != nil {
		return err
	}

	return ds.usedDyndat("array", offset, length, count)
}

func (ds *dataState) readString(tgt reflect.Value, rowdat []byte, dyndat []byte) error {
	var offset int64
	if ds.rowFormat.width.is64Bit() {
		offset = int64(binary.LittleEndian.Uint64(rowdat))
	} else {
		offset = int64(binary.LittleEndian.Uint32(rowdat))
	}

	str, err := ds.readStringFrom(dyndat, int(offset))
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

	offsetBase, offsetLength, offsetCount, err := ds.rawReadArray(offsets, rowdat, dyndat)
	if err != nil {
		return fmt.Errorf("unable to read string array offsets: %w", err)
	}

	count := offsets.Len()
	strs := make([]string, count)

	for i := 0; i < count; i++ {
		offset := offsets.Index(i).Int()
		str, err := ds.readStringFrom(dyndat, int(offset))
		if err != nil {
			return err
		}
		strs[i] = str
	}

	tgt.Set(reflect.ValueOf(strs))

	return ds.usedDyndat("offsets", offsetBase, offsetLength, offsetCount)
}

func (ds *dataState) readStringFrom(dyndat []byte, offset int) (string, error) {
	if offset < 8 {
		return "", fmt.Errorf("string offset too low (%x)", offset)
	}
	if offset > len(dyndat) {
		return "", fmt.Errorf("string offset too large (%x)", offset)
	}

	origOffset := offset
	if ds.rowFormat.width.isUTF32() {
		str := make([]rune, 0, 32)
		for {
			if offset+4 > len(dyndat) {
				break
			}
			ch := rune(binary.LittleEndian.Uint32(dyndat[offset:]))
			offset += 4
			if ch == 0 {
				return string(str), ds.usedDyndat("string", origOffset, offset-origOffset, 0)
			}
			str = append(str, ch)
		}
	} else {
		str := make([]uint16, 0, 32)
		for {
			if offset+2 > len(dyndat) {
				break
			}
			ch := binary.LittleEndian.Uint16(dyndat[offset:])
			offset += 2
			if ch == 0 {
				return string(utf16.Decode(str)), ds.usedDyndat("string", origOffset, offset-origOffset+2, 0) // +2? yep
			}
			str = append(str, ch)
		}
	}

	return "", io.EOF
}

func (ds *dataState) usedDyndat(purpose string, offset int, length int, count int) error {
	if ds.parser.debug == 0 && ds.parser.strict == 0 {
		return nil
	}
	message := ""
	warning := false
	if offset < ds.lastOffset {
		_, seen := ds.seenOffsets[offset]
		if seen {
			message = "(reused)"
		} else {
			message = "offset went backwards"
			warning = true
		}
	} else {
		if offset > ds.lastOffset {
			message = fmt.Sprintf("skipped %d bytes", offset-ds.lastOffset)
			warning = true
		}
		ds.lastOffset = offset + length
	}
	if purpose == "array" {
		purpose += fmt.Sprintf("[%d]", count)
	}
	if warning && ds.parser.strict > 0 {
		return fmt.Errorf("%s before %s %s", message, purpose, ds.curField)
	}
	if ds.parser.debug >= 2 {
		if warning {
			message = strings.ToUpper(message)
		}
		log.Printf(" ---> %10s %-10s @ dyn %x + %x -> %x %s", "", purpose, offset, length, offset+length, message)
	} else if warning { // implies ds.parser.debug > 0
		log.Printf("*** Row %d, %s %s, at %x + %x: %s", ds.curRow, purpose, ds.curField, offset, length, message)
	}
	ds.seenOffsets[offset] = true
	return nil
}
