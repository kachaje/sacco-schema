INSERT
OR IGNORE INTO loanType (
  name,
  amountLimit,
  periodLimitInMonths,
  maxInstallmentMonths,
  monthlyInterestRate
)
VALUES
  ("School fees", "Individual", 2000000, 6, 6, 0.05),
  ("Personal", "Individual", 2000000, 12, 12, 0.048),
  ("Emergency", "Individual", 50000, 3, 3, 0.1),
  (
    "Emergency",
    "Group/Institution",
    200000,
    3,
    3,
    0.1
  );