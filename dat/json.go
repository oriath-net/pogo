package dat

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/mcuadros/go-version"
)

type JsonFormat struct {
	File        string      `json:"file"`
	Fields      []JsonField `json:"fields"`
	Description string      `json:"description,omitempty"`
	Since       string      `json:"since,omitempty"`
	Until       string      `json:"until,omitempty"`
}

type JsonField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Since       string `json:"since,omitempty"`
	Until       string `json:"until,omitempty"`
	Unique      bool   `json:"unique,omitempty"`
	Ref         string `json:"ref,omitempty"`
	RefField    string `json:"ref-field,omitempty"`
	Path        string `json:"path,omitempty"`
}

func (dp *DataParser) typeFromJSON(jsonData []byte) (DataFormat, error) {
	jfmt := JsonFormat{}

	jdec := json.NewDecoder(bytes.NewReader(jsonData))
	jdec.DisallowUnknownFields()
	err := jdec.Decode(&jfmt)
	if err != nil {
		return DataFormat{}, err
	}

	if jfmt.File == "" {
		return DataFormat{}, fmt.Errorf("Missing \"file\" property in data format JSON")
	}

	fields := make([]DataField, 0, len(jfmt.Fields))

	for _, jf := range jfmt.Fields {
		if jf.Name == "" {
			return DataFormat{}, fmt.Errorf("Missing \"name\" property in field")
		}
		if jf.Type == "" {
			return DataFormat{}, fmt.Errorf("Missing \"type\" property in field \"%s\"", jf.Name)
		}

		df := DataField{
			Name: jf.Name,
			Type: FieldType(jf.Type),
		}
		if !df.Type.Valid() {
			return DataFormat{}, fmt.Errorf("Invalid type in field \"%s\"", jf.Name)
		}

		if jf.Since != "" && version.CompareSimple(jf.Since, dp.version) > 0 {
			continue
		}
		if jf.Until != "" && version.CompareSimple(jf.Until, dp.version) <= 0 {
			continue
		}

		fields = append(fields, df)
	}

	return DataFormat{
		Name:   jfmt.File,
		Fields: fields[:],
	}, nil
}
