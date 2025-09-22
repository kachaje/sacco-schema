CREATE TABLE IF NOT EXISTS account (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    number INTEGER NOT NULL,
    normal INTEGER NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS accountUpdated AFTER
UPDATE ON account FOR EACH ROW BEGIN
UPDATE account
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS accountJournal (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    accountId INTEGER NOT NULL,
    accountTransactionId INTEGER NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amount REAL NOT NULL,
    direction INTEGER NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (accountId) REFERENCES account (id) ON DELETE CASCADE,
    FOREIGN KEY (accountTransactionId) REFERENCES accountTransaction (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS accountJournalUpdated AFTER
UPDATE ON accountJournal FOR EACH ROW BEGIN
UPDATE accountJournal
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS accountStatement (
    accountId INTEGER NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    closingBalance REAL NOT NULL,
    totalCredit REAL NOT NULL,
    totalDebit REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (accountId) REFERENCES account (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS accountStatementUpdated AFTER
UPDATE ON accountStatement FOR EACH ROW BEGIN
UPDATE accountStatement
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS accountTransaction (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    description TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS accountTransactionUpdated AFTER
UPDATE ON accountTransaction FOR EACH ROW BEGIN
UPDATE accountTransaction
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS contributionNumberIdsCache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberContributionId INTEGER,
    idNumber TEXT NOT NULL,
    claimed INTEGER DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberContributionId) REFERENCES memberContribution (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS contributionNumberIdsCacheUpdated AFTER
UPDATE ON contributionNumberIdsCache FOR EACH ROW BEGIN
UPDATE contributionNumberIdsCache
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS insuranceProvider (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    phoneNumber TEXT NOT NULL,
    contactPerson TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS insuranceProviderUpdated AFTER
UPDATE ON insuranceProvider FOR EACH ROW BEGIN
UPDATE insuranceProvider
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS loanNumberIdsCache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER,
    loanNumber TEXT NOT NULL,
    claimed INTEGER DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS loanNumberIdsCacheUpdated AFTER
UPDATE ON loanNumberIdsCache FOR EACH ROW BEGIN
UPDATE loanNumberIdsCache
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS member (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberIdNumber TEXT,
    firstName TEXT NOT NULL,
    lastName TEXT NOT NULL,
    otherName TEXT,
    gender TEXT NOT NULL CHECK (gender IN ('Female', 'Male')),
    phoneNumber TEXT NOT NULL,
    title TEXT NOT NULL CHECK (
        title IN ('Mr', 'Mrs', 'Miss', 'Dr', 'Prof', 'Rev', 'Other')
    ),
    maritalStatus TEXT NOT NULL CHECK (
        maritalStatus IN ('Married', 'Single', 'Widowed', 'Divorced')
    ),
    dateOfBirth TEXT NOT NULL,
    nationalIdentifier TEXT NOT NULL,
    utilityBillType TEXT NOT NULL CHECK (utilityBillType IN ('ESCOM', 'Water Board')),
    utilityBillNumber TEXT NOT NULL,
    dateJoined TEXT DEFAULT CURRENT_TIMESTAMP,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS memberUpdated AFTER
UPDATE ON member FOR EACH ROW BEGIN
UPDATE member
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberBusiness (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    yearsInBusiness INTEGER NOT NULL,
    businessNature TEXT NOT NULL,
    businessName TEXT NOT NULL,
    tradingArea TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberBusinessUpdated AFTER
UPDATE ON memberBusiness FOR EACH ROW BEGIN
UPDATE memberBusiness
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberBusinessVerification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberBusinessId INTEGER NOT NULL,
    businessVerified TEXT DEFAULT No CHECK (businessVerified IN ('Yes', 'No')),
    grossIncomeVerified TEXT DEFAULT No CHECK (grossIncomeVerified IN ('Yes', 'No')),
    netIncomeVerified TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberBusinessId) REFERENCES memberBusiness (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberBusinessVerificationUpdated AFTER
UPDATE ON memberBusinessVerification FOR EACH ROW BEGIN
UPDATE memberBusinessVerification
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberContact (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    postalAddress TEXT NOT NULL,
    residentialAddress TEXT NOT NULL,
    email TEXT,
    homeVillage TEXT NOT NULL,
    homeTraditionalAuthority TEXT NOT NULL,
    homeDistrict TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberContactUpdated AFTER
UPDATE ON memberContact FOR EACH ROW BEGIN
UPDATE memberContact
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberContribution (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    contributionNumber TEXT,
    memberIdNumber TEXT NOT NULL,
    monthlyContribution REAL NOT NULL,
    nonRedeemableAmount REAL NOT NULL,
    totalAmount REAL DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberContributionUpdated AFTER
UPDATE ON memberContribution FOR EACH ROW BEGIN
UPDATE memberContribution
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberContributionDeposit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberContributionId INTEGER NOT NULL,
    description TEXT NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amount REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberContributionId) REFERENCES memberContribution (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberContributionDepositUpdated AFTER
UPDATE ON memberContributionDeposit FOR EACH ROW BEGIN
UPDATE memberContributionDeposit
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberContributionSchedule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberContributionId INTEGER NOT NULL,
    dueDate TEXT NOT NULL,
    expectedAmount REAL NOT NULL,
    paidAmount REAL DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberContributionId) REFERENCES memberContribution (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberContributionScheduleUpdated AFTER
UPDATE ON memberContributionSchedule FOR EACH ROW BEGIN
UPDATE memberContributionSchedule
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberContributionWithdraw (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberSavingId INTEGER NOT NULL,
    description TEXT NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amount REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberSavingId) REFERENCES memberSaving (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberContributionWithdrawUpdated AFTER
UPDATE ON memberContributionWithdraw FOR EACH ROW BEGIN
UPDATE memberContributionWithdraw
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberDependant (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    name TEXT NOT NULL,
    phoneNumber TEXT NOT NULL,
    address TEXT,
    percentage REAL NOT NULL,
    isNominee TEXT DEFAULT No CHECK (isNominee IN ('Yes', 'No')),
    relationship TEXT NOT NULL CHECK (
        relationship IN ('Spouse', 'Child', 'Sibling', 'Other')
    ),
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberDependantUpdated AFTER
UPDATE ON memberDependant FOR EACH ROW BEGIN
UPDATE memberDependant
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberIdsCache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER,
    idNumber TEXT NOT NULL,
    claimed INTEGER DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberIdsCacheUpdated AFTER
UPDATE ON memberIdsCache FOR EACH ROW BEGIN
UPDATE memberIdsCache
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLastYearBusinessHistory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberBusinessId INTEGER NOT NULL,
    financialYear INTEGER NOT NULL,
    totalIncome REAL NOT NULL,
    totalCostOfGoods REAL NOT NULL,
    employeeWages REAL NOT NULL,
    ownSalary REAL NOT NULL,
    transport REAL NOT NULL,
    loanInterest REAL NOT NULL,
    utilities REAL NOT NULL,
    rentals REAL NOT NULL,
    otherCosts REAL NOT NULL,
    totalCosts REAL NOT NULL,
    netProfitLoss REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberBusinessId) REFERENCES memberBusiness (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLastYearBusinessHistoryUpdated AFTER
UPDATE ON memberLastYearBusinessHistory FOR EACH ROW BEGIN
UPDATE memberLastYearBusinessHistory
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    loanNumber TEXT,
    loanStartDate TEXT DEFAULT CURRENT_TIMESTAMP,
    loanDueDate TEXT DEFAULT CURRENT_TIMESTAMP,
    loanPurpose TEXT NOT NULL,
    loanAmount REAL NOT NULL,
    repaymentPeriodInMonths INTEGER NOT NULL,
    loanType TEXT NOT NULL CHECK (
        loanType IN (
            'School fees',
            'Personal',
            'Business',
            'Agricultural',
            'Emergency'
        )
    ),
    loanCategory TEXT DEFAULT Individual CHECK (
        loanCategory IN ('Individual', 'Group/Institution')
    ),
    monthlyInstalments REAL DEFAULT 0,
    monthlyInterestRate REAL DEFAULT 0,
    monthlyInsuranceRate REAL DEFAULT 0,
    processingFeeRate REAL DEFAULT 0,
    penaltyRate REAL DEFAULT 0,
    amountPaid REAL DEFAULT 0,
    balanceAmount REAL DEFAULT 0,
    loanSchedule TEXT,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanUpdated AFTER
UPDATE ON memberLoan FOR EACH ROW BEGIN
UPDATE memberLoan
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanApproval (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    loanNumber TEXT NOT NULL,
    loanStatus TEXT DEFAULT PENDING CHECK (
        loanStatus IN (
            'PENDING',
            'APPROVED',
            'PARTIAL-APPROVAL',
            'REJECTED'
        )
    ),
    amountRecommended REAL NOT NULL,
    denialReason TEXT,
    partialApprovalReason TEXT,
    approvedBy TEXT,
    dateOfApproval TEXT DEFAULT CURRENT_TIMESTAMP,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanApprovalUpdated AFTER
UPDATE ON memberLoanApproval FOR EACH ROW BEGIN
UPDATE memberLoanApproval
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanDisbursement (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    description TEXT NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amount REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanDisbursementUpdated AFTER
UPDATE ON memberLoanDisbursement FOR EACH ROW BEGIN
UPDATE memberLoanDisbursement
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanLiability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    description TEXT NOT NULL,
    value REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanLiabilityUpdated AFTER
UPDATE ON memberLoanLiability FOR EACH ROW BEGIN
UPDATE memberLoanLiability
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanPayment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanPaymentScheduleId INTEGER,
    loanNumber TEXT NOT NULL,
    dueDate TEXT NOT NULL,
    description TEXT NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amountPaid REAL NOT NULL,
    availableCash REAL DEFAULT 0,
    totalDue REAL DEFAULT 0,
    interest REAL DEFAULT 0,
    insurance REAL DEFAULT 0,
    processingFee REAL DEFAULT 0,
    instalment REAL DEFAULT 0,
    penalty REAL DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanPaymentScheduleId) REFERENCES memberLoanPaymentSchedule (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanPaymentUpdated AFTER
UPDATE ON memberLoanPayment FOR EACH ROW BEGIN
UPDATE memberLoanPayment
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanPaymentSchedule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    dueDate TEXT NOT NULL,
    principal REAL NOT NULL,
    interest REAL NOT NULL,
    insurance REAL NOT NULL,
    processingFee REAL DEFAULT 0,
    instalment REAL NOT NULL,
    amountPaid REAL DEFAULT 0,
    amountRecommended REAL NOT NULL,
    loanNumber TEXT NOT NULL,
    penalty REAL DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanPaymentScheduleUpdated AFTER
UPDATE ON memberLoanPaymentSchedule FOR EACH ROW BEGIN
UPDATE memberLoanPaymentSchedule
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanSecurity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    description TEXT NOT NULL,
    value REAL NOT NULL,
    serialNumber TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanSecurityUpdated AFTER
UPDATE ON memberLoanSecurity FOR EACH ROW BEGIN
UPDATE memberLoanSecurity
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanSettlement (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    loanNumber TEXT NOT NULL,
    amountReserved REAL NOT NULL,
    amountClaimed REAL DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS memberLoanSettlementUpdated AFTER
UPDATE ON memberLoanSettlement FOR EACH ROW BEGIN
UPDATE memberLoanSettlement
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanTax (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanPaymentId INTEGER NOT NULL,
    description TEXT NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amount REAL NOT NULL,
    taxCategory TEXT NOT NULL CHECK (
        taxCategory IN ('Interest', 'Processing Fee', 'Penalty')
    ),
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanPaymentId) REFERENCES memberLoanPayment (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanTaxUpdated AFTER
UPDATE ON memberLoanTax FOR EACH ROW BEGIN
UPDATE memberLoanTax
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanType (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK (
        name IN (
            'School fees',
            'Personal',
            'Business',
            'Agricultural',
            'Emergency'
        )
    ),
    category TEXT NOT NULL CHECK (category IN ('Individual', 'Group/Institution')),
    amountLimit REAL NOT NULL,
    periodLimitInMonths INTEGER NOT NULL,
    maxInstalmentMonths INTEGER NOT NULL,
    processingFeeRate REAL DEFAULT 0.05,
    penaltyRate REAL DEFAULT 0.1,
    monthlyInterestRate REAL DEFAULT 0.05,
    monthlyInsuranceRate REAL DEFAULT 0.015,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS memberLoanTypeUpdated AFTER
UPDATE ON memberLoanType FOR EACH ROW BEGIN
UPDATE memberLoanType
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanVerification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanApprovalId INTEGER NOT NULL,
    loanNumber TEXT NOT NULL,
    verified TEXT NOT NULL CHECK (verified IN ('Yes', 'No')),
    verifiedBy TEXT,
    dateVerified TEXT DEFAULT CURRENT_TIMESTAMP,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanApprovalId) REFERENCES memberLoanApproval (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanVerificationUpdated AFTER
UPDATE ON memberLoanVerification FOR EACH ROW BEGIN
UPDATE memberLoanVerification
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberLoanWitness (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    witnessName TEXT NOT NULL,
    telephone TEXT NOT NULL,
    address TEXT,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberLoanWitnessUpdated AFTER
UPDATE ON memberLoanWitness FOR EACH ROW BEGIN
UPDATE memberLoanWitness
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberNextYearBusinessProjection (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberBusinessId INTEGER NOT NULL,
    financialYear INTEGER NOT NULL,
    totalIncome REAL NOT NULL,
    totalCostOfGoods REAL NOT NULL,
    employeeWages REAL NOT NULL,
    ownSalary REAL NOT NULL,
    transport REAL NOT NULL,
    loanInterest REAL NOT NULL,
    utilities REAL NOT NULL,
    rentals REAL NOT NULL,
    otherCosts REAL NOT NULL,
    totalCosts REAL NOT NULL,
    netProfitLoss REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberBusinessId) REFERENCES memberBusiness (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberNextYearBusinessProjectionUpdated AFTER
UPDATE ON memberNextYearBusinessProjection FOR EACH ROW BEGIN
UPDATE memberNextYearBusinessProjection
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberOccupation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    employerName TEXT NOT NULL,
    employerAddress TEXT NOT NULL,
    employerPhone TEXT NOT NULL,
    jobTitle TEXT NOT NULL,
    periodEmployedInMonths INTEGER NOT NULL,
    grossPay REAL NOT NULL,
    netPay REAL NOT NULL,
    highestQualification TEXT NOT NULL CHECK (
        highestQualification IN ('Tertiary', 'Secondary', 'Primary', 'None')
    ),
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberOccupationUpdated AFTER
UPDATE ON memberOccupation FOR EACH ROW BEGIN
UPDATE memberOccupation
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberOccupationVerification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberOccupationId INTEGER NOT NULL,
    jobVerified TEXT DEFAULT No CHECK (jobVerified IN ('Yes', 'No')),
    grossVerified TEXT NOT NULL,
    netVerified TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberOccupationId) REFERENCES memberOccupation (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberOccupationVerificationUpdated AFTER
UPDATE ON memberOccupationVerification FOR EACH ROW BEGIN
UPDATE memberOccupationVerification
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberSaving (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    memberSavingsIdNumber TEXT,
    savingsTypeId INTEGER NOT NULL,
    balance REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE,
    FOREIGN KEY (savingsTypeId) REFERENCES savingsType (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberSavingUpdated AFTER
UPDATE ON memberSaving FOR EACH ROW BEGIN
UPDATE memberSaving
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberSavingDeposit (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberSavingId INTEGER NOT NULL,
    description TEXT NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amount REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberSavingId) REFERENCES memberSaving (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberSavingDepositUpdated AFTER
UPDATE ON memberSavingDeposit FOR EACH ROW BEGIN
UPDATE memberSavingDeposit
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberSavingInterest (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberSavingId INTEGER NOT NULL,
    description TEXT NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amount REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberSavingId) REFERENCES memberSaving (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberSavingInterestUpdated AFTER
UPDATE ON memberSavingInterest FOR EACH ROW BEGIN
UPDATE memberSavingInterest
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberSavingWithdrawal (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberSavingId INTEGER NOT NULL,
    description TEXT NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    amount REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberSavingId) REFERENCES memberSaving (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberSavingWithdrawalUpdated AFTER
UPDATE ON memberSavingWithdrawal FOR EACH ROW BEGIN
UPDATE memberSavingWithdrawal
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS memberSavingsIdsCache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberSavingId INTEGER,
    idNumber TEXT NOT NULL,
    claimed INTEGER DEFAULT 0,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberSavingId) REFERENCES memberSaving (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS memberSavingsIdsCacheUpdated AFTER
UPDATE ON memberSavingsIdsCache FOR EACH ROW BEGIN
UPDATE memberSavingsIdsCache
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS notification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    message TEXT NOT NULL,
    msgDelivered TEXT DEFAULT No CHECK (msgDelivered IN ('Yes', 'No')),
    msgRead TEXT DEFAULT No CHECK (msgRead IN ('Yes', 'No')),
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS notificationUpdated AFTER
UPDATE ON notification FOR EACH ROW BEGIN
UPDATE notification
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS savingsRate (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    savingsTypeId INTEGER NOT NULL,
    monthlyRate REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (savingsTypeId) REFERENCES savingsType (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS savingsRateUpdated AFTER
UPDATE ON savingsRate FOR EACH ROW BEGIN
UPDATE savingsRate
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS savingsType (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    amountSize REAL NOT NULL,
    withdrawPattern TEXT NOT NULL,
    minWithdrawMonths INTEGER NOT NULL,
    maxWithdrawMonths INTEGER NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS savingsTypeUpdated AFTER
UPDATE ON savingsType FOR EACH ROW BEGIN
UPDATE savingsType
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS taxRate (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    value REAL NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS taxRateUpdated AFTER
UPDATE ON taxRate FOR EACH ROW BEGIN
UPDATE taxRate
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    userRole TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS userUpdated AFTER
UPDATE ON user FOR EACH ROW BEGIN
UPDATE user
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;

CREATE TABLE IF NOT EXISTS userRole (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS userRoleUpdated AFTER
UPDATE ON userRole FOR EACH ROW BEGIN
UPDATE userRole
SET
    updatedAt=CURRENT_TIMESTAMP
WHERE
    id=OLD.id;

END;
