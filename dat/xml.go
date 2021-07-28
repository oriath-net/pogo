package dat

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"github.com/mcuadros/go-version"
)

type xmlField struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Since string `xml:"since,attr"`
	Until string `xml:"until,attr"`
}

type xmlFormat struct {
	File   string     `xml:"file,attr"`
	Fields []xmlField `xml:"field"`
}

func (dp *DataParser) typeFromXML(xmlData []byte) (DataFormat, error) {
	xfmt := xmlFormat{}
	err := xml.Unmarshal(xmlData, &xfmt)
	if err != nil {
		return DataFormat{}, err
	}

	if xfmt.File == "" {
		return DataFormat{}, fmt.Errorf("Missing file attribute in <format>")
	}

	fields := make([]DataField, 0, len(xfmt.Fields))

	for _, xf := range xfmt.Fields {
		if xf.Name == "" {
			return DataFormat{}, fmt.Errorf("Missing name attribute in <field>")
		}
		if xf.Type == "" {
			return DataFormat{}, fmt.Errorf("Missing type attribute in <field name=\"%s\">", xf.Name)
		}

		df := DataField{
			Name: xf.Name,
			Type: FieldType(xf.Type),
		}
		if !df.Type.Valid() {
			return DataFormat{}, fmt.Errorf("Invalid type in <field name=\"%s\">", xf.Name)
		}

		if xf.Since != "" && version.CompareSimple(xf.Since, dp.version) > 0 {
			continue
		}
		if xf.Until != "" && version.CompareSimple(xf.Until, dp.version) <= 0 {
			continue
		}

		fields = append(fields, df)
	}

	return DataFormat{
		Name:   xfmt.File,
		Fields: fields[:],
	}, nil
}
