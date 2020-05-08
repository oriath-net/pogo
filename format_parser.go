package pogo

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type parserState struct {
	fset            *token.FileSet
	typeAliases     map[string]FieldType
	typeDefinitions map[string]DataFormat
}

func (ps *parserState) resolveTypeName(name string) (FieldType, error) {
	aliasTo, aliasExists := ps.typeAliases[name]
	if aliasExists {
		return FieldType(aliasTo), nil
	}

	dt := FieldType(name)
	if dt.Valid() {
		return dt, nil
	}

	return "", fmt.Errorf("unknown type '%s'", name)
}

func (ps *parserState) parseField(f ast.Field) (FieldType, error) {
	flt := f.Type
	var typename string
	isArray := false
	switch flt.(type) {
	case *ast.Ident:
		typename = flt.(*ast.Ident).Name
	case *ast.ArrayType:
		at := flt.(*ast.ArrayType)
		if at.Len != nil {
			return "", fmt.Errorf("array field may not have a length")
		}
		elt, elt_ok := at.Elt.(*ast.Ident)
		if !elt_ok {
			return "", fmt.Errorf("unknown or unsupported field type in array")
		}
		typename = elt.Name
		isArray = true
	default:
		return "", fmt.Errorf("unsupported field type")
	}

	dt, err := ps.resolveTypeName(typename)
	if err != nil {
		return "", err
	}
	if isArray {
		dt = "[]" + dt
	}

	if !dt.Valid() {
		return "", fmt.Errorf("unknown or unsupported field type %s", typename)
	}

	return dt, nil
}

func (ps *parserState) parseStructType(name string, sts *ast.StructType) (DataFormat, error) {
	fields := make([]DataField, 0, len(sts.Fields.List))

	for _, f := range sts.Fields.List {
		if f == nil {
			panic("HOW IS THIS NIL")
		}
		dt, err := ps.parseField(*f)
		if err != nil {
			return DataFormat{}, fmt.Errorf("%w at %s", err, ps.fset.Position(f.Pos()))
		}

		if len(f.Names) == 0 {
			return DataFormat{}, fmt.Errorf("field must have a name at %s", ps.fset.Position(f.Pos()))
		}

		for _, name := range f.Names {
			fields = append(fields, DataField{
				Name: name.Name,
				Type: dt,
			})
		}
	}

	df := DataFormat{
		Name:   name,
		Fields: fields,
	}
	return df, nil
}

func (ps *parserState) parseTypeSpec(ts *ast.TypeSpec) error {
	name := ts.Name.Name
	_, e := ps.typeAliases[name]
	if e {
		return fmt.Errorf("type %s already exists (as alias to %s)", name, ps.typeAliases[name])
	}
	_, e = ps.typeDefinitions[name]
	if e {
		return fmt.Errorf("type %s already exists", name)
	}

	switch ts.Type.(type) {
	case *ast.StructType:
		fmt, err := ps.parseStructType(name, ts.Type.(*ast.StructType))
		if err != nil {
			return err
		}
		ps.typeDefinitions[name] = fmt

	case *ast.Ident:
		id := ts.Type.(*ast.Ident)
		resolvedName, err := ps.resolveTypeName(id.Name)
		if err != nil {
			return err
		}
		ps.typeAliases[name] = resolvedName

	default:
		return fmt.Errorf("unsupported %T type declaration", ts.Type)
	}

	return nil
}

func typesFromFile(path string) ([]DataFormat, error) {
	ps := parserState{
		typeAliases:     make(map[string]FieldType),
		typeDefinitions: make(map[string]DataFormat),
	}

	ps.fset = token.NewFileSet()

	f, err := parser.ParseFile(ps.fset, path, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return []DataFormat{}, err
	}

	_ = ast.NewCommentMap(ps.fset, f, f.Comments)

	for _, d := range f.Decls {
		gd, gd_ok := d.(*ast.GenDecl)
		if !gd_ok {
			return []DataFormat{}, fmt.Errorf("non-type declaration at %s\n", ps.fset.Position(d.Pos()))
		}
		for _, s := range gd.Specs {
			ts, ts_ok := s.(*ast.TypeSpec)
			if !ts_ok {
				return []DataFormat{}, fmt.Errorf("non-type declaration at %s", ps.fset.Position(s.Pos()))
			}
			err := ps.parseTypeSpec(ts)
			if err != nil {
				return []DataFormat{}, fmt.Errorf("%w at %s", err, ps.fset.Position(s.Pos()))
			}
		}
	}

	df := make([]DataFormat, 0, len(ps.typeDefinitions))
	for _, td := range ps.typeDefinitions {
		df = append(df, td)
	}
	return df, nil
}
