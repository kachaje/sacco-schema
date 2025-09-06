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
    nationalId,
    utilityBillType,
    utilityBillNumber,
    fileNumber,
    oldFileNumber,
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
    "",
    "",
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
  memberNominee (memberId, name, phoneNumber, address)
VALUES
  (1, "John Banda", "0888888888", "Same as member");

INSERT INTO
  memberBeneficiary (memberId, name, percentage, contact)
VALUES
  (1, "Benefator 1", 10, "P.O. Box 1"),
  (1, "Benefator 2", 8, "P.O. Box 2");

INSERT INTO
  memberLoan (
    memberId,
    loanAmount,
    repaymentPeriodInMonths,
    loanPurpose,
    loanType
  )
VALUES
  (1, 200000, 12, "School fees", "PERSONAL");

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
    employeesWages,
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
    employeesWages,
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
  memberLoanApproval (
    memberLoanId,
    loanStatus,
    amountRecommended,
    approvedBy,
    approvalDate,
    amountApproved,
    verifiedBy,
    dateVerified
  )
VALUES
  (
    1,
    "APPROVED",
    200000,
    "me",
    "2025-08-30",
    200000,
    "me",
    "2025-08-30"
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
