CREATE TRIGGER IF NOT EXISTS allocateLoanPayment AFTER INSERT ON memberLoanPayment FOR EACH ROW BEGIN
INSERT INTO
  memberLoanPaymentDetail (
    memberLoanPaymentId,
    loanComponent,
    billedAmount,
    paidAmount
  )
WITH
  schedule AS (
    WITH
      vat AS (
        SELECT
          value AS tax
        FROM
          taxRate
        WHERE
          active = 1
          AND name = 'VAT'
        LIMIT
          1
      )
    SELECT
      memberLoanId,
      principal,
      processingFee,
      interest,
      instalment,
      insurance,
      processingFee * tax processingFeeTax,
      interest * tax interestTax,
      dueDate,
      amountRecommended,
      tax,
      amountPaid,
      (
        (processingFee * (1 + tax)) + (interest * (1 + tax)) + instalment + insurance
      ) AS totalDue
    FROM
      memberLoanPaymentSchedule,
      vat
    WHERE
      memberLoanId = 1
      AND amountPaid < totalDue
    ORDER BY
      dueDate ASC
    LIMIT
      1
  )
SELECT
  memberLoanPaymentId,
  'Interest',
  interest,
  interest
FROM
  schedule
UNION ALL
SELECT
  memberLoanPaymentId,
  'Instalment',
  instalment,
  instalment
FROM
  schedule
UNION ALL
SELECT
  memberLoanPaymentId,
  'Insurance',
  insurance,
  insurance
FROM
  schedule;

END;