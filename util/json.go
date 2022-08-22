package util

import (
	"encoding/json"
	"io"
	"os"
)

func ReadJson(r io.Reader, target any) error {
	jdec := json.NewDecoder(r)
	jdec.DisallowUnknownFields()
	return jdec.Decode(target)
}

func WriteJson(w io.Writer, data any, pretty bool) error {
	jenc := json.NewEncoder(w)
	jenc.SetEscapeHTML(false)
	if pretty {
		jenc.SetIndent("", "  ")
	}
	return jenc.Encode(data)
}

func ReadJsonFromFile(filename string, target any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return ReadJson(f, target)
}

func WriteJsonToFile(filename string, data any, pretty bool) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteJson(f, data, pretty)
}
