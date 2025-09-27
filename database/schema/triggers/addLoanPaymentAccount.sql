CREATE TRIGGER IF NOT EXISTS addLoanPaymentAccount AFTER INSERT ON memberLoanPayment FOR EACH ROW BEGIN
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
        accountType = 'REVENUE'
    ),
    CONCAT ('memberLoanPayment:', NEW.id),
    'Loan Payment',
    NEW.description,
    'CREDIT',
    NEW.amountPaid
  ),
  (
    (
      SELECT
        id
      FROM
        account
      WHERE
        accountType = 'ASSET'
    ),
    CONCAT ('memberLoanPayment:', NEW.id),
    'Loan Payment',
    NEW.description,
    'DEBIT',
    NEW.amountPaid
  );

END