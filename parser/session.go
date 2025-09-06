package parser

import (
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

	QueryFn    func(string, []string) (map[string]any, error)
	SkipFields []string

	Mu *sync.Mutex

	SessionToken    *string
	SessionUser     *string
	SessionUserRole *string
	SessionUserId   *int64

	Cache      map[string]string
	LastPrompt string
}

func NewSession(
	queryFn func(string, []string) (map[string]any, error),
	phoneNumber, sessionId *string,
) *Session {
	s := &Session{
		QueryFn:          queryFn,
		Mu:               &sync.Mutex{},
		AddedModels:      map[string]bool{},
		ActiveData:       map[string]any{},
		Data:             map[string]string{},
		SkipFields:       []string{"active"},
		CurrentMenu:      "main",
		WorkflowsMapping: map[string]*WorkFlow{},
		Cache:            map[string]string{},
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
