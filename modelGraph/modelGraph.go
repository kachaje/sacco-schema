package modelgraph

import (
	"fmt"
	"strings"
)

func Main() {

}

func checkParent(relationMaps map[string]any, model string) map[string]any {
	if _, ok := relationMaps[model]; !ok {
		relationMaps[model] = map[string]any{
			"singleChildren": []string{},
			"arrayChildren":  []string{},
		}
	}

	return relationMaps
}

func CreateGraph(rawData map[string]any) (map[string]any, error) {
	relationMaps := map[string]any{}

	for model, value := range rawData {
		if data, ok := value.(map[string]any); ok {
			parents := []string{}

			if val, ok := data["parents"].([]string); ok {
				parents = val
			} else if val, ok := data["parents"].([]any); ok {
				for _, key := range val {
					parents = append(parents, fmt.Sprintf("%v", key))
				}
			}

			for _, key := range parents {
				if strings.HasSuffix(model, "IdsCache") {
					checkParent(relationMaps, model)

					relationMaps[model].(map[string]any)["singleChildren"] = append(relationMaps[model].(map[string]any)["singleChildren"].([]string), key)
				} else {
					checkParent(relationMaps, key)

					if data["many"] != nil {
						relationMaps[key].(map[string]any)["arrayChildren"] = append(relationMaps[key].(map[string]any)["arrayChildren"].([]string), model)
					} else {
						relationMaps[key].(map[string]any)["singleChildren"] = append(relationMaps[key].(map[string]any)["singleChildren"].([]string), model)
					}
				}
			}

		}
	}

	for model, value := range relationMaps {
		if data, ok := value.(map[string]any); ok {
			for key, val := range data {
				if vl, ok := val.([]any); ok && len(vl) <= 0 {
					delete(relationMaps[model].(map[string]any), key)
				} else if vl, ok := val.([]string); ok && len(vl) <= 0 {
					delete(relationMaps[model].(map[string]any), key)
				}
			}
		}
	}

	return relationMaps, nil
}
