CREATE TRIGGER IF NOT EXISTS addLoanDisbursementAccount AFTER INSERT ON memberLoanDisbursement FOR EACH ROW BEGIN
INSERT INTO
  accountEntry (
    accountId,
    referenceNumber,
    name,
    description,
    debitCredit,
    amount
  )
VALUES
  (
    (
      SELECT
        id
      FROM
        account
      WHERE
        accountType = 'ASSET'
    ),
    CONCAT ('memberLoanDisbursement:', NEW.id),
    'Loan Disbursement',
    NEW.description,
    'CREDIT',
    NEW.amount
  ),
  (
    (
      SELECT
        id
      FROM
        account
      WHERE
        accountType = 'EXPENSE'
    ),
    CONCAT ('memberLoanDisbursement:', NEW.id),
    'Loan Disbursement',
    NEW.description,
    'DEBIT',
    NEW.amount
  );

END