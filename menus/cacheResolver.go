package menus

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
)

func ResolveCacheData(data map[string]any, groupRoot string) map[string]any {
	var result = map[string]any{}
	var incomingData = map[string]any{}

	keys := []string{}

	for key, value := range data {
		re := regexp.MustCompile(fmt.Sprintf(`%s(\d+)\.(.+)`, groupRoot))
		if strings.HasPrefix(key, groupRoot) && re.MatchString(key) {
			parts := re.FindAllStringSubmatch(key, -1)

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
			if regexp.MustCompile(`[A-Za-z]\.`).MatchString(groupRoot) {
				groupRoot = regexp.MustCompile(`\.`).ReplaceAllLiteralString(groupRoot, `\.`)
			}

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
