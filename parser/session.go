package parser

import (
	"fmt"
	"log"
	"regexp"
	"sacco/utils"
	"sync"
	"time"
)

type Session struct {
	CurrentMenu string
	Data        map[string]string

	PreferredLanguage  string
	SessionId          string
	PhoneNumber        string
	CurrentPhoneNumber string

	GlobalIds map[string]any

	WorkflowsMapping map[string]*WorkFlow

	AddedModels map[string]bool

	ActiveData map[string]any
	LoanRates  map[string]any

	QueryFn        func(string, []string) (map[string]any, error)
	GenericQueryFn func(query string) ([]map[string]any, error)
	SkipFields     []string

	Mu *sync.Mutex

	SessionToken    *string
	SessionUser     *string
	SessionUserRole *string
	SessionUserId   *int64

	Cache      map[string]any
	LastPrompt string
}

func NewSession(
	queryFn func(string, []string) (map[string]any, error),
	phoneNumber, sessionId *string,
	genericQueryFn func(query string) ([]map[string]any, error),
) *Session {
	s := &Session{
		QueryFn:          queryFn,
		GenericQueryFn:   genericQueryFn,
		Mu:               &sync.Mutex{},
		AddedModels:      map[string]bool{},
		ActiveData:       map[string]any{},
		LoanRates:        map[string]any{},
		Data:             map[string]string{},
		SkipFields:       []string{"active", "createdAt", "updatedAt"},
		CurrentMenu:      "main",
		WorkflowsMapping: map[string]*WorkFlow{},
		Cache:            map[string]any{},
		LastPrompt:       "",
		GlobalIds:        map[string]any{},
	}

	if phoneNumber != nil {
		s.CurrentPhoneNumber = *phoneNumber
	}
	if sessionId != nil {
		s.SessionId = *sessionId
	}

	return s
}

func (s *Session) updateActiveData(data map[string]any, retries int) {
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			s.updateActiveData(data, retries)
			return
		}
	}
	defer s.Mu.Unlock()

	s.ActiveData = utils.FlattenMap(data, false)

	idsData := utils.FlattenMap(data, true)

	s.GlobalIds = idsData

	s.AddedModels = map[string]bool{}

	for key := range idsData {
		model := key[:len(key)-2]

		model = regexp.MustCompile(`[^A-Za-z]`).ReplaceAllLiteralString(model, "")

		s.AddedModels[model] = true
	}
}

func (s *Session) WriteToMap(key string, value any, retries int) {
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			s.WriteToMap(key, value, retries)
			return
		}
	}
	defer s.Mu.Unlock()

	if s.ActiveData == nil {
		s.ActiveData = map[string]any{}
	}

	s.ActiveData[key] = value
}

func (s *Session) ReadFromMap(key string, retries int) any {
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			return s.ReadFromMap(key, retries)
		}
	}
	defer s.Mu.Unlock()

	return s.ActiveData[key]
}

func (s *Session) ClearSession() {
	s.ActiveData = map[string]any{}
	s.Data = map[string]string{}
	s.AddedModels = map[string]bool{}
	s.GlobalIds = map[string]any{}
}

func (s *Session) RefreshSession() (map[string]any, error) {
	if s.CurrentPhoneNumber != "" && s.QueryFn != nil {
		if s.GenericQueryFn != nil {
			rows, err := s.GenericQueryFn("SELECT * FROM memberLoanType WHERE active = 1")
			if err != nil {
				s.LoanRates = map[string]any{}

				log.Println(err)
			} else {
				loanRates := map[string]any{}

				for _, row := range rows {
					key := fmt.Sprintf("%v:%v", row["name"], row["category"])

					loanRates[key] = row
				}

				s.LoanRates = map[string]any{
					"loanRates": loanRates,
				}

			}
		}

		data, err := s.QueryFn(s.CurrentPhoneNumber, s.SkipFields)
		if err != nil {
			s.updateActiveData(map[string]any{}, 0)

			return nil, err
		}

		s.updateActiveData(data, 0)

		return data, nil
	}
	return s.ActiveData, nil
}
