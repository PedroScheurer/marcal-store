package modules

import "encoding/json"

func Encode(titles []string) string {
	if len(titles) == 0 {
		return "[]"
	}

	encoded, err := json.Marshal(titles)
	if err != nil {
		return "[]"
	}

	return string(encoded)
}

func Decode(raw string) []string {
	if raw == "" || raw == "[]" {
		return nil
	}

	var titles []string
	if err := json.Unmarshal([]byte(raw), &titles); err != nil {
		return nil
	}

	return titles
}

func ResolveCount(titles []string, modules int) int {
	if len(titles) > 0 {
		return len(titles)
	}
	if modules >= 0 {
		return modules
	}
	return 0
}
