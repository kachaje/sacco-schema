package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func flattenRecursive(m map[string]any, prefix string, flat map[string]any, idMapOnly bool) {
	for key, value := range m {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]any:
			flattenRecursive(v, newKey, flat, idMapOnly)
		case []map[string]any, []any:
			arr := []map[string]any{}

			if vv, ok := v.([]any); ok {
				for _, vi := range vv {
					if vc, ok := vi.(map[string]any); ok {
						arr = append(arr, vc)
					}
				}
			} else if vv, ok := v.([]map[string]any); ok {
				arr = vv
			}

			for i, vc := range arr {
				newKey = fmt.Sprintf("%s.%v", prefix+"."+key, i)

				flattenRecursive(vc, newKey, flat, idMapOnly)
				if idMapOnly {
					break
				}
			}
		default:
			if idMapOnly && strings.HasSuffix(newKey, ".id") {
				re := regexp.MustCompile(`([A-Z-a-z]+)\.*0*\.id$`)

				if re.MatchString(newKey) {
					model := re.FindAllStringSubmatch(newKey, -1)[0][1]

					flat[model+"Id"] = map[string]any{
						"key":   newKey,
						"value": fmt.Sprintf("%v", v),
					}
				}
			} else if !idMapOnly {
				flat[newKey] = v
			}
		}
	}
}

func FlattenMap(m map[string]any, idMapOnly bool) map[string]any {
	flat := make(map[string]any)
	flattenRecursive(m, "", flat, idMapOnly)
	return flat
}

func SetNestedValue(m map[string]any, key string, value any) {
	parts := strings.Split(key, ".")
	currentMap := m

	for i, part := range parts {
		if i == len(parts)-1 {
			currentMap[part] = value
		} else {
			if _, ok := currentMap[part]; !ok {
				currentMap[part] = make(map[string]any)
			}

			currentMap = currentMap[part].(map[string]any)
		}
	}
}
