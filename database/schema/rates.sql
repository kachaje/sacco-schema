INSERT
OR IGNORE INTO memberLoanType (
  name,
  amountLimit,
  periodLimitInMonths,
  maxInstalmentMonths,
  monthlyInterestRate
)
VALUES
  ("School fees", "Individual", 2000000, 6, 0.05),
  ("Personal", "Individual", 2000000, 12, 0.048),
  ("Emergency", "Individual", 50000, 3, 0.1),
  ("Emergency", "Group/Institution", 200000, 3, 0.1);