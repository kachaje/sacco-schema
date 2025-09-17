CREATE TRIGGER IF NOT EXISTS disburseLoan AFTER INSERT ON memberLoanVerification WHEN NEW.verified = "Yes" BEGIN
INSERT INTO
  memberLoanDisbursement (memberLoanId, description, amount)
SELECT
  a.memberLoanId,
  CONCAT ("Disbursement for loan number ", v.loanNumber) description,
  a.amountRecommended AS amount
FROM
  memberLoanVerification v
  LEFT JOIN memberLoanApproval a ON a.id = v.memberLoanApprovalId
WHERE
  v.id = NEW.id;

INSERT INTO
  memberLoanProcessingFee (memberLoanId, amount)
SELECT
  l.id AS memberLoanId,
  l.processingFeeRate * a.amountRecommended AS amount
FROM
  memberLoan l
  LEFT JOIN memberLoanApproval a;

INSERT INTO
  memberLoanTax (memberLoanId, amount, description)
SELECT
  memberLoanId,
  amount * (
    SELECT
      value
    FROM
      taxRate
    WHERE
      name = "VAT"
    LIMIT
      1
  ),
  description
FROM
  memberLoanDisbursement;

WITH RECURSIVE
  cnt (x) AS (
    SELECT
      1
    UNION ALL
    SELECT
      x + 1
    FROM
      cnt
    LIMIT
      (
        SELECT
          repaymentPeriodInMonths
        FROM
          memberLoan
        WHERE
          id = 1
      )
  )
SELECT
  DATE (
    current_timestamp,
    CONCAT ('+', x, ' month'),
    'start of month'
  ) AS dueDate,
  (i.amountRecommended - ((x -1) * i.instalment)) AS principal,
  (i.amountRecommended - ((x -1) * i.instalment)) * i.monthlyInterestRate AS interest,
  (i.amountRecommended - ((x -1) * i.instalment)) * i.monthlyInsuranceRate AS insurance,
  CASE
    WHEN x = 1 THEN (i.amountRecommended - ((x -1) * i.instalment)) * i.processingFeeRate
    ELSE 0
  END AS processingFee
FROM
  cnt,
  (
    SELECT
      a.amountRecommended / l.repaymentPeriodInMonths AS instalment,
      l.repaymentPeriodInMonths,
      a.amountRecommended,
      l.monthlyInterestRate,
      l.monthlyInsuranceRate,
      l.processingFeeRate
    FROM
      memberLoan l
      LEFT JOIN memberLoanApproval a ON a.memberLoanId = l.id
  ) AS i;

END;