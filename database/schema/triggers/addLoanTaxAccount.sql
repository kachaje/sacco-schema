CREATE TRIGGER IF NOT EXISTS addLoanTaxAccount AFTER INSERT ON memberLoanTax FOR EACH ROW BEGIN
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
        accountType = 'LIABILITY'
    ),
    CONCAT ('memberLoanTax:', NEW.id),
    'Loan Tax Reserve',
    CONCAT (NEW.taxCategory, ' Tax'),
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
        accountType = 'ASSET'
    ),
    CONCAT ('memberLoanTax:', NEW.id),
    'Loan Tax Collection',
    CONCAT (NEW.taxCategory, ' Tax'),
    'DEBIT',
    NEW.amount
  );

END