package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"
)

var JsonSchema = cli.Command{
	Name:      "jsonschema",
	Usage:     "Detect schema of a JSON file (such as the stash tab river)",
	UsageText: "pogo jsonschema <file.json>",

	Flags: []cli.Flag{},

	Action: do_jsonschema,
}

func do_jsonschema(c *cli.Context) error {
	var err error

	r, err := os.Open(c.Args().Get(0))
	if err != nil {
		return err
	}

	d := json.NewDecoder(r)

	var input interface{}
	err = d.Decode(&input)
	if err != nil {
		return err
	}

	dumper := &schemaDumper{
		seen: make(map[string]bool),
	}

	dumper.dump(input, "")

	return nil
}

type schemaDumper struct {
	seen map[string]bool
}

func (d *schemaDumper) dump(val interface{}, path string) {
	switch val := val.(type) {
	case nil:
		d.showPathType(path + " NULL")
	case bool:
		d.showPathType(path + " bool")
	case string:
		d.showPathType(path + " string")
	case float64:
		if float64(int(val)) != val {
			d.showPathType(path + " float")
		} else {
			d.showPathType(path + " int")
		}
	case map[string]interface{}:
		for key, kval := range val {
			d.dump(kval, path+"."+key)
		}
	case []interface{}:
		for i := range val {
			d.dump(val[i], path+"[]")
		}
	default:
		panic(fmt.Sprintf("%s has type %T ???\n", path, val))
	}
}

func (d *schemaDumper) showPathType(t string) {
	_, seen := d.seen[t]
	if !seen {
		fmt.Println(t)
		d.seen[t] = true
	}
}
