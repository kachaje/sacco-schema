CREATE TRIGGER IF NOT EXISTS allocateLoanPayment AFTER INSERT ON memberLoanPayment FOR EACH ROW BEGIN
INSERT INTO
  memberLoanPaymentDetail (memberLoanPaymentId, loanComponent, amount)
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
  NEW.id,
  'Interest',
  interest
FROM
  schedule
UNION ALL
SELECT
  NEW.id,
  'Instalment',
  instalment
FROM
  schedule
UNION ALL
SELECT
  NEW.id,
  'Insurance',
  insurance
FROM
  schedule
WHERE
  insurance > 0
UNION ALL
SELECT
  NEW.id,
  'Processing Fee',
  processingFee
FROM
  schedule
WHERE
  processingFee > 0
UNION ALL
SELECT
  NEW.id,
  'Settlement Overflow',
  (NEW.amountPaid - totalDue)
FROM
  schedule
WHERE
  (NEW.amountPaid - totalDue) > 0;

END;