package menus

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
)

func ResolveCacheData(data map[string]any, groupRoot string) map[string]any {
	var result = map[string]any{}
	var incomingData = map[string]any{}

	keys := []string{}

	if regexp.MustCompile(`\.`).MatchString(groupRoot) {
		groupRoot = regexp.MustCompile(`\.`).ReplaceAllLiteralString(groupRoot, `\.`)
	}

	if regexp.MustCompile(`0`).MatchString(groupRoot) {
		groupRoot = regexp.MustCompile(`0`).ReplaceAllLiteralString(groupRoot, `(\d+)`)
	}

	groupRoot = fmt.Sprintf(`^%s`, groupRoot)

	for key, value := range data {
		var re1 = regexp.MustCompile(fmt.Sprintf(`%s(\d+)\.([A-Za-z]+)$`, groupRoot))
		var re2 = regexp.MustCompile(fmt.Sprintf(`%s([A-Za-z]+)$`, groupRoot))

		if regexp.MustCompile(groupRoot).MatchString(key) &&
			(re1.MatchString(key) || regexp.MustCompile(`\\d+`).MatchString(key)) {
			var parts [][]string

			if re1.MatchString(key) {
				parts = re1.FindAllStringSubmatch(key, -1)
			} else if re2.MatchString(key) {
				parts = re2.FindAllStringSubmatch(key, -1)
			} else {
				continue
			}

			indexKey := parts[0][1]
			field := parts[0][2]

			if incomingData[indexKey] == nil {
				incomingData[indexKey] = map[string]any{}

				if !slices.Contains(keys, indexKey) {
					keys = append(keys, indexKey)
				}
			}

			incomingData[indexKey].(map[string]any)[field] = value
		} else if regexp.MustCompile(fmt.Sprintf(`^%s[A-Za-z]+$`, groupRoot)).MatchString(key) {
			field := regexp.MustCompile(groupRoot).ReplaceAllLiteralString(key, "")

			result[field] = value
		}
	}

	sort.Strings(keys)

	for i, key := range keys {
		value := incomingData[key]

		if val, ok := value.(map[string]any); ok {
			for k, v := range val {
				newKey := fmt.Sprintf("%v%v", k, i+1)
				result[newKey] = v
			}
		}
	}

	return result
}
