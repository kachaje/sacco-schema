package database

import (
	"fmt"
	"os"
	"regexp"
	"sacco/utils"
	"strings"
)

func Main(model, destinationFile string, sourceData map[string]any) (*string, map[string][]string, map[string]bool, []string, error) {
	data := map[string]any{
		"model": model,
		"formSummary": map[string]any{
			"type": "quitScreen",
		},
	}

	parentModels := []string{}
	floatFields := map[string]bool{}
	relationships := map[string][]string{}
	j := 0
	lastTag := ""

	if rawData, ok := sourceData[model].(map[string]any); ok {
		count := 1
		arrayModel := false

		if rawData["rootQuery"] != nil {
			data["rootQuery"] = rawData["rootQuery"]

			if strings.HasSuffix(fmt.Sprintf("%v", rawData["rootQuery"]), ".0") {
				arrayModel = true
			}
		}

		if rawData["settings"] != nil {
			if val, ok := rawData["settings"].(map[string]any); ok && val["hasLoops"] != nil && val["totalLoops"] != nil {
				if totalLoops, ok := val["totalLoops"].(int); ok {
					count = totalLoops
				} else if totalLoops, ok := val["totalLoops"].(int64); ok {
					count = int(totalLoops)
				} else if totalLoops, ok := val["totalLoops"].(float64); ok {
					count = int(totalLoops)
				}
			}
		}

		if rawData["hasMany"] != nil {
			if val, ok := rawData["hasMany"].([]any); ok {
				values := []string{}

				for _, v := range val {
					if vs, ok := v.(string); ok {
						values = append(values, vs)
					}
				}

				relationships["hasMany"] = values
			}
		}
		if rawData["hasOne"] != nil {
			if val, ok := rawData["hasOne"].([]any); ok {
				values := []string{}

				for _, v := range val {
					if vs, ok := v.(string); ok {
						values = append(values, vs)
					}
				}

				relationships["hasOne"] = values
			}
		}
		if rawData["belongsTo"] != nil {
			if val, ok := rawData["belongsTo"].([]any); ok {
				values := []string{}

				for _, v := range val {
					if vs, ok := v.(string); ok {
						values = append(values, vs)
					}
				}

				parentModels = values
			}
		}

		cacheQueries := map[string]string{}

		for index := range count {
			suffix := ""
			if count > 1 || arrayModel {
				suffix = fmt.Sprint(index + 1)
			}

			fields := []any{}

			if rawData["arrayFields"] != nil {
				if val, ok := rawData["arrayFields"].([]any); ok {
					fields = val
				} else if val, ok := rawData["fields"].([]any); ok {
					fields = val
				}
			} else if val, ok := rawData["fields"].([]any); ok {
				fields = val
			}

			for _, row := range fields {
				if val, ok := row.(map[string]any); ok {
					for key, rawValue := range val {
						if value, ok := rawValue.(map[string]any); ok {
							if regexp.MustCompile(`\d+$`).MatchString(key) {
								suffix = ""
							}

							tag := fmt.Sprintf("enter%s%v", utils.CapitalizeFirstLetter(key), suffix)

							if len(parentModels) > 0 {
								parentIds := []string{}

								for _, name := range parentModels {
									parentIds = append(parentIds, fmt.Sprintf("%sId", name))
								}

								data["parentIds"] = parentIds
							}

							if data["initialScreen"] == nil && value["hidden"] == nil {
								data["initialScreen"] = tag
							}

							inputIdentifier := fmt.Sprintf("%s%v", key, suffix)

							data[tag] = map[string]any{
								"inputIdentifier": inputIdentifier,
							}

							if rawData["rootQuery"] != nil {
								rootQuery := fmt.Sprintf("%v", rawData["rootQuery"])

								re := regexp.MustCompile(`(#\d+#)`)

								matches := re.FindAllStringSubmatch(rootQuery, -1)

								if len(matches) > 0 {
									size := len(matches)

									for i, val := range matches {
										if i == size-1 {
											rootQuery = regexp.MustCompile(val[0]).ReplaceAllLiteralString(rootQuery, fmt.Sprint(index))
										} else {
											rootQuery = regexp.MustCompile(val[0]).ReplaceAllLiteralString(rootQuery, "0")
										}
									}
								}

								localSuffix := key

								reSuffix := regexp.MustCompile(`\d+$`)
								if reSuffix.MatchString(key) {
									localSuffix = reSuffix.ReplaceAllLiteralString(key, "")
								}

								cacheQuery := fmt.Sprintf("%v.%s", rootQuery, localSuffix)

								cacheQueries[inputIdentifier] = cacheQuery
							}

							if value["readOnly"] != nil ||
								value["formula"] != nil {
								data[tag].(map[string]any)["readOnly"] = true
								data[tag].(map[string]any)["type"] = "inputScreen"

								text := utils.IdentifierToLabel(key)

								data[tag].(map[string]any)["text"] = map[string]any{
									"en": text,
								}

								if value["formula"] != nil {
									data[tag].(map[string]any)["formula"] = value["formula"].(string)
								}

								if value["numericField"] != nil {
									floatFields[key] = true
								}
							} else if value["hidden"] == nil {
								j++

								text := utils.IdentifierToLabel(key)

								data[tag].(map[string]any)["text"] = map[string]any{
									"en": text,
								}

								data[tag].(map[string]any)["order"] = j
								data[tag].(map[string]any)["type"] = "inputScreen"
								data[tag].(map[string]any)["nextScreen"] = "formSummary"

								if value["optional"] != nil {
									data[tag].(map[string]any)["optional"] = true
								}

								if value["numericField"] != nil {
									data[tag].(map[string]any)["validationRule"] = "^\\d+\\.*\\d*$"

									floatFields[key] = true
								}

								if value["validationRule"] != nil {
									data[tag].(map[string]any)["validationRule"] = value["validationRule"].(string)
								}

								if value["terminateBlockOnEmpty"] != nil {
									data[tag].(map[string]any)["terminateBlockOnEmpty"] = true
								}

								if value["adminOnly"] != nil {
									data[tag].(map[string]any)["adminOnly"] = true
								}

								if value["formula"] != nil {
									data[tag].(map[string]any)["formula"] = value["formula"].(string)
								}

								if value["scheduleFormula"] != nil {
									data[tag].(map[string]any)["scheduleFormula"] = value["scheduleFormula"].(string)
									data[tag].(map[string]any)["optional"] = true
								}

								if value["options"] != nil {
									if opts, ok := value["options"].([]any); ok {
										options := []any{}

										for i, opt := range opts {
											option := map[string]any{
												"position": i + 1,
												"label": map[string]any{
													"en": opt,
												},
											}

											options = append(options, option)
										}

										data[tag].(map[string]any)["options"] = options
									}
								}

								if lastTag != "" {
									data[lastTag].(map[string]any)["nextScreen"] = tag
								}

								lastTag = tag
							} else {
								data[tag].(map[string]any)["hidden"] = true
								data[tag].(map[string]any)["type"] = "hiddenField"
							}
						}
					}
				}
			}
		}

		data["cacheQueries"] = cacheQueries
	}

	yamlString, err := utils.DumpYaml(data)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = os.WriteFile(destinationFile, []byte(*yamlString), 0644)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return yamlString, relationships, floatFields, parentModels, nil
}
