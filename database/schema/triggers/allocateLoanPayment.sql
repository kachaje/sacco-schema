CREATE TRIGGER IF NOT EXISTS allocateLoanPayment AFTER INSERT ON memberLoanPayment FOR EACH ROW BEGIN
UPDATE memberLoanPaymentSchedule
SET
  penalty = COALESCE(
    (
      SELECT
        COALESCE(penaltyRate, 0) * (
          1 + COALESCE(
            (
              SELECT
                value
              FROM
                taxRate
              WHERE
                name = 'VAT'
                AND active = 1
            ),
            0
          )
        )
      FROM
        memberLoan
      WHERE
        loanNumber = NEW.loanNumber
    ),
    0
  ) * instalment
WHERE
  JULIANDAY (NEW.date) - JULIANDAY (dueDate) > 30
  AND id = NEW.memberLoanPaymentScheduleId;

UPDATE memberLoanPayment
SET
  penalty = COALESCE(
    (
      SELECT
        penalty
      FROM
        memberLoanPaymentSchedule
      WHERE
        id = NEW.memberLoanPaymentScheduleId
    ),
    0
  )
WHERE
  id = NEW.id;

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
  ),
  totalDue = (
    interest + insurance + processingFee + instalment + penalty
  )
WHERE
  id = NEW.id;

UPDATE memberLoanPaymentSchedule
SET
  amountPaid = amountPaid + COALESCE(
    (
      SELECT
        totalDue
      FROM
        memberLoanPayment
      WHERE
        id = NEW.id
        AND availableCash >= totalDue
    ),
    0
  )
WHERE
  id = NEW.memberLoanPaymentScheduleId;

INSERT INTO
  memberLoanSettlement (loanNumber, amountReserved, amountClaimed)
SELECT
  loanNumber,
  amountPaid,
  totalDue
FROM
  memberLoanPayment
WHERE
  id = NEW.id
  AND availableCash > 0;

INSERT INTO
  memberLoanTax (
    memberLoanPaymentId,
    description,
    amount,
    taxCategory
  )
WITH
  vat AS (
    SELECT
      value AS tax
    FROM
      taxRate
    WHERE
      name = 'VAT'
      AND active = 1
  )
SELECT
  NEW.id,
  CONCAT ('Tax: ', NEW.description),
  COALESCE(vat.tax, 0) * NEW.interest,
  'Interest'
FROM
  vat
WHERE
  vat.tax > 0
UNION ALL
SELECT
  NEW.id,
  CONCAT ('Tax: ', NEW.description),
  COALESCE(vat.tax, 0) * NEW.processingFee,
  'Processing Fee'
FROM
  vat
WHERE
  vat.tax > 0
UNION ALL
SELECT
  NEW.id,
  CONCAT ('Tax: ', NEW.description),
  COALESCE(vat.tax, 0) * NEW.penalty,
  'Penalty'
FROM
  vat
WHERE
  vat.tax > 0;

END