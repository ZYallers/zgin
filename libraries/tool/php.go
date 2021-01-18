package tool

import (
	"github.com/syyongx/php2go"
	"github.com/techoner/gophp/serialize"
	"strings"
)

// PhpUnserialize
func PhpUnserialize(str string) map[string]interface{} {
	vars := make(map[string]interface{}, 10)
	offset := 0
	strlen := php2go.Strlen(str)
	for offset < strlen {
		if index := strings.Index(php2go.Substr(str, uint(offset), -1), "|"); index < 0 {
			break
		}

		pos := php2go.Strpos(str, "|", offset)
		num := pos - offset

		varname := php2go.Substr(str, uint(offset), num)
		offset += num + 1
		data, _ := serialize.UnMarshal([]byte(php2go.Substr(str, uint(offset), -1)))
		vars[varname] = data

		jsonbyte, _ := serialize.Marshal(data)
		offset += php2go.Strlen(string(jsonbyte))
	}
	return vars
}

// PhpSerialize
func PhpSerialize(vars map[string]interface{}) (str string) {
	for k, v := range vars {
		shal, _ := serialize.Marshal(v)
		str += k + "|" + string(shal)
	}
	return
}
