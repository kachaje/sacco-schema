CREATE TRIGGER IF NOT EXISTS addMemberLoanInvoiceDetail AFTER INSERT ON memberLoanInvoice FOR EACH ROW BEGIN
INSERT INTO
  memberLoanInvoiceDetail (memberLoanInvoiceId, loanComponent, billedAmount)
VALUES
  (NEW.id, "Interest", NEW.interest),
  (NEW.id, "Instalment", NEW.instalment),
  (NEW.id, "Insurance", NEW.insurance);

INSERT INTO
  memberLoanInvoiceDetail (memberLoanInvoiceId, loanComponent, billedAmount)
SELECT
  NEW.id,
  "Processing Fee",
  NEW.processingFee
WHERE
  NEW.processingFee > 0;

END;

CREATE TRIGGER IF NOT EXISTS addMemberLoanInvoice AFTER INSERT ON memberLoanPaymentSchedule FOR EACH ROW BEGIN
INSERT INTO
  memberLoanInvoice (
    memberLoanPaymentScheduleId,
    loanNumber,
    description,
    totalDue,
    interest,
    insurance,
    processingFee,
    instalment
  )
SELECT
  NEW.id,
  loanNumber,
  CONCAT (
    "Loan payment on ",
    NEW.dueDate,
    " for ",
    loanNumber
  ) AS description,
  (
    NEW.interest + NEW.insurance + NEW.processingFee + NEW.instalment
  ) AS totalDue,
  NEW.interest,
  NEW.insurance,
  NEW.processingFee,
  NEW.instalment
FROM
  memberLoan
WHERE
  id = NEW.memberLoanId;

END;