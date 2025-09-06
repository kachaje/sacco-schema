package menufuncs

import (
	"log"
	"os"
	"path/filepath"
	"sacco/parser"
	"sync"
)

var mu sync.Mutex

func DoExit(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	mu.Lock()
	defer mu.Unlock()

	var phoneNumber string
	var cacheFolder string

	if data != nil {
		if data["session"] != nil {
			if val, ok := data["session"].(*parser.Session); ok {
				session = val
			}
		}
		if data["phoneNumber"] != nil {
			if val, ok := data["phoneNumber"].(string); ok {
				phoneNumber = val
			}
		}
		if data["cacheFolder"] != nil {
			if val, ok := data["cacheFolder"].(string); ok {
				cacheFolder = val
			}
		}

		session.Cache = map[string]string{}
		session.LastPrompt = "username"
		session.SessionToken = nil

		if phoneNumber != "" {
			delete(Sessions, phoneNumber)

			if cacheFolder != "" {
				folderName := filepath.Join(cacheFolder, phoneNumber)

				_, err := os.Stat(folderName)
				if !os.IsNotExist(err) {
					files, err := os.ReadDir(folderName)
					if err == nil && len(files) == 0 {
						err = os.RemoveAll(folderName)
						if err != nil {
							log.Printf("server.menus.menu.removeFolder: %s\n", err.Error())
						}
					}
				}
			}
		}
	}

	return "END Thank you for using our service"
}
