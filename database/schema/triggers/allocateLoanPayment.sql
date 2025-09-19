CREATE TRIGGER IF NOT EXISTS allocateLoanPayment AFTER INSERT ON memberLoanPayment FOR EACH ROW BEGIN
UPDATE memberLoanPayment
SET
  availableCash = (
    amountPaid + COALESCE(
      (
        SELECT
          SUM(amountReserved) - SUM(amountClaimed)
        FROM
          memberLoanSettlement
        WHERE
          loanNumber = NEW.loanNumber
        GROUP BY
          loanNumber
      ),
      0
    )
  )
WHERE
  id = NEW.id;

UPDATE memberLoanPayment
SET
  totalDue = (
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
          (
            (s.processingFee * (1 + tax)) + (s.interest * (1 + tax)) + s.instalment + s.insurance
          ) AS totalDue
        FROM
          memberLoanPaymentSchedule s,
          vat
          LEFT OUTER JOIN memberLoanPayment p ON p.loanNumber = s.loanNumber
        WHERE
          s.dueDate = NEW.dueDate
          AND s.loanNumber = NEW.loanNumber
          AND p.availableCash >= totalDue
      )
    SELECT
      totalDue
    FROM
      schedule
  )
WHERE
  id = NEW.id;

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
      s.dueDate,
      amountRecommended,
      tax,
      s.amountPaid,
      (
        (processingFee * (1 + tax)) + (interest * (1 + tax)) + instalment + insurance
      ) AS totalDue
    FROM
      memberLoanPaymentSchedule s,
      vat
      LEFT OUTER JOIN memberLoanPayment p ON p.loanNumber = s.loanNumber
    WHERE
      s.dueDate = NEW.dueDate
      AND s.loanNumber = NEW.loanNumber
      AND p.availableCash >= totalDue
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

INSERT INTO
  memberLoanTax (
    memberLoanPaymentId,
    description,
    amount,
    taxCategory
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
      processingFee * tax processingFeeTax,
      interest * tax interestTax,
      (
        (processingFee * (1 + tax)) + (interest * (1 + tax)) + instalment + insurance
      ) AS totalDue
    FROM
      memberLoanPaymentSchedule s,
      vat
      LEFT OUTER JOIN memberLoanPayment p ON p.loanNumber = s.loanNumber
    WHERE
      s.dueDate = NEW.dueDate
      AND s.loanNumber = NEW.loanNumber
      AND p.availableCash >= totalDue
  )
SELECT
  NEW.id,
  'Tax on Interest',
  interestTax,
  'Interest'
FROM
  schedule
WHERE
  interestTax > 0
UNION ALL
SELECT
  NEW.id,
  'Tax on Processing Fee',
  processingFeeTax,
  'Processing Fee'
FROM
  schedule
WHERE
  processingFeeTax > 0;

UPDATE memberLoanSettlement
SET
  amountClaimed = (
    SELECT
      totalDue
    FROM
      memberLoanPayment
    WHERE
      id = NEW.id
  )
WHERE
  id = (
    SELECT
      id
    FROM
      memberLoanSettlement
    WHERE
      loanNumber = NEW.loanNumber
      AND COALESCE(amountClaimed, 0) = 0
    ORDER BY
      id
    LIMIT
      1
  );

INSERT INTO
  memberLoanSettlement (loanNumber, amountReserved)
WITH
  settlement AS (
    SELECT
      COALESCE(
        (
          SELECT
            SUM(amountReserved) - SUM(amountClaimed)
          FROM
            memberLoanSettlement
          WHERE
            loanNumber = NEW.loanNumber
          GROUP BY
            loanNumber
        ),
        0
      ) - NEW.totalDue AS balance
  )
SELECT
  NEW.loanNumber,
  balance
FROM
  settlement
WHERE
  balance > 0;

END;