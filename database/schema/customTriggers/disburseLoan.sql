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

END;