package menus

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	menufuncs "github.com/kachaje/sacco-schema/menus/menuFuncs"
	"github.com/kachaje/sacco-schema/parser"
	"github.com/kachaje/utils/utils"
)

//go:embed workflows/*
var RawWorkflows embed.FS

//go:embed menus/*
var menuFiles embed.FS

type Menus struct {
	ActiveMenus  map[string]any
	Titles       map[string]string
	Workflows    map[string]any
	RootQueries  map[string]string
	Functions    map[string]any
	FunctionsMap map[string]func(
		func(
			string, *parser.Session,
			string, string, string,
		) string,
		map[string]any,
		*parser.Session,
	) string
	TargetKeys    map[string][]string
	LabelWorkflow map[string]any
	CacheQueries  map[string]any

	mu sync.Mutex

	DevModeActive bool
	DemoMode      bool

	Cache      map[string]string
	LastPrompt string
}

func init() {
	err := fs.WalkDir(RawWorkflows, ".", func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(file, ".yml") {
			return nil
		}

		content, err := RawWorkflows.ReadFile(file)
		if err != nil {
			return err
		}

		data, err := utils.LoadYaml(string(content))
		if err != nil {
			log.Fatal(err)
		}

		model := strings.Split(filepath.Base(file), ".")[0]

		menufuncs.WorkflowsData[model] = data

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func NewMenus(devMode, demoMode *bool) *Menus {
	m := &Menus{
		ActiveMenus: map[string]any{},
		Titles:      map[string]string{},
		Workflows:   map[string]any{},
		RootQueries: map[string]string{},
		Functions:   map[string]any{},
		FunctionsMap: map[string]func(
			func(
				string, *parser.Session,
				string, string, string,
			) string,
			map[string]any,
			*parser.Session,
		) string{},
		TargetKeys:    map[string][]string{},
		CacheQueries:  map[string]any{},
		LabelWorkflow: map[string]any{},
		mu:            sync.Mutex{},

		Cache:      map[string]string{},
		LastPrompt: "username",
	}

	if devMode != nil {
		m.DevModeActive = *devMode
	}
	if demoMode != nil {
		m.DemoMode = *demoMode
	}

	m.FunctionsMap = menufuncs.FunctionsMap

	err := m.populateMenus()
	if err != nil {
		log.Panic(err)
	}

	return m
}

func (m *Menus) populateMenus() error {
	return fs.WalkDir(menuFiles, ".", func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(file, ".yml") {
			return nil
		}

		content, err := menuFiles.ReadFile(file)
		if err != nil {
			return err
		}

		data, err := utils.LoadYaml(string(content))
		if err != nil {
			log.Fatal(err)
		}

		re := regexp.MustCompile("Menu$")

		group := re.ReplaceAllLiteralString(strings.Split(filepath.Base(file), ".")[0], "")

		m.ActiveMenus[group] = map[string]any{}

		if val, ok := data["title"].(string); ok {
			m.Titles[group] = val
		}

		if data["allowedRoles"] != nil {
			allowedRoles := []string{}

			if val, ok := data["allowedRoles"].([]any); ok {
				for _, key := range val {
					allowedRoles = append(allowedRoles, fmt.Sprintf("%v", key))
				}
			} else if val, ok := data["allowedRoles"].([]string); ok {
				allowedRoles = append(allowedRoles, val...)
			}

			m.ActiveMenus[group].(map[string]any)["allowedRoles"] = allowedRoles
		}

		m.LabelWorkflow[group] = map[string]any{}

		if val, ok := data["fields"].(map[string]any); ok {
			keys := []string{}
			values := []string{}
			kv := map[string]any{}
			devMenus := map[string]any{}
			menuRoles := map[string][]string{}

			for key, row := range val {
				keys = append(keys, key)

				if val, ok := row.(map[string]any); ok {
					if val["id"] != nil && val["label"] != nil && val["label"].(map[string]any)["en"] != nil {
						if val["devOnly"] != nil && !m.DevModeActive {
							continue
						}

						id := fmt.Sprintf("%v", val["id"])
						label := fmt.Sprintf("%v", val["label"].(map[string]any)["en"])

						if val["allowedRoles"] != nil {
							allowedRoles := []string{}

							if vc, ok := val["allowedRoles"].([]any); ok {
								for _, key := range vc {
									allowedRoles = append(allowedRoles, fmt.Sprintf("%v", key))
								}
							} else if vc, ok := val["allowedRoles"].([]string); ok {
								allowedRoles = append(allowedRoles, vc...)
							}

							menuRoles[id] = allowedRoles
						}

						value := fmt.Sprintf("%v. %v\n", key, label)

						values = append(values, value)

						kv[key] = map[string]any{
							"menu":  id,
							"key":   key,
							"value": value,
						}

						if val["workflow"] != nil {
							if v, ok := val["workflow"].(string); ok {
								m.Workflows[id] = v

								if menufuncs.WorkflowsData[v]["rootQuery"] != nil {
									m.RootQueries[v] = fmt.Sprintf("%v", menufuncs.WorkflowsData[v]["rootQuery"])
								}

								parentIds := []string{}

								if menufuncs.WorkflowsData[v]["parentIds"] != nil {
									if val, ok := menufuncs.WorkflowsData[v]["parentIds"].([]string); ok {
										parentIds = val
									} else if val, ok := menufuncs.WorkflowsData[v]["parentIds"].([]any); ok {
										for _, vc := range val {
											parentIds = append(parentIds, fmt.Sprintf("%v", vc))
										}
									}
								}

								if menufuncs.WorkflowsData[v]["cacheQueries"] != nil {
									m.CacheQueries[v] = menufuncs.WorkflowsData[v]["cacheQueries"]
								}

								m.LabelWorkflow[group].(map[string]any)[value] = map[string]any{
									"model": v,
									"id":    id,
								}

								if len(parentIds) > 0 {
									m.LabelWorkflow[group].(map[string]any)[value].(map[string]any)["parentIds"] = parentIds
								}
							}
						}
						if val["function"] != nil {
							if v, ok := val["function"].(string); ok {
								m.Functions[id] = v
							}
						}
						if val["targetKeys"] != nil {
							if v, ok := val["targetKeys"].([]any); ok {
								m.TargetKeys[id] = []string{}

								for _, e := range v {
									if s, ok := e.(string); ok {
										m.TargetKeys[id] = append(m.TargetKeys[id], s)
									}
								}
							}
						}

						if val["allowedRoles"] != nil {
							allowedRoles := []string{}

							if vc, ok := val["allowedRoles"].([]any); ok {
								for _, key := range vc {
									allowedRoles = append(allowedRoles, fmt.Sprintf("%v", key))
								}
							} else if vc, ok := val["allowedRoles"].([]string); ok {
								allowedRoles = append(allowedRoles, vc...)
							}

							if m.ActiveMenus[group].(map[string]any)[id] == nil {
								m.ActiveMenus[group].(map[string]any)[id] = map[string]any{}
							}

							m.ActiveMenus[group].(map[string]any)[id].(map[string]any)["menuRoles"] = allowedRoles
						}
					}
				}
			}

			m.ActiveMenus[group].(map[string]any)["keys"] = keys
			m.ActiveMenus[group].(map[string]any)["kv"] = kv
			m.ActiveMenus[group].(map[string]any)["values"] = values
			m.ActiveMenus[group].(map[string]any)["devMenus"] = devMenus
			m.ActiveMenus[group].(map[string]any)["menuRoles"] = menuRoles
		}

		return nil
	})
}

func (m *Menus) LoadMenu(menuName string, session *parser.Session, phoneNumber, text, preferencesFolder string) string {
	var response string

	preferredLanguage := menufuncs.CheckPreferredLanguage(phoneNumber, preferencesFolder)

	if preferredLanguage != nil {
		session.PreferredLanguage = *preferredLanguage
	}

	if session == nil {
		return response
	}

	if session.SessionToken == nil && !m.DemoMode {
		switch session.CurrentMenu {
		case "signIn":
			return menufuncs.SignIn(
				m.LoadMenu,
				map[string]any{
					"phoneNumber":       phoneNumber,
					"session":           session,
					"preferredLanguage": preferredLanguage,
					"preferencesFolder": preferencesFolder,
					"text":              text,
				}, session)
		case "signUp":
			return menufuncs.SignUp(
				m.LoadMenu,
				map[string]any{
					"phoneNumber":       phoneNumber,
					"session":           session,
					"preferredLanguage": preferredLanguage,
					"preferencesFolder": preferencesFolder,
					"text":              text,
				}, session)
		default:
			return menufuncs.Landing(
				m.LoadMenu,
				map[string]any{
					"phoneNumber":       phoneNumber,
					"session":           session,
					"preferredLanguage": preferredLanguage,
					"preferencesFolder": preferencesFolder,
					"text":              text,
				}, session)
		}
	}

	keys := []string{}
	values := []string{}
	kv := map[string]string{}
	menuRoles := map[string][]string{}

	if val, ok := m.ActiveMenus[menuName].(map[string]any); ok {
		if val["menuRoles"] != nil {
			if v, ok := val["menuRoles"].(map[string][]string); ok {
				menuRoles = v
			}
		}
		if val["kv"] != nil {
			if v, ok := val["kv"].(map[string]any); ok {
				for key, value := range v {
					if vs, ok := value.(map[string]any); ok {
						menuId := fmt.Sprintf("%v", vs["menu"])
						menuKey := fmt.Sprintf("%v", vs["key"])
						menuValue := fmt.Sprintf("%v", vs["value"])

						if menuRoles[menuId] != nil {
							found := false

							for _, role := range menuRoles[menuId] {
								if session.SessionUserRole == nil {
									break
								} else if strings.EqualFold(*session.SessionUserRole, role) {
									found = true
									break
								}
							}

							if !found {
								continue
							}
						}

						keys = append(keys, menuKey)
						values = append(values, menuValue)
						kv[key] = menuId
					}
				}
			}
		}
	}

	filterValues := func(values []string) []string {
		newValues := []string{}

		keys = []string{}

		if m.LabelWorkflow[menuName] != nil && session != nil {
			for _, value := range values {
				if m.LabelWorkflow[menuName].(map[string]any)[value] != nil {
					model := fmt.Sprintf("%v", m.LabelWorkflow[menuName].(map[string]any)[value].(map[string]any)["model"])

					suffix := ""

					if session.AddedModels[model] {
						suffix = "(*)"
					}

					if m.LabelWorkflow[menuName].(map[string]any)[value].(map[string]any)["parentIds"] != nil {
						found := true

						checkIfExists := func(val []string) bool {
							found := true

							for _, key := range val {
								if session == nil || len(session.GlobalIds) <= 0 {
									found = false
								} else {
									_, found = session.GlobalIds[fmt.Sprintf("%v", key)]
								}

								if !found {
									break
								}
							}

							return found
						}

						if val, ok := m.LabelWorkflow[menuName].(map[string]any)[value].(map[string]any)["parentIds"].([]any); ok {
							vals := []string{}

							for _, key := range val {
								vals = append(vals, fmt.Sprintf("%v", key))
							}

							found = checkIfExists(vals)
						} else if val, ok := m.LabelWorkflow[menuName].(map[string]any)[value].(map[string]any)["parentIds"].([]string); ok {
							found = checkIfExists(val)
						}

						if !found {
							continue
						}
					}

					newValues = append(newValues, fmt.Sprintf("%s %s\n", strings.TrimSpace(value), suffix))
				} else {
					newValues = append(newValues, value)
				}

				keys = append(keys, strings.Split(value, ".")[0])
			}
		} else {
			newValues = values
		}

		return newValues
	}

	newValues := filterValues(values)

	if slices.Contains(keys, text) {
		target := text
		text = "000"

		if session != nil {
			session.CurrentMenu = kv[target]
			session.Cache = map[string]any{}

			model := fmt.Sprintf("%v", m.Workflows[kv[target]])

			if m.CacheQueries[model] != nil && m.RootQueries[model] != "" {
				groupRoot := fmt.Sprintf("%v.", m.RootQueries[model])

				if groupRoot != "" {
					session.Cache = ResolveCacheData(session.ActiveData, groupRoot)
				}
			}
		}

		return m.LoadMenu(kv[target], session, phoneNumber, text, preferencesFolder)
	} else if session != nil && m.Workflows[session.CurrentMenu] != nil {
		workingMenu := session.CurrentMenu
		model := fmt.Sprintf("%v", m.Workflows[workingMenu])

		if session.Cache != nil {
			if regexp.MustCompile(`^\d+$`).MatchString(phoneNumber) && session.WorkflowsMapping != nil &&
				session.WorkflowsMapping[model] != nil {
				if m.TargetKeys[workingMenu] != nil {
					for key, value := range session.Cache {
						if _, ok := session.WorkflowsMapping[model].Data[key]; !ok {
							session.WorkflowsMapping[model].Data[key] = value
						}
					}
				}

				if session.WorkflowsMapping[model].Data["phoneNumber"] == nil {
					session.WorkflowsMapping[model].Data["phoneNumber"] = session.CurrentPhoneNumber
				}
			}
		}

		if session.WorkflowsMapping != nil &&
			session.WorkflowsMapping[model] != nil {
			response = session.WorkflowsMapping[model].NavNext(text)

			if text == "00" {
				session.CurrentMenu = "main"
				text = "0"
				return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder)
			} else if strings.TrimSpace(response) == "" {
				if text == "0" {
					session.AddedModels[model] = true
				}

				parentMenu := "main"

				if regexp.MustCompile(`\.\d+$`).MatchString(session.CurrentMenu) {
					parentMenu = regexp.MustCompile(`\.\d+$`).ReplaceAllLiteralString(session.CurrentMenu, "")
				}

				session.CurrentMenu = parentMenu
				text = ""

				return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder)
			}
		} else {
			if text == "00" {
				session.CurrentMenu = "main"
				text = "0"
				return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder)
			}

			response = "NOT IMPLEMENTED YET\n\n" +
				"00. Main Menu\n"
		}

	} else {
		var menuRoot string

		if session != nil {
			menuRoot = session.CurrentMenu

			if regexp.MustCompile(`\.\d+$`).MatchString(session.CurrentMenu) {
				menuRoot = strings.Split(session.CurrentMenu, ".")[0]
			}

			if m.Functions[menuRoot] == nil && m.Functions[session.CurrentMenu] != nil {
				menuRoot = session.CurrentMenu
			}
		}

		if session != nil && m.Functions[menuRoot] != nil {
			if text == "00" {
				session.CurrentMenu = "main"
				text = "0"
				return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder)
			} else {
				if fnName, ok := m.Functions[menuRoot].(string); ok && m.FunctionsMap[fnName] != nil {
					response = m.FunctionsMap[fnName](
						m.LoadMenu,
						map[string]any{
							"phoneNumber":       phoneNumber,
							"preferredLanguage": preferredLanguage,
							"preferencesFolder": preferencesFolder,
							"text":              text,
						},
						session)
				} else {
					response = fmt.Sprintf("Function %s not found\n\n", m.Functions[menuRoot]) +
						"00. Main Menu\n"
				}
			}
		} else {
			utils.SortSlice(newValues)

			index := utils.Index(newValues, "99. Cancel\n")

			if index >= 0 {
				newValues = append(newValues[:index], newValues[index+1:]...)

				newValues = append(newValues, "\n99. Cancel")
			}

			index = utils.Index(newValues, "00. Main Menu\n")

			if index >= 0 {
				newValues = append(newValues[:index], newValues[index+1:]...)

				newValues = append(newValues, "\n00. Main Menu\n")
			}

			response = fmt.Sprintf("CON %s\n%s", m.Titles[menuName], strings.Join(newValues, ""))
		}
	}

	return response
}
