CREATE TRIGGER IF NOT EXISTS applyLoanRates AFTER INSERT ON memberLoan FOR EACH ROW BEGIN
UPDATE memberLoan
SET
  monthlyInterestRate = loanType.monthlyInterestRate,
  monthlyInsuranceRate = loanType.monthlyInsuranceRate,
  processingFeeRate = loanType.processingFeeRate,
  penaltyRate = loanType.penaltyRate
FROM
  loanType
WHERE
  loanType.name = memberLoan.loanType
  AND loanType.category = memberLoan.loanCategory;

END;