CREATE TRIGGER IF NOT EXISTS updateLoanSchedule AFTER INSERT ON memberLoanPaymentDetail FOR EACH ROW BEGIN
UPDATE memberLoanPaymentSchedule
SET
  amountPaid = amountPaid + NEW.amount
WHERE
  id IN (
    SELECT
      s.id
    FROM
      memberLoanPaymentSchedule s
      LEFT OUTER JOIN memberLoanPayment p ON p.memberLoanPaymentScheduleId = s.id
    WHERE
      p.id = NEW.memberLoanPaymentId
  );

END;