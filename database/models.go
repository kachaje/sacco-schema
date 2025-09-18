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
	MemberLoanApprovalSingleChildren = []string{
		"memberLoanVerification",
	}
	MemberLoanArrayChildren = []string{
		"memberLoanInsurance",
		"memberLoanLiability",
		"memberLoanPaymentSchedule",
		"memberLoanSecurity",
		"memberLoanWitness",
	}
	MemberLoanInvoiceArrayChildren = []string{
		"memberLoanInvoiceDetail",
		"memberLoanRepayment",
		"memberLoanTax",
	}
	MemberLoanPaymentScheduleSingleChildren = []string{
		"memberLoanInvoice",
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
		"LoanNumberIdsCacheSingleChildren":        LoanNumberIdsCacheSingleChildren,
		"MemberBusinessSingleChildren":            MemberBusinessSingleChildren,
		"MemberIdsCacheSingleChildren":            MemberIdsCacheSingleChildren,
		"MemberLoanApprovalSingleChildren":        MemberLoanApprovalSingleChildren,
		"MemberLoanPaymentScheduleSingleChildren": MemberLoanPaymentScheduleSingleChildren,
		"MemberLoanSingleChildren":                MemberLoanSingleChildren,
		"MemberOccupationSingleChildren":          MemberOccupationSingleChildren,
		"MemberSavingsIdsCacheSingleChildren":     MemberSavingsIdsCacheSingleChildren,
		"MemberSharesIdsCacheSingleChildren":      MemberSharesIdsCacheSingleChildren,
		"MemberSingleChildren":                    MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"InsuranceProviderArrayChildren":  InsuranceProviderArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberLoanArrayChildren":         MemberLoanArrayChildren,
		"MemberLoanInvoiceArrayChildren":  MemberLoanInvoiceArrayChildren,
		"MemberSavingArrayChildren":       MemberSavingArrayChildren,
		"SavingsTypeArrayChildren":        SavingsTypeArrayChildren,
	}
	FloatFields = []string{
		"amount",
		"amountAllocated",
		"amountLimit",
		"amountRecommended",
		"amountSize",
		"balance",
		"billedAmount",
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
		"monthlyInstalments",
		"monthlyInsuranceRate",
		"monthlyInterestRate",
		"monthlyPremium",
		"monthlyRate",
		"netPay",
		"netProfitLoss",
		"normal",
		"number",
		"numberOfShares",
		"otherCosts",
		"ownSalary",
		"paidAmount",
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
		"totalDue",
		"totalIncome",
		"totalValue",
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
		"dividends": {
			"member",
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
		"memberLoanInvoice": {
			"memberLoanPaymentSchedule",
		},
		"memberLoanInvoiceDetail": {
			"memberLoanInvoice",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanPaymentSchedule": {
			"memberLoan",
		},
		"memberLoanRepayment": {
			"memberLoanInvoice",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberLoanTax": {
			"memberLoanInvoice",
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
