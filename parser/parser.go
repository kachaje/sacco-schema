package parser

import (
	"fmt"
	"log"
	"regexp"
	"sacco/utils"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

const (
	INITIAL_SCREEN = "initialScreen"
	INPUT_SCREEN   = "inputScreen"
	QUIT_SCREEN    = "quitScreen"
	LANG_EN        = "1"
	LANG_NY        = "2"
	LANG_EN_LABEL  = "en"
	LANG_NY_LABEL  = "ny"
)

type WorkFlow struct {
	Tree           map[string]any
	Data           map[string]any
	OptionalFields map[string]bool

	CurrentScreen         string
	NextScreen            string
	PreviousScreen        string
	CurrentLanguage       string
	CurrentPhoneNumber    string
	CurrentModel          string
	CurrentSessionId      string
	ScreenIdMap           map[string]string
	FormulaFields         map[string]string
	ScheduleFormulaFields map[string]string
	MatchModelFields      map[string]string
	ConditionFields       map[string]string
	ScreenOrder           map[int]string
	ReadOnlyFields        []string
	SubmitCallback        func(
		d any, m *string, p *string, f *string,
		addFn func(
			map[string]any,
			string,
			int,
		) (*int64, error),
		ss map[string]*Session,
		refData map[string]any,
	) error
	History          map[int]string
	HistoryIndex     int
	PreferenceFolder string
	AddFunc          func(
		map[string]any,
		string,
		int,
	) (*int64, error)

	Sessions map[string]*Session
}

func NewWorkflow(
	tree map[string]any,
	callbackFunc func(
		any, *string, *string, *string,
		func(
			map[string]any,
			string,
			int,
		) (*int64, error),
		map[string]*Session,
		map[string]any,
	) error,
	preferredLanguage, phoneNumber, sessionId, preferenceFolder *string, addFunc func(
		map[string]any,
		string,
		int,
	) (*int64, error), sessions map[string]*Session, refData map[string]any) *WorkFlow {

	w := &WorkFlow{
		Tree:                  tree,
		Data:                  map[string]any{},
		OptionalFields:        map[string]bool{},
		CurrentScreen:         INITIAL_SCREEN,
		CurrentLanguage:       LANG_EN,
		ScreenIdMap:           map[string]string{},
		ScreenOrder:           map[int]string{},
		SubmitCallback:        callbackFunc,
		History:               map[int]string{},
		HistoryIndex:          -1,
		FormulaFields:         map[string]string{},
		ScheduleFormulaFields: map[string]string{},
		MatchModelFields:      map[string]string{},
		ConditionFields:       map[string]string{},
		ReadOnlyFields:        []string{},
	}

	if sessions != nil {
		w.Sessions = sessions
	}
	if addFunc != nil {
		w.AddFunc = addFunc
	}
	if preferenceFolder != nil {
		w.PreferenceFolder = *preferenceFolder
	}
	if sessionId != nil {
		w.CurrentSessionId = *sessionId
	}
	if phoneNumber != nil {
		w.CurrentPhoneNumber = *phoneNumber
	}

	if preferredLanguage != nil {
		switch *preferredLanguage {
		case LANG_NY_LABEL:
			w.CurrentLanguage = LANG_NY
		default:
			w.CurrentLanguage = LANG_EN
		}
	}

	for key, value := range tree {
		if key == "model" {
			val, ok := value.(string)
			if ok {
				w.CurrentModel = val
			}
		} else {
			row, ok := value.(map[string]any)
			if ok {
				if row["inputIdentifier"] != nil {
					id := fmt.Sprintf("%v", row["inputIdentifier"])

					if row["hidden"] == nil {
						w.ScreenIdMap[id] = key
					}

					if row["readOnly"] != nil {
						w.ReadOnlyFields = append(w.ReadOnlyFields, id)
					}

					if row["formula"] != nil {
						w.FormulaFields[id] = row["formula"].(string)
					}

					if row["scheduleFormula"] != nil {
						w.ScheduleFormulaFields[id] = row["scheduleFormula"].(string)
					}

					if row["matchModel"] != nil {
						w.MatchModelFields[id] = row["matchModel"].(string)
					}

					if row["condition"] != nil {
						w.ConditionFields[id] = row["condition"].(string)
					}

					if row["order"] != nil {
						if row["skipSummary"] != nil {
							val, ok := row["skipSummary"].(bool)
							if ok && val {
								continue
							}
						}

						i, err := strconv.Atoi(fmt.Sprintf("%v", row["order"]))

						if err == nil {
							w.ScreenOrder[i] = id
						}
					}

					if row["optional"] != nil || row["terminateBlockOnEmpty"] != nil {
						w.OptionalFields[id] = true
					}
				}
			}
		}
	}

	return w
}

func (w *WorkFlow) CalculateFormulae(wait chan bool) error {
	defer func() {
		wait <- true
	}()

	for key, value := range w.FormulaFields {
		tokens := GetTokens(value)

		result, err := ResultFromFormulae(tokens, w.Data)
		if err != nil {
			return err
		}

		w.Data[key] = fmt.Sprintf("%0.2f", *result)
	}

	return nil
}

func (w *WorkFlow) FetchModelId(model, field, value string) (*int64, error) {
	if w.Sessions != nil && w.CurrentPhoneNumber != "" {
		session := w.Sessions[w.CurrentPhoneNumber]

		if session.GenericQueryFn != nil {
			result, err := session.GenericQueryFn(fmt.Sprintf(`SELECT id FROM %s WHERE %s = "%s"`, model, field, value))
			if err != nil {
				return nil, err
			}

			if len(result) > 0 && result[0]["id"] != nil {
				id, err := strconv.ParseInt(fmt.Sprintf("%v", result[0]["id"]), 10, 64)
				if err != nil {
					return nil, err
				}

				return &id, nil
			}
		}
	}

	return nil, nil
}

func (w *WorkFlow) EvaluateScheduleFormulae(wait chan bool) error {
	defer func() {
		wait <- true
	}()

	data := map[string]any{}

	if w.Sessions != nil && w.CurrentPhoneNumber != "" {
		session := w.Sessions[w.CurrentPhoneNumber]

		if w.CurrentModel == "memberLoan" {
			if loanRates, ok := session.LoanRates["loanRates"].(map[string]any); ok {
				key := fmt.Sprintf("%v:%v", w.Data["loanType1"], w.Data["loanCategory1"])

				data["loanAmount"] = w.Data["loanAmount1"]
				data["repaymentPeriodInMonths"] = w.Data["repaymentPeriodInMonths1"]

				if row, ok := loanRates[key]; ok {
					if val, ok := row.(map[string]any); ok {
						for _, key := range []string{
							"processingFeeRate",
							"monthlyInterestRate",
							"monthlyInsuranceRate",
						} {
							data[key] = val[key]
						}
					}
				}
			}
		}

		for key, query := range w.ScheduleFormulaFields {
			result, err := GenerateSchedule(query, data)
			if err != nil {
				return err
			}

			w.Data[key] = utils.Map2Table(result, []string{"principal", "totalDue"})
		}
	} else {
		for key, query := range w.ScheduleFormulaFields {
			result, err := GenerateSchedule(query, data)
			if err != nil {
				return err
			}

			w.Data[key] = utils.Map2Table(result, nil)
		}
	}

	return nil
}

func (w *WorkFlow) GetNode(screen string) map[string]any {
	if w.Tree[screen] != nil {
		node, ok := w.Tree[screen].(map[string]any)
		if ok {
			return node
		}
	}

	return nil
}

func (w *WorkFlow) InputIncluded(input string, options []any) (bool, string) {
	var nextRoute string
	found := false

	for _, opt := range options {
		option, ok := opt.(map[string]any)
		if ok && option["position"] != nil {
			var value int

			val, ok := option["position"].(int)
			if ok {
				value = val
			} else {
				val, ok := option["position"].(float64)
				if ok {
					value = int(val)
				}
			}

			if fmt.Sprint(value) == input {
				found = true

				if option["nextScreen"] != nil {
					nextRoute = fmt.Sprintf("%s", option["nextScreen"])
				}
				break
			}
		}
	}

	return found, nextRoute
}

func (w *WorkFlow) NodeOptions(input string) []string {
	options := []string{}

	node := w.GetNode(input)
	if node != nil && node["options"] != nil {
		opts, ok := node["options"].([]any)
		if ok {
			for _, row := range opts {
				optVal, ok := row.(map[string]any)
				if ok {
					position := fmt.Sprintf("%v", optVal["position"])

					activeLang := LANG_EN_LABEL

					if w.CurrentLanguage == LANG_NY {
						activeLang = LANG_NY_LABEL
					}

					val, ok := optVal["label"].(map[string]any)
					if ok {
						if val["all"] != nil {
							entry := fmt.Sprintf("%s. %s", position, val["all"])

							options = append(options, entry)
						} else if w.CurrentLanguage != "" && val[activeLang] != nil {
							entry := fmt.Sprintf("%s. %s", position, val[activeLang])

							options = append(options, entry)
						}
					}
				}
			}
		}
	}

	return options
}

func (w *WorkFlow) CheckLanguage() {
	if w.Data != nil && w.Data["language"] != nil {
		val, ok := w.Data["language"].(string)
		if ok {
			switch val {
			case LANG_NY:
				w.CurrentLanguage = LANG_NY
			default:
				w.CurrentLanguage = LANG_EN
			}
		}
	}
}

func (w *WorkFlow) NextNode(input string) (map[string]any, error) {
	var node map[string]any
	var nextScreen string
	var ok bool

	defer func() {
		w.History[w.HistoryIndex] = w.CurrentScreen
	}()

	switch input {
	case "99":
		// Cancel
		w.Data = map[string]any{}
		w.CurrentScreen = INITIAL_SCREEN
		w.CurrentLanguage = LANG_EN
		w.PreviousScreen = ""
		w.History = map[int]string{}
		w.HistoryIndex = -1

		return nil, nil
	case "0":
		if w.CurrentScreen == "formSummary" {
			// Submit
			if w.SubmitCallback != nil {
				data := w.ResolveData(w.Data, true)

				if w.Data["id"] != nil {
					data["id"] = w.Data["id"]
				}

				err := w.SubmitCallback(
					data, &w.CurrentModel, &w.CurrentPhoneNumber,
					&w.PreferenceFolder, w.AddFunc, w.Sessions, w.Data,
				)
				if err != nil {
					log.Println(err)
					return nil, err
				}
			}

			w.CurrentScreen = INITIAL_SCREEN
			w.CurrentLanguage = LANG_EN
			w.PreviousScreen = ""
			w.History = map[int]string{}
			w.HistoryIndex = -1

			w.Data = map[string]any{}

			return nil, nil
		}
	case "00":
		// Main Menu
		w.Data = map[string]any{}
		w.CurrentScreen = INITIAL_SCREEN
		w.CurrentLanguage = LANG_EN
		w.PreviousScreen = ""
		w.History = map[int]string{}
		w.HistoryIndex = -1

		return nil, nil
	case "98":
		if w.PreviousScreen != "" {
			nextScreen = w.PreviousScreen

			if w.HistoryIndex > 0 {
				w.HistoryIndex--

				prevIndex := w.HistoryIndex - 1

				if val, ok := w.History[prevIndex]; ok {
					w.PreviousScreen = val
				}
			} else {
				w.PreviousScreen = ""
			}

			w.CurrentScreen = nextScreen

			node = w.GetNode(nextScreen)

			if node["scheduleFormula"] != nil {
				wait := make(chan bool, 1)

				err := w.EvaluateScheduleFormulae(wait)
				<-wait
				if err != nil {
					log.Println(err)
				}
			}

			return node, nil
		}
	}

	if input == "01" || input == "02" {
		node = w.GetNode(w.CurrentScreen)

		if node["terminateBlockOnEmpty"] != nil && input == "02" {
			nextScreen = "formSummary"

			node = w.GetNode(nextScreen)
		} else if node["nextScreen"] != nil {
			nextScreen = fmt.Sprintf("%v", node["nextScreen"])

			node = w.GetNode(nextScreen)
		}
	} else {
		if w.CurrentScreen == INITIAL_SCREEN {
			nextScreen, ok = w.Tree[INITIAL_SCREEN].(string)
			if ok {
				node = w.GetNode(nextScreen)
			}
		} else {
			node = w.GetNode(w.CurrentScreen)

			if node["options"] != nil {
				options := node["options"]

				val, ok := options.([]any)
				if ok {
					valid, nextRoute := w.InputIncluded(input, val)

					if !valid {
						return node, nil
					}

					if nextRoute != "" {
						if node["inputIdentifier"] != nil {
							inputIdentifier := fmt.Sprintf("%v", node["inputIdentifier"])

							w.Data[inputIdentifier] = input

							if inputIdentifier == "language" {
								w.CheckLanguage()
							}
						}

						w.PreviousScreen = w.CurrentScreen
						w.CurrentScreen = nextRoute

						node = w.GetNode(w.CurrentScreen)

						w.HistoryIndex++

						return node, nil
					}
				}

				if node != nil && node["inputIdentifier"] != nil {
					inputIdentifier := fmt.Sprintf("%v", node["inputIdentifier"])

					w.Data[inputIdentifier] = input

					if inputIdentifier == "language" {
						w.CheckLanguage()
					}
				}

				if node["nextScreen"] != nil {
					nextScreen = fmt.Sprintf("%v", node["nextScreen"])

					node = w.GetNode(nextScreen)
				}
			} else {
				if node["validationRule"] != nil {
					val, ok := node["validationRule"].(string)
					if ok {
						re := regexp.MustCompile(val)

						if !re.MatchString(input) {
							return node, nil
						}
					}
				}

				if node["optional"] == nil && len(strings.TrimSpace(input)) == 0 {
					return node, nil
				}

				if node != nil && node["inputIdentifier"] != nil {
					inputIdentifier := fmt.Sprintf("%v", node["inputIdentifier"])

					if node["matchModel"] != nil {
						if model, ok := node["matchModel"].(string); ok {
							id, err := w.FetchModelId(model, inputIdentifier, input)
							if err != nil {
								return nil, err
							}

							if id != nil {
								w.Data[fmt.Sprintf("%sId", model)] = *id
							} else {
								return node, nil
							}
						}
					}

					w.Data[inputIdentifier] = input

					if inputIdentifier == "language" {
						w.CheckLanguage()
					}
				}

				if node["nextScreen"] != nil {
					nextScreen = fmt.Sprintf("%v", node["nextScreen"])

					node = w.GetNode(nextScreen)
				}
			}
		}
	}

	if node["scheduleFormula"] != nil {
		wait := make(chan bool, 1)

		err := w.EvaluateScheduleFormulae(wait)
		<-wait
		if err != nil {
			log.Println(err)
		}
	}

	w.PreviousScreen = w.CurrentScreen
	w.CurrentScreen = nextScreen

	w.HistoryIndex++

	if node["condition"] != nil {
		data := w.ResolveData(w.Data, false)

		if !w.EvalCondition(fmt.Sprintf("%v", node["condition"]), data) {
			return w.NextNode("")
		}
	}

	return node, nil
}

func (w *WorkFlow) EvalCondition(condition string, data map[string]any) bool {
	result := false

	if regexp.MustCompile(`=`).MatchString(condition) {
		parts := strings.Split(condition, "=")

		if len(parts) > 1 {
			identifier := parts[0]
			value := parts[1]

			if val, ok := data[identifier]; ok {
				if v, ok := val.(string); ok {
					if regexp.MustCompile(`^[A-Za-z]+$`).MatchString(value) {
						return v == value
					} else if regexp.MustCompile(`^([A-Z]+)\[(.+)\]$`).MatchString(value) {
						var values = []string{}

						parts := regexp.MustCompile(`^([A-Z]+)\[(.+)\]$`).FindAllStringSubmatch(value, -1)[0]

						op := parts[1]

						for key := range strings.SplitSeq(parts[2], ",") {
							values = append(values, key)
						}

						switch op {
						case "IN":
							return slices.Contains(values, v)
						}
					}
				}
			}
		}
	}

	return result
}

func (w *WorkFlow) OptionValue(options []any, input string) (string, *string) {
	var result string
	var code string

	for _, row := range options {
		optVal, ok := row.(map[string]any)
		if ok {
			position := fmt.Sprintf("%v", optVal["position"])

			if position == input {
				val, ok := optVal["label"].(map[string]any)
				if ok {
					if optVal["code"] != nil {
						code = fmt.Sprintf("%v", optVal["code"])
					}

					if val["all"] != nil {
						result = fmt.Sprintf("%v", val["all"])
						break
					} else if val[LANG_EN_LABEL] != nil {
						result = fmt.Sprintf("%s", val[LANG_EN_LABEL])
						break
					}
				}
			}
		}
	}

	return result, &code
}

func (w *WorkFlow) ResolveData(data map[string]any, preferCode bool) map[string]any {
	result := map[string]any{}

	wait := make(chan bool, 1)

	err := w.CalculateFormulae(wait)
	<-wait
	if err != nil {
		log.Println(err)
	}

	formatLabel := func(key string, value any) {
		if !preferCode &&
			regexp.MustCompile(`^[0-9\.\+e]+$`).MatchString(fmt.Sprintf("%v", value)) &&
			!regexp.MustCompile(utils.NUMBER_FORMAT_ESCAPE).MatchString(strings.ToLower(key)) {
			p := message.NewPrinter(language.English)

			var vn float64

			vr, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
			if err == nil {
				vn = vr
			}

			result[key] = p.Sprintf("%f", number.Decimal(vn))
		} else {
			result[key] = value
		}
	}

	for key, value := range data {
		if _, ok := w.ScheduleFormulaFields[key]; ok {
			continue
		}

		nodeId := w.ScreenIdMap[key]

		if nodeId != "" {
			if w.Tree[nodeId] != nil {
				val, ok := w.Tree[nodeId].(map[string]any)
				if ok {
					if val["options"] != nil {
						opts, ok := val["options"].([]any)

						if ok {
							mappedValue, code := w.OptionValue(opts, fmt.Sprintf("%v", value))

							if mappedValue != "" {
								result[key] = mappedValue
							} else {
								result[key] = value
							}

							if code != nil && *code != "" && preferCode {
								result[key] = *code
							}
						}
					} else {
						if strings.ToLower(key) == "password" && !preferCode {
							result[key] = "********"
							result["password.hidden"] = value
						} else {
							formatLabel(key, value)
						}
					}
				}
			}
		} else {
			formatLabel(key, value)
		}
	}

	return result
}

func (w *WorkFlow) LoadLabel(key string) string {
	dispLabel := key

	if w.ScreenIdMap[key] != "" && w.Tree[w.ScreenIdMap[key]] != nil {
		val, ok := w.Tree[w.ScreenIdMap[key]].(map[string]any)
		if ok {
			if val["text"] != nil {
				vl, ok := val["text"].(map[string]any)
				if ok {
					if vl["all"] != nil {
						dispLabel = fmt.Sprintf("%v", vl["all"])
					} else if w.CurrentLanguage != "" {
						langLabel := LANG_EN_LABEL

						switch w.CurrentLanguage {
						case LANG_NY:
							langLabel = LANG_NY_LABEL
						default:
							langLabel = LANG_EN_LABEL
						}

						if vl[langLabel] != nil {
							dispLabel = fmt.Sprintf("%v", vl[langLabel])
						}
					}
				}
			}
		}
	}

	return dispLabel
}

func (w *WorkFlow) GetLabel(node map[string]any, input string) string {
	var label string

	if node != nil {
		var nodeType string
		var startLabel string
		var title string

		if node["type"] != nil {
			nodeType = fmt.Sprintf("%s", node["type"])
		}

		if nodeType == "" {
			return label
		}

		if nodeType == QUIT_SCREEN {
			data := w.ResolveData(w.Data, false)

			result := []string{}

			if w.CurrentLanguage == LANG_NY {
				result = append(result, "Zomwe Mwalemba")
			} else {
				result = append(result, "Summary")
			}

			indices := make([]int, 0, len(w.ScreenOrder))

			for k := range w.ScreenOrder {
				indices = append(indices, k)
			}

			sort.Ints(indices)

			for _, i := range indices {
				key := w.ScreenOrder[i]

				if data[key] != nil {
					dispLabel := w.LoadLabel(key)

					result = append(result, fmt.Sprintf("- %s: %v", dispLabel, data[key]))
				}
			}
			for _, key := range w.ReadOnlyFields {
				if data[key] != nil {
					dispLabel := w.LoadLabel(key)

					result = append(result, fmt.Sprintf("- %s: %v", dispLabel, data[key]))
				}
			}

			if w.CurrentLanguage == LANG_NY {
				result = append(result, "")
				result = append(result, "0. Zatheka")
				result = append(result, "00. Tiyambirenso")
				result = append(result, "98. Bwererani")
				result = append(result, "99. Basi")
			} else {
				result = append(result, "")
				result = append(result, "0. Submit")
				result = append(result, "00. Main Menu")
				result = append(result, "98. Back")
				result = append(result, "99. Cancel")
			}

			label = strings.Join(result, "\n")
		} else {
			if w.Tree[INITIAL_SCREEN] != nil {
				startLabel = fmt.Sprintf("%s", w.Tree[INITIAL_SCREEN])
			}

			var id string

			if node["inputIdentifier"] != nil {
				id = fmt.Sprintf("%v", node["inputIdentifier"])

				dispLabel := w.LoadLabel(id)

				var existingData string

				if w.Data[id] != nil && fmt.Sprintf("%v", w.Data[id]) != "" {
					if strings.ToLower(id) == "password" {
						existingData = "(*******)"
					} else {
						if node["scheduleFormula"] != nil {
							existingData = fmt.Sprintf("\n\n%v", w.Data[id])
						} else {
							var value string

							if regexp.MustCompile(`^[0-9\.\+e]+$`).MatchString(fmt.Sprintf("%v", w.Data[id])) &&
								!regexp.MustCompile(utils.NUMBER_FORMAT_ESCAPE).MatchString(strings.ToLower(id)) {
								p := message.NewPrinter(language.English)

								var vn float64

								vr, err := strconv.ParseFloat(fmt.Sprintf("%v", w.Data[id]), 64)
								if err == nil {
									vn = vr
								}

								value = p.Sprintf("%f", number.Decimal(vn))
							} else {
								value = fmt.Sprintf("%v", w.Data[id])
							}

							existingData = fmt.Sprintf("(%v)", value)
						}
					}
				}

				title = fmt.Sprintf("%s: %s", dispLabel, existingData)
			}

			options := w.NodeOptions(input)
			options = append(options, "")

			if w.CurrentLanguage == LANG_NY {
				if input != startLabel {
					options = append(options, "00. Tiyambirenso")
				}

				if id != "" && w.Data[id] != nil {
					options = append(options, "01. Momwemo")
				}

				if id != "" && w.OptionalFields[id] {
					options = append(options, "02. Tidumphe")
				}

				if input != startLabel {
					options = append(options, "98. Bwererani")
				}

				options = append(options, "99. Basi")
			} else {
				if input != startLabel {
					options = append(options, "00. Main Menu")
				}

				if id != "" && w.Data[id] != nil {
					options = append(options, "01. Keep")
				}

				if id != "" && w.OptionalFields[id] {
					options = append(options, "02. Skip")
				}

				if input != startLabel {
					options = append(options, "98. Back")
				}

				options = append(options, "99. Cancel")
			}

			label = fmt.Sprintf(`%s
%s
`, title, strings.Join(options, "\n"))
		}
	}

	return label
}

func (w *WorkFlow) NavNext(input string) string {
	node, err := w.NextNode(input)
	if err != nil {
		return "Transaction error\n\n00. Main Menu\n"
	}

	label := w.GetLabel(node, w.CurrentScreen)

	return label
}
