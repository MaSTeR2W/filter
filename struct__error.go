package filter

import (
	"encoding/json"
	"strconv"
	"strings"
)

var OmitVal = &[]int8{}

type FilterErr struct {
	Key     string
	Value   any
	Path    []any
	Message string
}

func (f *FilterErr) Error() string {
	return f.Message
}

func (f *FilterErr) MarshalJSON() ([]byte, error) {
	var js = ""
	if f.Key != "" {
		js += `,"key":"` + strings.ReplaceAll(f.Key, `"`, `\"`) + `"`
	}

	if f.Value != OmitVal {
		js += `,"value":`
		var err error
		var jsVal []byte

		if jsVal, err = json.Marshal(f.Value); err != nil {
			return nil, err
		}

		js += string(jsVal)
	}

	var els = ""
	for _, el := range f.Path {
		switch t := el.(type) {
		case int:
			els += "," + strconv.Itoa(t)
		case string:
			els += `,"` + strings.ReplaceAll(t, `"`, `\"`) + `"`
		}
	}

	if len(els) > 0 {
		els = els[1:]
	}

	js += `,"path":[` + els + `]`

	if f.Message != "" {
		js += `,"message":"` + strings.ReplaceAll(f.Message, `"`, `\"`) + `"`
	}

	if len(js) == 0 {
		return []byte("{}"), nil
	}

	return []byte("{" + js[1:] + "}"), nil

}
