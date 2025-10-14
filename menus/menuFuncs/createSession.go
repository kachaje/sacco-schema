package menufuncs

import (
	filehandling "github.com/kachaje/sacco-schema/fileHandling"
	"github.com/kachaje/sacco-schema/parser"
)

func CreateNewSession(phoneNumber, sessionId, preferencesFolder, preferredLanguage string, demoMode bool) *parser.Session {
	mu.Lock()
	session, exists := Sessions[phoneNumber]
	if !exists {
		session = parser.NewSession(
			DB.MemberByPhoneNumber,
			&phoneNumber,
			&sessionId,
			DB.SQLQuery,
		)

		for model, data := range WorkflowsData {
			session.WorkflowsMapping[model] = parser.NewWorkflow(data, filehandling.SaveData, &preferredLanguage, &phoneNumber, &sessionId, &preferencesFolder, DB.GenericsSaveData, Sessions, nil)
		}

		if preferredLanguage != "" {
			session.PreferredLanguage = preferredLanguage
		}

		if demoMode {
			defaultUser := "admin"
			defaultUserId := int64(1)
			defaultRole := "Admin"

			session.SessionUser = &defaultUser
			session.SessionUserId = &defaultUserId
			session.SessionUserRole = &defaultRole
		}

		session.CurrentPhoneNumber = phoneNumber

		Sessions[phoneNumber] = session
	}
	mu.Unlock()

	return session
}
