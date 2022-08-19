package cmd

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/oriath-net/pogo/dat"
	"github.com/oriath-net/pogo/util"

	cli "github.com/urfave/cli/v2"
)

type schemaTopLevel struct {
	Version      int                 `json:"version"`
	CreatedAt    int                 `json:"createdAt"`
	Tables       []schemaTable       `json:"tables"`
	Enumerations []schemaEnumeration `json:"enumerations"`
}

type schemaTable struct {
	Name    string         `json:"name"`
	Columns []schemaColumn `json:"columns"`
}

type schemaColumn struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Array       bool   `json:"array"`
	Type        string `json:"type"`
	Unique      bool   `json:"unique"`
	Localized   bool   `json:"localized"` // TODO
	References  struct {
		Table  string `json:"table"`
		Column string `json:"column"`
	} `json:"references"`
	Since string   `json:"since"` // doesn't exist yet, but it should
	Until string   `json:"until"`
	File  string   `json:"file"`  // TODO
	Files []string `json:"files"` // TODO
}

type schemaEnumeration struct {
	Name        string   `json:"name"`
	Indexing    int      `json:"indexing"`
	Enumerators []string `json:"enumerators"`
}

type schemaMetaData struct {
	Description string `json:"description"`
	Version     int    `json:"version"`
	CreatedAt   int    `json:"createdAt"`
}

var Schema2json = cli.Command{
	Name:      "schema2json",
	Usage:     "Convert schema.min.json from poe-tool-dev to a directory of JSON format definitions",
	UsageText: "pogo schema2json <schema.min.json> <output directory>",
	Action:    do_schema2json,
}

func do_schema2json(c *cli.Context) error {
	if c.NArg() < 2 {
		return errNotEnoughArguments
	}
	if c.NArg() > 2 {
		return errTooManyArguments
	}

	schema := schemaTopLevel{}
	err := util.ReadJsonFromFile(c.Args().Get(0), &schema)
	if err != nil {
		return err
	}

	outdir := c.Args().Get(1)
	err = os.Mkdir(outdir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	util.WriteJsonToFile(filepath.Join(outdir, "_META.json"), schemaMetaData{
		"Created by pogo schema2json",
		schema.Version,
		schema.CreatedAt,
	}, true)

	identifierRegexp := regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*$`)

	for _, tbl := range schema.Tables {
		jfmt := dat.JsonFormat{
			File:   tbl.Name,
			Fields: []dat.JsonField{},
		}

		for i, sf := range tbl.Columns {
			jfield := dat.JsonField{
				Name:        sf.Name,
				Description: sf.Description,
				Unique:      sf.Unique,
				Ref:         sf.References.Table,
				RefField:    sf.References.Column,
				Since:       sf.Since,
				Until:       sf.Until,
			}

			if sf.Name != "" && !identifierRegexp.MatchString(sf.Name) {
				log.Printf("Invalid name '%s' in %s - using placeholder name", sf.Name, tbl.Name)
				sf.Name = ""
			}

			if sf.Name == "" {
				if sf.Type == "bool" {
					jfield.Name = "Flag" + strconv.Itoa(i)
				} else if sf.Type == "foreignrow" || sf.Type == "row" {
					jfield.Name = "Key" + strconv.Itoa(i)
				} else {
					jfield.Name = "Unknown" + strconv.Itoa(i)
				}
			}

			switch sf.Type {
			case "foreignrow":
				jfield.Type = "longid"
			case "row":
				jfield.Type = "shortid"
			case "array":
				jfield.Type = "void"
			default:
				jfield.Type = sf.Type
			}
			if sf.Array == true {
				jfield.Type = jfield.Type + "[]"
			}

			jfmt.Fields = append(jfmt.Fields, jfield)
		}
		util.WriteJsonToFile(filepath.Join(outdir, tbl.Name+".json"), jfmt, true)
	}

	return nil
}
