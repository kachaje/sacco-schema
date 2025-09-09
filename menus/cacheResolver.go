package menus

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
)

func ResolveCacheData(data, cacheMap map[string]any) map[string]any {
	var result = map[string]any{}
	var groupRoot string

	if cacheMap != nil && cacheMap["cacheQueries"] != nil && data != nil {
		if cacheQueries, ok := cacheMap["cacheQueries"].(map[string]any); ok {
			for _, value := range cacheQueries {
				if groupRoot == "" {
					reRoot := regexp.MustCompile(`^([^0]+)`)
					if val, ok := value.(string); ok {
						if reRoot.MatchString(val) {
							groupRoot = reRoot.FindString(val)
							break
						}
					}
				}
			}
		}
	}

	incomingData := map[string]any{}

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
