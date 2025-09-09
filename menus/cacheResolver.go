package menus

import (
	"fmt"
	"regexp"
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

	for key, value := range data {
		re := regexp.MustCompile(fmt.Sprintf(`%s(\d+)\.(.+)`, groupRoot))

		if strings.HasPrefix(key, groupRoot) && re.MatchString(key) {
			parts := re.FindAllStringSubmatch(key, -1)

			indexKey := parts[0][1]
			field := parts[0][2]

			if incomingData[indexKey] == nil {
				incomingData[indexKey] = map[string]any{}
			}

			incomingData[indexKey].(map[string]any)[field] = value
		}
	}

	i := 0
	for _, value := range incomingData {
		i++

		if val, ok := value.(map[string]any); ok {
			for k, v := range val {
				newKey := fmt.Sprintf("%v%v", k, i)
				result[newKey] = v
			}
		}
	}

	return result
}
