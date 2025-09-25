package database

var (
	AccountArrayChildren = []string{
		"accountJournal",
		"accountStatement",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	ContributionNumberIdsCacheSingleChildren = []string{
		"memberContribution",
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
	MemberContributionArrayChildren = []string{
		"memberContributionDeposit",
		"memberContributionSchedule",
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
		"memberContributionWithdraw",
		"memberSavingInterest",
		"memberSavingTransaction",
	}
	MemberSavingIdsCacheSingleChildren = []string{
		"memberSaving",
	}
	MemberSingleChildren = []string{
		"memberContact",
	}
	SavingsTypeArrayChildren = []string{
		"memberSaving",
	}
	SingleChildren = map[string][]string{
		"ContributionNumberIdsCacheSingleChildren": ContributionNumberIdsCacheSingleChildren,
		"LoanNumberIdsCacheSingleChildren":         LoanNumberIdsCacheSingleChildren,
		"MemberBusinessSingleChildren":             MemberBusinessSingleChildren,
		"MemberIdsCacheSingleChildren":             MemberIdsCacheSingleChildren,
		"MemberLoanApprovalSingleChildren":         MemberLoanApprovalSingleChildren,
		"MemberLoanPaymentScheduleSingleChildren":  MemberLoanPaymentScheduleSingleChildren,
		"MemberLoanSingleChildren":                 MemberLoanSingleChildren,
		"MemberOccupationSingleChildren":           MemberOccupationSingleChildren,
		"MemberSavingIdsCacheSingleChildren":       MemberSavingIdsCacheSingleChildren,
		"MemberSingleChildren":                     MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberContributionArrayChildren": MemberContributionArrayChildren,
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
		"balance",
		"claimed",
		"closingBalance",
		"deposit",
		"direction",
		"employeeWages",
		"expectedAmount",
		"financialYear",
		"grossPay",
		"instalment",
		"insurance",
		"interest",
		"interestRate",
		"loanAmount",
		"loanInterest",
		"maxInstalmentMonths",
		"maxWithdrawMonths",
		"minWithdrawMonths",
		"minimumAmount",
		"monthlyContribution",
		"monthlyInstalments",
		"monthlyInsuranceRate",
		"monthlyInterestRate",
		"netPay",
		"netProfitLoss",
		"nonRedeemableAmount",
		"normal",
		"number",
		"otherCosts",
		"overflowAmount",
		"ownSalary",
		"paidAmount",
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
		"withdrawal",
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
		"contributionNumberIdsCache": {
			"memberContribution",
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
		"memberContributionDeposit": {
			"memberContribution",
		},
		"memberContributionSchedule": {
			"memberContribution",
		},
		"memberContributionWithdraw": {
			"memberSaving",
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
		"memberSavingIdsCache": {
			"memberSaving",
		},
		"memberSavingInterest": {
			"memberSaving",
		},
		"memberSavingTransaction": {
			"memberSaving",
		},
		"notification": {
			"member",
		},
	}
)
