CREATE TRIGGER IF NOT EXISTS handleLoanPaymentDetail AFTER INSERT ON memberLoanPaymentDetail FOR EACH ROW BEGIN
INSERT INTO
  memberLoanSettlement (memberId, amountReserved)
SELECT
  l.id,
  NEW.amount
FROM
  memberLoan l
  LEFT OUTER JOIN memberLoanPaymentSchedule s ON s.memberLoanId = l.id
  LEFT OUTER JOIN memberLoanPayment p ON p.memberLoanPaymentScheduleId = s.id
WHERE
  p.id = NEW.memberLoanPaymentId AND NEW.loanComponent = 'Settlement Overflow';

END;