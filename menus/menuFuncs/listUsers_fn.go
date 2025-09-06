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

	result, err := DB.SQLQuery("SELECT id, username, role, created_at, updated_at FROM user WHERE active = 1")
	if err != nil {
		content = fmt.Sprintf("%s\n", err.Error())
	} else {
		rows := []string{
			fmt.Sprintf("%2s | %-10s | %-8s | %-20s | %-20s", "id", "username", "role", "created_at", "updated_at"),
			strings.Repeat("-", 70),
		}

		for _, row := range result {
			id := fmt.Sprintf("%v", row["id"])
			username := fmt.Sprintf("%v", row["username"])
			role := fmt.Sprintf("%v", row["role"])
			createdAt := fmt.Sprintf("%v", row["created_at"])
			updatedAt := fmt.Sprintf("%v", row["updated_at"])

			entry := fmt.Sprintf("%2s | %-10s | %-8s | %-20s | %-20s", id, username, role, createdAt, updatedAt)

			rows = append(rows, entry)
		}

		content = strings.Join(rows, "\n")
	}

	response = fmt.Sprintf("%s\n\n%s\n\n00. Main Menu\n", title, content)

	return response
}
