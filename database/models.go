package database

var (
	AccountArrayChildren = []string{
		"accountJournal",
		"accountStatement",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	LoanNumberIdsCacheSingleChildren = []string{
		"memberLoan",
	}
	MemberArrayChildren = []string{
		"memberContribution",
		"memberDependant",
		"memberLoan",
		"memberSaving",
		"notification",
	}
	MemberBusinessSingleChildren = []string{
		"memberBusinessVerification",
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	MemberContributionIdsCacheSingleChildren = []string{
		"memberContribution",
	}
	MemberIdsCacheSingleChildren = []string{
		"member",
	}
	MemberLoanApprovalSingleChildren = []string{
		"memberLoanVerification",
	}
	MemberLoanArrayChildren = []string{
		"memberLoanLiability",
		"memberLoanPaymentSchedule",
		"memberLoanSecurity",
		"memberLoanWitness",
	}
	MemberLoanPaymentArrayChildren = []string{
		"memberLoanTax",
	}
	MemberLoanPaymentScheduleSingleChildren = []string{
		"memberLoanPayment",
	}
	MemberLoanSingleChildren = []string{
		"memberBusiness",
		"memberLoanApproval",
		"memberLoanDisbursement",
		"memberOccupation",
	}
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
	}
	MemberSavingArrayChildren = []string{
		"contributionWithdraw",
		"memberSavingDeposit",
		"memberSavingInterest",
		"memberSavingWithdrawal",
	}
	MemberSavingsIdsCacheSingleChildren = []string{
		"memberSaving",
	}
	MemberSingleChildren = []string{
		"memberContact",
	}
	SavingsTypeArrayChildren = []string{
		"memberSaving",
		"savingsRate",
	}
	SingleChildren = map[string][]string{
		"LoanNumberIdsCacheSingleChildren":         LoanNumberIdsCacheSingleChildren,
		"MemberBusinessSingleChildren":             MemberBusinessSingleChildren,
		"MemberContributionIdsCacheSingleChildren": MemberContributionIdsCacheSingleChildren,
		"MemberIdsCacheSingleChildren":             MemberIdsCacheSingleChildren,
		"MemberLoanApprovalSingleChildren":         MemberLoanApprovalSingleChildren,
		"MemberLoanPaymentScheduleSingleChildren":  MemberLoanPaymentScheduleSingleChildren,
		"MemberLoanSingleChildren":                 MemberLoanSingleChildren,
		"MemberOccupationSingleChildren":           MemberOccupationSingleChildren,
		"MemberSavingsIdsCacheSingleChildren":      MemberSavingsIdsCacheSingleChildren,
		"MemberSingleChildren":                     MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberLoanArrayChildren":         MemberLoanArrayChildren,
		"MemberLoanPaymentArrayChildren":  MemberLoanPaymentArrayChildren,
		"MemberSavingArrayChildren":       MemberSavingArrayChildren,
		"SavingsTypeArrayChildren":        SavingsTypeArrayChildren,
	}
	FloatFields = []string{
		"amount",
		"amountClaimed",
		"amountLimit",
		"amountRecommended",
		"amountReserved",
		"amountSize",
		"balance",
		"claimed",
		"closingBalance",
		"direction",
		"employeeWages",
		"financialYear",
		"grossPay",
		"instalment",
		"insurance",
		"interest",
		"loanAmount",
		"loanInterest",
		"maxInstalmentMonths",
		"maxWithdrawMonths",
		"minWithdrawMonths",
		"monthlyContribution",
		"monthlyInstalments",
		"monthlyInsuranceRate",
		"monthlyInterestRate",
		"monthlyRate",
		"netPay",
		"netProfitLoss",
		"normal",
		"number",
		"otherCosts",
		"ownSalary",
		"penalty",
		"penaltyRate",
		"percentage1",
		"percentage2",
		"percentage3",
		"percentage4",
		"periodEmployedInMonths",
		"periodLimitInMonths",
		"principal",
		"processingFee",
		"processingFeeRate",
		"rentals",
		"repaymentPeriodInMonths",
		"totalCostOfGoods",
		"totalCosts",
		"totalCredit",
		"totalDebit",
		"totalIncome",
		"transport",
		"utilities",
		"value",
		"value1",
		"value2",
		"yearsInBusiness",
	}
	ParentModels = map[string][]string{
		"accountJournal": {
			"account",
			"accountTransaction",
		},
		"accountStatement": {
			"account",
		},
		"contributionWithdraw": {
			"memberSaving",
		},
		"loanNumberIdsCache": {
			"memberLoan",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberBusinessVerification": {
			"memberBusiness",
		},
		"memberContact": {
			"member",
		},
		"memberContribution": {
			"member",
		},
		"memberContributionIdsCache": {
			"memberContribution",
		},
		"memberDependant": {
			"member",
		},
		"memberIdsCache": {
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
		"memberLoanDisbursement": {
			"memberLoan",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanPayment": {
			"memberLoanPaymentSchedule",
		},
		"memberLoanPaymentSchedule": {
			"memberLoan",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberLoanTax": {
			"memberLoanPayment",
		},
		"memberLoanVerification": {
			"memberLoanApproval",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"memberSaving": {
			"member",
			"savingsType",
		},
		"memberSavingDeposit": {
			"memberSaving",
		},
		"memberSavingInterest": {
			"memberSaving",
		},
		"memberSavingWithdrawal": {
			"memberSaving",
		},
		"memberSavingsIdsCache": {
			"memberSaving",
		},
		"notification": {
			"member",
		},
		"savingsRate": {
			"savingsType",
		},
	}
)
