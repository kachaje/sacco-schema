package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kachaje/sacco-schema/database"
	"github.com/kachaje/utils/utils"
	"github.com/kachaje/workflow-parser/parser"
)

func SaveModelData(data any, model, phoneNumber *string,
	saveFunc func(map[string]any, string, int) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
	if rawData, ok := data.(map[string]any); ok {
		dataRows := utils.UnpackData(rawData)

		if refData != nil {
			var unpackedRefData []map[string]any

			// Check if refData is nested under model name (e.g., {"memberDependant": [...]})
			if model != nil && refData[*model] != nil {
				if arr, ok := refData[*model].([]map[string]any); ok {
					unpackedRefData = arr
				} else if arr, ok := refData[*model].([]any); ok {
					// Convert []any to []map[string]any
					unpackedRefData = make([]map[string]any, len(arr))
					for i, v := range arr {
						if m, ok := v.(map[string]any); ok {
							unpackedRefData[i] = m
						}
					}
				}
			} else {
				// Try to unpack as flat map
				unpackedRefData = utils.UnpackData(refData)
			}

			missingIds := utils.GetSkippedRefIds(dataRows, unpackedRefData)

			for _, row := range missingIds {
				if row["id"] != nil {
					updateData := map[string]any{
						"id":     row["id"],
						"active": 0,
					}
					dataRows = append(dataRows, updateData)
				}
			}
		}

		for _, modelData := range dataRows {
			if model != nil {
				for _, key := range database.FloatFields {
					if modelData[key] != nil {
						nv, ok := modelData[key].(string)
						if ok {
							real, err := strconv.ParseFloat(nv, 64)
							if err == nil {
								modelData[key] = real
							}
						}
					}
				}

				if saveFunc == nil {
					return fmt.Errorf("server.SaveModelData.%s:missing saveFunc", *model)
				}

				if sessions[*phoneNumber] != nil {
					if sessions[*phoneNumber].GlobalIds == nil {
						sessions[*phoneNumber].GlobalIds = map[string]any{}
					}
					if sessions[*phoneNumber].AddedModels == nil {
						sessions[*phoneNumber].AddedModels = map[string]bool{}
					}

					if database.ParentModels[*model] != nil {
						for _, value := range database.ParentModels[*model] {
							key := fmt.Sprintf("%sId", value)
							if sessions[*phoneNumber].GlobalIds[key] != nil {
								if val, ok := sessions[*phoneNumber].GlobalIds[key].(map[string]any); ok {
									vr, err := strconv.Atoi(fmt.Sprintf("%v", val["value"]))
									if err == nil {
										modelData[key] = vr
									}
								}
							}
						}
					}
				}

				if len(modelData) < 2 {
					continue
				}

				_, err := saveFunc(modelData, *model, 0)
				if err != nil {
					log.Println(err)
				}

				if *model == "member" && modelData["phoneNumber"] != nil && sessions[*phoneNumber] != nil {
					if val, ok := modelData["phoneNumber"].(string); ok {
						sessions[*phoneNumber].CurrentPhoneNumber = val
					}
				}
			}

			if sessions != nil && sessions[*phoneNumber] != nil {
				_, err := sessions[*phoneNumber].RefreshSession()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	return nil
}

func SaveData(
	data any, model, phoneNumber, preferenceFolder *string,
	saveFunc func(
		map[string]any,
		string,
		int,
	) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
	switch *model {
	case "language":
		val, ok := data.(map[string]any)
		if ok {
			if val["language"] != nil && phoneNumber != nil {
				language, ok := val["language"].(string)
				if ok {
					err := SavePreference(*phoneNumber, "language", language, *preferenceFolder)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

	default:
		return SaveModelData(data, model, phoneNumber, saveFunc, sessions, refData)
	}

	return nil
}

func SavePreference(phoneNumber, key, value, preferencesFolder string) error {
	settingsFile := filepath.Join(preferencesFolder, phoneNumber)

	data := map[string]any{}

	_, err := os.Stat(settingsFile)
	if !os.IsNotExist(err) {
		content, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Println(err)
			return err
		}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	data[key] = value

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}

	return os.WriteFile(settingsFile, payload, 0644)
}
