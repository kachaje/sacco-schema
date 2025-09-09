package database

var (
	AccountArrayChildren = []string{
		"accountJournal",
		"accountStatement",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	InsuranceProviderArrayChildren = []string{
		"memberLoanInsurance",
	}
	LoanNumberIdsCacheSingleChildren = []string{
		"memberLoan",
	}
	LoanTypeSingleChildren = []string{
		"loanRate",
	}
	MemberArrayChildren = []string{
		"dividends",
		"memberDependant",
		"memberLoan",
		"memberSaving",
		"memberShares",
		"notification",
	}
	MemberBusinessSingleChildren = []string{
		"memberBusinessVerification",
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	MemberIdsCacheSingleChildren = []string{
		"member",
	}
	MemberLoanArrayChildren = []string{
		"memberLoanInsurance",
		"memberLoanLiability",
		"memberLoanPaymentSchedule",
		"memberLoanProcessingFee",
		"memberLoanSecurity",
		"memberLoanTax",
		"memberLoanWitness",
	}
	MemberLoanPaymentScheduleArrayChildren = []string{
		"memberLoanReceipt",
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
		"memberSavingDeposit",
		"memberSavingInterest",
		"memberSavingWithdrawal",
		"sharesDepositReceipt",
		"sharesDepositWithdraw",
	}
	MemberSavingsIdsCacheSingleChildren = []string{
		"memberSaving",
	}
	MemberSharesIdsCacheSingleChildren = []string{
		"memberShares",
	}
	MemberSingleChildren = []string{
		"memberContact",
	}
	SavingsTypeArrayChildren = []string{
		"memberSaving",
		"savingsRate",
	}
	SingleChildren = map[string][]string{
		"LoanNumberIdsCacheSingleChildren":    LoanNumberIdsCacheSingleChildren,
		"LoanTypeSingleChildren":              LoanTypeSingleChildren,
		"MemberBusinessSingleChildren":        MemberBusinessSingleChildren,
		"MemberIdsCacheSingleChildren":        MemberIdsCacheSingleChildren,
		"MemberLoanSingleChildren":            MemberLoanSingleChildren,
		"MemberOccupationSingleChildren":      MemberOccupationSingleChildren,
		"MemberSavingsIdsCacheSingleChildren": MemberSavingsIdsCacheSingleChildren,
		"MemberSharesIdsCacheSingleChildren":  MemberSharesIdsCacheSingleChildren,
		"MemberSingleChildren":                MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":                   AccountArrayChildren,
		"AccountTransactionArrayChildren":        AccountTransactionArrayChildren,
		"InsuranceProviderArrayChildren":         InsuranceProviderArrayChildren,
		"MemberArrayChildren":                    MemberArrayChildren,
		"MemberLoanArrayChildren":                MemberLoanArrayChildren,
		"MemberLoanPaymentScheduleArrayChildren": MemberLoanPaymentScheduleArrayChildren,
		"MemberSavingArrayChildren":              MemberSavingArrayChildren,
		"SavingsTypeArrayChildren":               SavingsTypeArrayChildren,
	}
	FloatFields = []string{
		"amount",
		"amount",
		"amount",
		"amount",
		"amount",
		"amount",
		"amount",
		"amount",
		"amount",
		"amount",
		"amount",
		"amount",
		"amountDue",
		"amountLimit",
		"amountRecommended",
		"amountSize",
		"balance",
		"balanceAmount",
		"claimed",
		"claimed",
		"claimed",
		"claimed",
		"closingBalance",
		"direction",
		"employeeWages",
		"employeeWages",
		"financialYear",
		"financialYear",
		"grossPay",
		"interestRate",
		"loanAmount",
		"loanInterest",
		"loanInterest",
		"maxInstalmentMonths",
		"maxWithdrawMonths",
		"minWithdrawMonths",
		"monthlyInstalments",
		"monthlyPremium",
		"monthlyRate",
		"monthlyRate",
		"netPay",
		"netProfitLoss",
		"netProfitLoss",
		"normal",
		"number",
		"numberOfShares",
		"otherCosts",
		"otherCosts",
		"ownSalary",
		"ownSalary",
		"percentage",
		"periodEmployedInMonths",
		"periodLimitInMonths",
		"rentals",
		"rentals",
		"repaymentPeriodInMonths",
		"totalCostOfGoods",
		"totalCostOfGoods",
		"totalCosts",
		"totalCosts",
		"totalCredit",
		"totalDebit",
		"totalIncome",
		"totalIncome",
		"totalValue",
		"transport",
		"transport",
		"utilities",
		"utilities",
		"value",
		"value",
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
		"dividends": {
			"member",
		},
		"loanNumberIdsCache": {
			"memberLoan",
		},
		"loanRate": {
			"loanType",
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
		"memberLoanInsurance": {
			"memberLoan",
			"insuranceProvider",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanPaymentSchedule": {
			"memberLoan",
		},
		"memberLoanProcessingFee": {
			"memberLoan",
		},
		"memberLoanReceipt": {
			"memberLoanPaymentSchedule",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberLoanTax": {
			"memberLoan",
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
		"memberShares": {
			"member",
		},
		"memberSharesIdsCache": {
			"memberShares",
		},
		"notification": {
			"member",
		},
		"savingsRate": {
			"savingsType",
		},
		"sharesDepositReceipt": {
			"memberSaving",
		},
		"sharesDepositWithdraw": {
			"memberSaving",
		},
	}
)
