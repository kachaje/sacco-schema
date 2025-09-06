package modelgraph

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

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

			if len(parents) > 0 {
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
			} else {
				checkParent(relationMaps, model)
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
				} else {
					sort.Strings(relationMaps[model].(map[string]any)[key].([]string))
				}
			}
		}
	}

	return relationMaps, nil
}

func CreateModelQuery(model string, modelsData, seedData map[string]any) (*string, error) {
	var query string

	if modelsData != nil && modelsData[model] != nil {
		if val, ok := modelsData[model].(map[string]any); ok {
			if val["fields"] != nil {
				if data, ok := val["fields"].(map[string]any); ok {
					seed := map[string]any{}

					if seedData != nil {
						seed = seedData
					}

					orderMap := map[int]string{}

					for key, value := range data {
						if v, ok := value.(map[string]any); ok {
							if v["order"] != nil {
								index, err := strconv.Atoi(fmt.Sprintf("%v", v["order"]))
								if err == nil {
									orderMap[index] = key
								}
							}
						}
					}

					fields := []string{}
					values := []string{}

					for i := range len(orderMap) {
						key := orderMap[i]

						if vv, ok := data[key].(map[string]any); ok &&
							vv["autoIncrement"] == nil &&
							vv["optional"] == nil && vv["default"] == nil {
							var entry string

							if seed[key] != nil {
								entry = fmt.Sprintf("%v", seed[key])
							} else if vv["type"] != nil && fmt.Sprintf("%v", vv["type"]) == "text" {
								entry = key
							} else {
								entry = "10"
							}

							fields = append(fields, key)

							if vv["type"] != nil && fmt.Sprintf("%v", vv["type"]) == "text" {
								values = append(values, fmt.Sprintf(`"%s"`, entry))
							} else {
								values = append(values, entry)
							}
						}
					}

					query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s);`, model, strings.Join(fields, ", "), strings.Join(values, ", "))
				}
			}
		}
	}

	return &query, nil
}
