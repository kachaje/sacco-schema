package database

import (
	"fmt"
	"log"
	"maps"
	"strconv"

	"github.com/kachaje/utils/utils"
)

func (d *Database) LoadSingleChildren(parentKey, model string, parentId int64) (map[string]any, error) {
	data := map[string]any{}

	capModel := utils.CapitalizeFirstLetter(model)

	if chidren, ok := SingleChildren[fmt.Sprintf("%sSingleChildren", capModel)]; ok {
		for _, childModel := range chidren {
			parentKey := fmt.Sprintf("%sId", model)

			results, err := d.GenericModels[childModel].FilterBy(fmt.Sprintf(`WHERE %s = %v AND active = 1 ORDER by updatedAt DESC LIMIT 1`, parentKey, parentId))
			if err != nil {
				log.Println("LoadSingleChildren 1:", childModel, err)
				continue
			}

			if len(results) > 0 {
				row := results[0]

				for _, field := range d.SkipFields {
					delete(row, field)
				}

				if id, err := strconv.ParseInt(fmt.Sprintf("%v", row["id"]), 10, 64); err == nil {
					result, err := d.LoadModelChildren(childModel, id)
					if err != nil {
						log.Println("LoadSingleChildren 2:", childModel, err)
					} else {
						maps.Copy(row, result)
					}
				}

				data[childModel] = row
			}
		}
	}

	return data, nil
}

func (d *Database) LoadArrayChildren(parentKey, model string, parentId int64) (map[string]any, error) {
	data := map[string]any{}

	capModel := utils.CapitalizeFirstLetter(model)

	if arrayChidren, ok := ArrayChildren[fmt.Sprintf("%sArrayChildren", capModel)]; ok {
		for _, childModel := range arrayChidren {
			parentKey := fmt.Sprintf("%sId", model)

			results, err := d.GenericModels[childModel].FilterBy(fmt.Sprintf(`WHERE %s = %v AND active = 1`, parentKey, parentId))
			if err != nil {
				log.Println("LoadArrayChildren 1:", childModel, err)
				continue
			}

			if len(results) > 0 {
				rows := map[string]any{}

				for _, row := range results {

					for _, field := range d.SkipFields {
						delete(row, field)
					}

					if id, err := strconv.ParseInt(fmt.Sprintf("%v", row["id"]), 10, 64); err == nil {
						result, err := d.LoadModelChildren(childModel, id)
						if err != nil {
							log.Println("LoadArrayChildren 2:", childModel, err)
						} else {
							maps.Copy(row, result)
						}

						rows[fmt.Sprintf("%v", id)] = row
					}
				}

				data[childModel] = rows
			}
		}
	}

	return data, nil
}

func (d *Database) LoadModelChildren(model string, id int64) (map[string]any, error) {
	data, err := d.GenericModels[model].FetchById(id)
	if err != nil {
		return nil, err
	}

	if len(data) <= 0 {
		return nil, fmt.Errorf("no match found")
	}

	for _, field := range d.SkipFields {
		delete(data, field)
	}

	parentKey := fmt.Sprintf("%sId", model)

	arrayChildren, err := d.LoadArrayChildren(parentKey, model, id)
	if err != nil {
		return nil, err
	}

	maps.Copy(data, arrayChildren)

	singleChidren, err := d.LoadSingleChildren(parentKey, model, id)
	if err != nil {
		return nil, err
	}

	maps.Copy(data, singleChidren)

	if data["memberDependant"] != nil {
		if val, ok := data["memberDependant"].(map[string]any); ok {
			for _, value := range val {
				if child, ok := value.(map[string]any); ok {
					if child["isNominee"] != nil && fmt.Sprintf("%v", child["isNominee"]) == "Yes" {
						data["memberNominee"] = child
						break
					}
				}
			}
		}
	}

	return data, nil
}

func (d *Database) FullMemberRecord(phoneNumber string) (map[string]any, error) {
	var data = map[string]any{}

	results, err := d.GenericModels["member"].FilterBy(fmt.Sprintf(`WHERE phoneNumber = "%s" AND active = 1`, phoneNumber))
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		if id, err := strconv.ParseInt(fmt.Sprintf("%v", results[0]["id"]), 10, 64); err == nil {
			result, err := d.LoadModelChildren("member", id)
			if err != nil {
				return nil, err
			}

			data = result
		}
	} else {
		return nil, fmt.Errorf("no match found")
	}

	return map[string]any{
		"member": data,
	}, nil
}
