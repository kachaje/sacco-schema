package menufuncs

import (
	"fmt"
	"sacco/parser"
	"strings"
)

func ListUsers(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
	session *parser.Session,
) string {
	var response, content string

	title := "Users List\n----------"

	result, err := DB.SQLQuery("SELECT id, username, userRole, createdAt, updatedAt FROM user WHERE active = 1")
	if err != nil {
		content = fmt.Sprintf("%s\n", err.Error())
	} else {
		rows := []string{
			fmt.Sprintf("%2s | %-10s | %-8s | %-20s | %-20s", "id", "username", "role", "createdAt", "updatedAt"),
			strings.Repeat("-", 70),
		}

		for _, row := range result {
			id := fmt.Sprintf("%v", row["id"])
			username := fmt.Sprintf("%v", row["username"])
			userRole := fmt.Sprintf("%v", row["userRole"])
			createdAt := fmt.Sprintf("%v", row["createdAt"])
			updatedAt := fmt.Sprintf("%v", row["updatedAt"])

			entry := fmt.Sprintf("%2s | %-10s | %-8s | %-20s | %-20s", id, username, userRole, createdAt, updatedAt)

			rows = append(rows, entry)
		}

		content = strings.Join(rows, "\n")
	}

	response = fmt.Sprintf("%s\n\n%s\n\n00. Main Menu\n", title, content)

	return response
}
