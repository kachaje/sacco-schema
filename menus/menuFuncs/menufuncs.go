package menufuncs

import (
	"sacco/database"
	"sacco/parser"
)

var (
	DB       *database.Database
	Sessions = map[string]*parser.Session{}

	WorkflowsData = map[string]map[string]any{}

	DemoMode bool

	FunctionsMap = map[string]func(
		func(
			string, *parser.Session,
			string, string, string,
		) string,
		map[string]any,
		*parser.Session,
	) string{}

	ReRouteRemaps = map[string]any{}
)

func init() {
	FunctionsMap["bankingDetails"] = BankingDetails
	FunctionsMap["blockUser"] = BlockUser
	FunctionsMap["businessSummary"] = BusinessSummary
	FunctionsMap["changePassword"] = ChangePassword
	FunctionsMap["checkBalance"] = CheckBalance
	FunctionsMap["devConsole"] = DevConsole
	FunctionsMap["doExit"] = DoExit
	FunctionsMap["editUser"] = EditUser
	FunctionsMap["employmentSummary"] = EmploymentSummary
	FunctionsMap["landing"] = Landing
	FunctionsMap["listUsers"] = ListUsers
	FunctionsMap["memberLoansSummary"] = MemberLoansSummary
	FunctionsMap["setPhoneNumber"] = SetPhoneNumber
	FunctionsMap["signIn"] = SignIn
	FunctionsMap["signUp"] = SignUp
	FunctionsMap["viewMemberDetails"] = ViewMemberDetails
}
