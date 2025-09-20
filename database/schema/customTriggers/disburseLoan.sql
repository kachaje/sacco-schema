CREATE TRIGGER IF NOT EXISTS disburseLoan AFTER INSERT ON memberLoanVerification FOR EACH ROW WHEN NEW.verified = "Yes" BEGIN
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
  memberLoanPaymentSchedule (
    memberLoanId,
    dueDate,
    principal,
    interest,
    insurance,
    processingFee,
    instalment,
    amountRecommended,
    loanNumber
  )
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
          l.repaymentPeriodInMonths
        FROM
          memberLoan l
          LEFT JOIN memberLoanApproval a ON a.memberLoanId = l.id
          LEFT JOIN memberLoanVerification v ON v.memberLoanApprovalId = a.id
        WHERE
          v.id = NEW.id
      )
  )
SELECT
  i.memberLoanId,
  DATE (
    CURRENT_TIMESTAMP,
    CONCAT ('+', x, ' month'),
    'start of month'
  ) AS dueDate,
  (i.amountRecommended - ((x -1) * i.instalment)) AS principal,
  (i.amountRecommended - ((x -1) * i.instalment)) * i.monthlyInterestRate * (1 + tax) AS interest,
  CASE
    WHEN x = 1 THEN i.amountRecommended * i.monthlyInsuranceRate
    ELSE 0
  END AS insurance,
  CASE
    WHEN x = 1 THEN i.amountRecommended * i.processingFeeRate * (1 + tax)
    ELSE 0
  END AS processingFee,
  i.instalment,
  i.amountRecommended,
  i.loanNumber
FROM
  cnt,
  (
    SELECT
      l.id AS memberLoanId,
      a.amountRecommended / l.repaymentPeriodInMonths AS instalment,
      l.repaymentPeriodInMonths,
      a.amountRecommended,
      l.monthlyInterestRate,
      l.monthlyInsuranceRate,
      l.processingFeeRate,
      l.loanNumber,
      COALESCE(
        (
          SELECT
            value
          FROM
            taxRate
          WHERE
            name = 'VAT'
        ),
        0
      ) AS tax
    FROM
      memberLoan l
      LEFT JOIN memberLoanApproval a ON a.memberLoanId = l.id
      LEFT JOIN memberLoanVerification v ON v.memberLoanApprovalId = a.id
    WHERE
      v.id = NEW.id
  ) AS i;

END;