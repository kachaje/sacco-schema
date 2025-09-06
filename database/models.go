package database

var (
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	MemberArrayChildren = []string{
		"memberBeneficiary",
		"memberShares",
		"memberLoan",
	}
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	MemberLoanArrayChildren = []string{
		"memberLoanLiability",
		"memberLoanSecurity",
	}
	MemberLoanSingleChildren = []string{
		"memberBusiness",
		"memberOccupation",
		"memberLoanWitness",
		"memberLoanApproval",
	}
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
	}
	MemberSingleChildren = []string{
		"memberContact",
		"memberDependant",
	}
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
		"MemberSingleChildren":           MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"MemberLoanArrayChildren":         MemberLoanArrayChildren,
	}
	FloatFields = []string{
		"amountRecommended",
		"amountRecommended",
		"credit",
		"debit",
		"employeeWages",
		"employeeWages",
		"financialYear",
		"financialYear",
		"grossPay",
		"loanAmount",
		"loanInterest",
		"loanInterest",
		"netPay",
		"netProfitLoss",
		"netProfitLoss",
		"numberOfShares",
		"otherCosts",
		"otherCosts",
		"ownSalary",
		"ownSalary",
		"password",
		"periodEmployedInMonths",
		"pricePerShare",
		"rentals",
		"rentals",
		"repaymentPeriodInMonths",
		"totalCostOfGoods",
		"totalCostOfGoods",
		"totalCosts",
		"totalCosts",
		"totalIncome",
		"totalIncome",
		"transport",
		"transport",
		"utilities",
		"utilities",
		"value",
		"value",
	}
	ParentModels = map[string][]string{
		"accountJournal": {
			"account",
		},
		"accountTransaction": {
			"account",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberContact": {
			"member",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberLoan": {
			"member",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberDependant": {
			"member",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"memberShares": {
			"member",
		},
	}
)
