package shared

import "encoding/json"

func ToJSON(v any) string {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(json)
}
