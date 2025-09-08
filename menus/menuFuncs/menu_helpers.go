package menufuncs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func CheckPreferredLanguage(phoneNumber, preferencesFolder string) *string {
	settingsFile := filepath.Join(preferencesFolder, phoneNumber)

	_, err := os.Stat(settingsFile)
	if !os.IsNotExist(err) {
		content, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Println(err)
			return nil
		}

		data := map[string]any{}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Println(err)
			return nil
		}

		var preferredLanguage string

		if data["language"] != nil {
			val, ok := data["language"].(string)
			if ok {
				preferredLanguage = val
			}
		}

		return &preferredLanguage
	}

	return nil
}

func LoadGroupMembers(data map[string]any, target string) []map[string]any {
	rows := []map[string]any{}

	filteredRows := map[string]any{}
	keys := []string{}

	re := regexp.MustCompile(fmt.Sprintf(`^(.+%s\.\d+)\.([^\.]+)`, target))

	for key, value := range data {
		if re.MatchString(key) {
			newKey := re.FindAllStringSubmatch(key, -1)[0][1]
			nodeKey := re.FindAllStringSubmatch(key, -1)[0][2]

			if filteredRows[newKey] == nil {
				filteredRows[newKey] = map[string]any{}

				keys = append(keys, newKey)
			}

			filteredRows[newKey].(map[string]any)[nodeKey] = value
		}
	}

	sort.Strings(keys)

	for _, key := range keys {
		row := filteredRows[key]

		if val, ok := row.(map[string]any); ok {
			rows = append(rows, val)
		}
	}

	return rows
}

