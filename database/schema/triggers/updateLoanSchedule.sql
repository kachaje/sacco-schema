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
      LEFT OUTER JOIN memberLoanPayment p ON p.loanNumber = s.loanNumber
    WHERE
      p.id = NEW.memberLoanPaymentId AND s.id = (SELECT id FROM memberLoanPaymentSchedule WHERE loanNumber = p.loanNumber AND amountPaid = 0 ORDER BY dueDate ASC LIMIT 1)
      AND NEW.loanComponent IN ('Interest', 'Processing Fee', 'Instalment')
  );

END;