package prettyjson

import "encoding/json"

func Stringify(v any) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
