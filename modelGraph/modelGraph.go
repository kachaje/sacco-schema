package modelgraph

import "fmt"

func Main() {

}

func checkParent(relationMaps map[string]any, model string) map[string]any {
	if _, ok := relationMaps[model]; !ok {
		relationMaps[model] = map[string]any{
			"children": []string{},
		}
	}

	return relationMaps
}

func CreateGraph(rawData map[string]any) (map[string]any, error) {
	relationMaps := map[string]any{}

	for model, value := range rawData {
		if data, ok := value.(map[string]any); ok {
			checkParent(relationMaps, model)

			parents := []string{}

			if val, ok := data["parents"].([]string); ok {
				parents = val
			} else if val, ok := data["parents"].([]any); ok {
				for _, key := range val {
					parents = append(parents, fmt.Sprintf("%v", key))
				}
			}

			for _, key := range parents {
				checkParent(relationMaps, key)

				relationMaps[key].(map[string]any)["children"] = append(relationMaps[key].(map[string]any)["children"].([]string), model)
			}

		}
	}

	return relationMaps, nil
}
