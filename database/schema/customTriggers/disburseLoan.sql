CREATE TRIGGER IF NOT EXISTS disburseLoan AFTER INSERT ON memberLoanVerification WHEN NEW.verified = "Yes" BEGIN
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
  memberLoanProcessingFee (memberLoanId, amount)
SELECT
  l.id AS memberLoanId,
  l.processingFeeRate * a.amountRecommended AS amount
FROM
  memberLoan l
  LEFT JOIN memberLoanApproval a;

INSERT INTO
  memberLoanTax (amount, description)
SELECT
  amount * (
    SELECT
      value
    FROM
      taxRate
    WHERE
      name = "VAT"
    LIMIT
      1
  ),
  description
FROM
  memberLoanDisbursement;

END;