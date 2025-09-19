INSERT INTO
  member (
    id,
    firstName,
    lastName,
    otherName,
    gender,
    title,
    maritalStatus,
    dateOfBirth,
    nationalIdentifier,
    utilityBillType,
    utilityBillNumber,
    phoneNumber
  )
VALUES
  (
    1,
    "Mary",
    "Banda",
    "",
    "Female",
    "Miss",
    "Single",
    "1999-09-01",
    "DHFYR8475",
    "ESCOM",
    "29383746",
    "09999999999"
  );

INSERT INTO
  memberContact (
    memberId,
    postalAddress,
    residentialAddress,
    homeVillage,
    homeTraditionalAuthority,
    homeDistrict
  )
VALUES
  (
    1,
    "P.O. Box 3200, Blantyre",
    "Chilomoni, Blantrye",
    "Thumba",
    "Kabudula",
    "Lilongwe"
  );

INSERT INTO
  memberContact (
    memberId,
    postalAddress,
    residentialAddress,
    homeVillage,
    homeTraditionalAuthority,
    homeDistrict
  )
VALUES
  (
    1,
    "P.O. Box 1000, Lilongwe",
    "Area 2, Lilongwe",
    "Songwe",
    "Kyungu",
    "Karonga"
  );

INSERT INTO
  memberLoan (
    memberId,
    loanAmount,
    repaymentPeriodInMonths,
    loanPurpose,
    loanType,
    loanDueDate,
    loanStartDate
  )
VALUES
  (
    1,
    200000,
    12,
    "School fees",
    "Personal",
    "2025-09-08 07:09:32",
    "2025-09-08 07:09:32"
  );

INSERT INTO
  memberBusiness (
    memberLoanId,
    yearsInBusiness,
    businessNature,
    businessName,
    tradingArea
  )
VALUES
  (1, 3, "Vendor", "Vendors Galore", "Mtandire");

INSERT INTO
  memberLastYearBusinessHistory (
    memberBusinessId,
    financialYear,
    totalIncome,
    totalCostOfGoods,
    employeeWages,
    ownSalary,
    transport,
    loanInterest,
    utilities,
    rentals,
    otherCosts,
    totalCosts,
    netProfitLoss
  )
VALUES
  (
    1,
    2024,
    2000000,
    1000000,
    50000,
    100000,
    50000,
    0,
    35000,
    50000,
    0,
    1285000,
    715000
  );

INSERT INTO
  memberNextYearBusinessProjection (
    memberBusinessId,
    financialYear,
    totalIncome,
    totalCostOfGoods,
    employeeWages,
    ownSalary,
    transport,
    loanInterest,
    utilities,
    rentals,
    otherCosts,
    totalCosts,
    netProfitLoss
  )
VALUES
  (
    1,
    2025,
    2500000,
    1500000,
    50000,
    100000,
    50000,
    0,
    35000,
    50000,
    0,
    1285000,
    715000
  );

INSERT INTO
  memberOccupation (
    memberLoanId,
    employerName,
    grossPay,
    netPay,
    jobTitle,
    employerAddress,
    employerPhone,
    periodEmployedInMonths,
    highestQualification
  )
VALUES
  (
    1,
    "SOBO",
    100000,
    90000,
    "Driver",
    "Kanengo",
    "0999888474",
    36,
    "Secondary"
  );

INSERT INTO
  memberOccupationVerification (
    memberOccupationId,
    jobVerified,
    grossVerified,
    netVerified
  )
VALUES
  (1, "Yes", "Yes", "Yes");

INSERT INTO
  memberDependant (
    memberId,
    name,
    phoneNumber,
    address,
    percentage,
    isNominee,
    relationship
  )
VALUES
  (
    1,
    "Benefator 1",
    "0888888888",
    "P.O. Box 1",
    10,
    "Yes",
    "Spouse"
  ),
  (
    1,
    "Benefator 2",
    "0888888887",
    "P.O. Box 2",
    8,
    "No",
    "Child"
  ),
  (
    1,
    "Benefator 3",
    "0888888886",
    "P.O. Box 3",
    5,
    "No",
    "Sibling"
  ),
  (
    1,
    "Benefator 4",
    "0888888885",
    "P.O. Box 4",
    2,
    "No",
    "Other"
  );

INSERT INTO
  memberLoanLiability (memberLoanId, description, value)
VALUES
  (1, "Liability 1", 100000),
  (1, "Liability 2", 50000);

INSERT INTO
  memberLoanSecurity (memberLoanId, description, value, serialNumber)
VALUES
  (1, "Security 1", 50000, "123456"),
  (1, "Security 2", 50000, "456789");

INSERT INTO
  memberLoanWitness (memberLoanId, witnessName, telephone, date)
VALUES
  (
    1,
    "Witness 1",
    "09928388727",
    "2025-09-15 14:21:52"
  ),
  (
    1,
    "Witness 2",
    "08858574646",
    "2025-09-15 14:21:52"
  );

INSERT INTO
  memberLoanApproval (
    memberLoanId,
    loanNumber,
    loanStatus,
    amountRecommended,
    partialApprovalReason,
    approvedBy
  )
SELECT
  id AS memberLoanId,
  loanNumber,
  "PARTIAL-APPROVAL",
  10000,
  "Shortage of funds",
  "admin"
FROM
  memberLoan
WHERE
  id = 1;

INSERT INTO
  memberLoanVerification (
    memberLoanApprovalId,
    loanNumber,
    verified,
    verifiedBy
  )
SELECT
  id,
  loanNumber,
  "Yes",
  "admin"
FROM
  memberLoanApproval
WHERE
  id = 1;

INSERT INTO
  memberLoanPayment (loanNumber, dueDate, description, amountPaid)
SELECT
  loanNumber,
  dueDate,
  'Repayment Month 1',
  1963.33
FROM
  memberLoanPaymentSchedule
WHERE
  id = 1;

INSERT INTO
  memberLoanPayment (loanNumber, dueDate, description, amountPaid)
SELECT
  loanNumber,
  dueDate,
  'Repayment Month 2',
  1273
FROM
  memberLoanPaymentSchedule
WHERE
  id = 2;

INSERT INTO
  memberLoanPayment (loanNumber, dueDate, description, amountPaid)
SELECT
  loanNumber,
  dueDate,
  'Repayment Month 3',
  5000
FROM
  memberLoanPaymentSchedule
WHERE
  id = 3;