func LoadTemplateData(data map[string]any, template map[string]any) map[string]any {
	result := map[string]any{}

	for key, value := range template {
		result[key] = map[string]any{}

		var level string

		val, ok := value.(map[string]any)
		if ok {
			if val["level"] != nil {
				level = fmt.Sprintf("%v", val["level"])
			}

			fieldData := map[string]any{}

			v, ok := val["data"].(map[string]any)
			if ok {
				fieldData = v
			}

			switch level {
			case "memberDependant":
				groupData := LoadGroupMembers(data, level)

				var j float64

				keys := []string{}

				for key := range fieldData {
					keys = append(keys, key)
				}

				sort.Strings(keys)

				for i, row := range groupData {
					for _, field := range keys {
						kids := fieldData[field]
						if vf, ok := kids.(map[string]any); ok {
							localKey := fmt.Sprintf("%s%v", field, i+1)

							j++

							result[key].(map[string]any)[localKey] = map[string]any{
								"order": j,
								"label": fmt.Sprintf("%v", vf["label"]),
							}

							if value, ok := row[field]; ok {
								result[key].(map[string]any)[localKey].(map[string]any)["value"] = value
							}
						}
					}
				}
			default:
				for field, kids := range fieldData {
					if vf, ok := kids.(map[string]any); ok {
						if order, err := strconv.ParseFloat(fmt.Sprintf("%v", vf["order"]), 64); err == nil {
							result[key].(map[string]any)[field] = map[string]any{
								"order": order,
								"label": fmt.Sprintf("%v", vf["label"]),
							}

							if vf["cachQuery"] != nil {
								if query, ok := vf["cachQuery"].(string); ok {
									if value, ok := data[query]; ok {
										result[key].(map[string]any)[field].(map[string]any)["value"] = value
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}

func TabulateData(data map[string]any) []string {
	result := []string{}

	keys := []string{}

	for key := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		result = append(result, key)

		keysMap := map[float64]string{}
		floatKeys := []float64{}

		row, ok := data[key].(map[string]any)
		if ok {
			childData := map[string]map[string]any{}

			for k, v := range row {
				val, ok := v.(map[string]any)
				if ok && val["order"] != nil {
					childData[k] = val

					order, err := strconv.ParseFloat(fmt.Sprintf("%v", val["order"]), 64)
					if err == nil {
						keysMap[order] = k
						floatKeys = append(floatKeys, order)
					}
				}
			}

			sort.Float64s(floatKeys)

			if key == "E. BENEFICIARIES DETAILS" {
				row1 := "--- --------------------- --------- ------------"
				row2 := "No | Name of Beneficiary | Percent | Contact"

				result = append(result, row1)
				result = append(result, row2)
				result = append(result, row1)

				payload, _ := json.MarshalIndent(childData, "", "  ")

				fmt.Println(string(payload))

				for i := range 4 {
					index := i + 1

					nameLabel := fmt.Sprintf("name%d", index)
					percentageLabel := fmt.Sprintf("percentage%d", index)
					contactLabel := fmt.Sprintf("contact%d", index)

					if childData[nameLabel] == nil {
						break
					}

					var name string
					var percentage float64
					var contact string

					name = fmt.Sprintf("%v", childData[nameLabel]["value"])

					if childData[percentageLabel] != nil {
						v, err := strconv.ParseFloat(fmt.Sprintf("%v", childData[percentageLabel]["value"]), 64)
						if err == nil {
							percentage = v
						}
					}
					if childData[contactLabel] != nil {
						contact = fmt.Sprintf("%v", childData[contactLabel]["value"])
					}

					entry := fmt.Sprintf("%-3d| %-19s | %7.1f | %s", index, name, percentage, contact)

					result = append(result, entry)
				}
			} else {
				for _, order := range floatKeys {
					var label string
					var value string

					childKey := keysMap[order]

					if childData[childKey]["label"] != nil {
						label = fmt.Sprintf("%v:", childData[childKey]["label"])
					}
					if childData[childKey]["value"] != nil {
						value = fmt.Sprintf("%v", childData[childKey]["value"])
					}

					entry := fmt.Sprintf("   %-25s| %s", label, value)

					result = append(result, entry)
				}
			}

			result = append(result, "")
		}
	}

	if false {
		payload, _ := json.MarshalIndent(data, "", "  ")

		fmt.Println(string(payload))
	}

	return result
}

func LoadLoanApplicationForm(data map[string]any, template map[string]any) map[string]any {
	result := map[string]any{}

	keys := []string{}

	for key := range template {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	loadData := func(fieldData, parentMap map[string]any, key string) {
		for field, values := range fieldData {
			if value, ok := values.(map[string]any); ok {
				if order, err := strconv.ParseFloat(fmt.Sprintf("%v", value["order"]), 64); err == nil {
					parentMap[key].(map[string]any)[field] = map[string]any{
						"order": order,
						"label": fmt.Sprintf("%v", value["label"]),
					}

					if value["cachQuery"] != nil {
						if query, ok := value["cachQuery"].(string); ok {
							if val, ok := data[query]; ok {
								if value["formula"] != nil {
									if formula, ok := value["formula"].(string); ok {
										fmt.Println("######## TODO:", formula, val)
									}
								} else if len(strings.TrimSpace(fmt.Sprintf("%v", val))) > 0 {
									parentMap[key].(map[string]any)[field].(map[string]any)["value"] = val
								}
							}
						}
					}
				}
			}
		}
	}

	for _, key := range keys {
		result[key] = map[string]any{}

		if rawData, ok := template[key].(map[string]any); ok {
			if fieldData, ok := rawData["data"].(map[string]any); ok {
				loadData(fieldData, result, key)
			}

			if tables, ok := rawData["tables"]; ok {
				if tablesData, ok := tables.(map[string]any); ok {
					var tableLabel string

					if tablesData["label"] != nil {
						if val, ok := tablesData["label"].(string); ok {
							tableLabel = val
						}
					}

					if sectionsData, ok := tablesData["sections"].(map[string]any); ok {
						for section, sectionData := range sectionsData {
							if value, ok := sectionData.(map[string]any); ok {
								if fieldData, ok := value["data"].(map[string]any); ok {
									result[key].(map[string]any)[tableLabel] = map[string]any{
										section: map[string]any{},
									}

									loadData(fieldData, result[key].(map[string]any)[tableLabel].(map[string]any), section)
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}
