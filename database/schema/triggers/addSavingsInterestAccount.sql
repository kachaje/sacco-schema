CREATE TRIGGER IF NOT EXISTS addSavingInterestAccount AFTER INSERT ON memberSavingInterest FOR EACH ROW BEGIN
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
    CONCAT ('memberSavingInterest:', NEW.id),
    'Savings Interest',
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
    CONCAT ('memberSavingInterest:', NEW.id),
    'Savings Interest',
    NEW.description,
    'DEBIT',
    NEW.amount
  );

END