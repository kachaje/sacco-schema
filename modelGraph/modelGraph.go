package modelgraph

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sacco/utils"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func Main(folder *string) error {
	workingFolder := filepath.Join(".", "database", "schema", "models")

	if folder != nil {
		workingFolder = *folder
	}

	if _, err := os.Stat(workingFolder); os.IsNotExist(err) {
		return fmt.Errorf("folder %s not found", workingFolder)
	}

	modelsData := map[string]any{}

	content, err := os.ReadFile(filepath.Join(workingFolder, "modelsData.json"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &modelsData)
	if err != nil {
		return err
	}

	graphsData, err := CreateGraph(modelsData)
	if err != nil {
		return err
	}

	payload, err := json.MarshalIndent(graphsData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(workingFolder, "graph.json"), payload, 0644)
	if err != nil {
		return err
	}

	modelsWorkflowData, err := CreateWorkflowGraph(modelsData, graphsData)
	if err != nil {
		return err
	}

	payload, err = json.MarshalIndent(modelsWorkflowData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(workingFolder, "models.json"), payload, 0644)
	if err != nil {
		return err
	}

	yamlData, err := utils.DumpYaml(modelsWorkflowData)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(workingFolder, "models.yml"), []byte(*yamlData), 0644)
	if err != nil {
		return err
	}

	return nil
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

							if vv["options"] != nil {
								if vo, ok := vv["options"].([]any); ok && len(vo) > 0 {
									values = append(values, fmt.Sprintf(`"%s"`, vo[0]))
								} else if vo, ok := vv["options"].([]string); ok && len(vo) > 0 {
									values = append(values, fmt.Sprintf(`"%s"`, vo[0]))
								}
							} else if vv["type"] != nil && fmt.Sprintf("%v", vv["type"]) == "text" {
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

func CreateWorkflowGraph(modelsData, graphData map[string]any) (map[string]any, error) {
	if modelsData == nil || graphData == nil {
		return nil, fmt.Errorf("all inputs required")
	}

	result := map[string]any{}

	for model, value := range modelsData {
		result[model] = map[string]any{
			"rootQuery": model,
			"fields":    []map[string]any{},
		}

		if graphData[model] != nil {
			if val, ok := graphData[model].(map[string]any); ok {
				if val["singleChildren"] != nil {
					if v, ok := val["singleChildren"].([]any); ok {
						result[model].(map[string]any)["hasOne"] = v
					} else if v, ok := val["singleChildren"].([]string); ok {
						result[model].(map[string]any)["hasOne"] = v
					}
				}

				if val["arrayChildren"] != nil {
					if v, ok := val["arrayChildren"].([]any); ok {
						result[model].(map[string]any)["hasMany"] = v
					} else if v, ok := val["arrayChildren"].([]string); ok {
						result[model].(map[string]any)["hasMany"] = v
					}
				}
			}
		}

		if val, ok := value.(map[string]any); ok {
			if val["parents"] != nil {
				if v, ok := val["parents"].([]any); ok && len(v) > 0 {
					result[model].(map[string]any)["belongsTo"] = v

					result[model].(map[string]any)["rootQuery"] = fmt.Sprintf("%v.%v", v[0], model)
				}
			}
			if val["fields"] != nil {
				if vv, ok := val["fields"].(map[string]any); ok {
					keysOrder := map[int]string{}

					for k, v := range vv {
						if vo, ok := v.(map[string]any); ok {
							if vo["order"] != nil {
								vi, err := strconv.Atoi(fmt.Sprintf("%v", vo["order"]))
								if err == nil {
									keysOrder[vi] = k
								}
							}
						}
					}

					for i := range len(keysOrder) {
						if k, ok := keysOrder[i]; ok {
							v := vv[k]

							row := map[string]any{}

							if regexp.MustCompile("id$").MatchString(strings.ToLower(k)) {
								row["hidden"] = true

								if vf, ok := v.(map[string]any); ok {
									if vf["order"] != nil {
										vi, err := strconv.Atoi(fmt.Sprintf("%v", vf["order"]))
										if err == nil {
											row["order"] = vi
										}
									}
								}
							} else {
								if vf, ok := v.(map[string]any); ok {
									for kf, vf := range vf {
										switch kf {
										case "hidden":
											row["optional"] = true
											row["hidden"] = true
										case "default", "optional":
											row["optional"] = true
										case "type", "order":
											if slices.Contains([]string{"int", "real"}, fmt.Sprintf("%v", vf)) {
												row["numericField"] = true
											} else {
												row[kf] = vf
											}
										case "options":
											if vo, ok := vf.([]any); ok {
												row[kf] = vo
											} else if vo, ok := vf.([]string); ok {
												row[kf] = vo
											}
										}
									}
								}
							}

							result[model].(map[string]any)["fields"] = append(result[model].(map[string]any)["fields"].([]map[string]any), map[string]any{
								k: row,
							})
						}
					}
				}
			}
		}
	}

	return result, nil
}
