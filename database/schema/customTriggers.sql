CREATE TRIGGER IF NOT EXISTS applyLoanRates AFTER INSERT ON memberLoan FOR EACH ROW BEGIN
UPDATE memberLoan
SET
  monthlyInterestRate = (
    SELECT
      monthlyInterestRate
    FROM
      memberLoanType lt
    WHERE
      lt.name = NEW.loanType
      AND category = NEW.loanCategory
  ),
  monthlyInsuranceRate = (
    SELECT
      monthlyInsuranceRate
    FROM
      memberLoanType lt
    WHERE
      lt.name = NEW.loanType
      AND category = NEW.loanCategory
  ),
  processingFeeRate = (
    SELECT
      processingFeeRate
    FROM
      memberLoanType lt
    WHERE
      lt.name = NEW.loanType
      AND category = NEW.loanCategory
  ),
  penaltyRate = (
    SELECT
      penaltyRate
    FROM
      memberLoanType lt
    WHERE
      lt.name = NEW.loanType
      AND category = NEW.loanCategory
  ),
  monthlyInstalments = NEW.loanAmount / NEW.repaymentPeriodInMonths
WHERE
  memberLoan.id = NEW.id;

END;