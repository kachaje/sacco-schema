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
  ),
  totalDue = (interest + insurance + processingFee + instalment + penalty)
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
  id = NEW.id AND availableCash > 0;

END