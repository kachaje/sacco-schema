CREATE TRIGGER IF NOT EXISTS addContributionDepositAccount AFTER INSERT ON memberContributionDeposit FOR EACH ROW BEGIN
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
    CONCAT ('memberContributionDeposit:', NEW.id),
    'Contribution Deposit',
    NEW.description,
    'DEBIT',
    NEW.amount
  ),
  (
    (
      SELECT
        id
      FROM
        account
      WHERE
        accountType = 'EQUITY'
    ),
    CONCAT ('memberContributionDeposit:', NEW.id),
    'Contribution Deposit',
    NEW.description,
    'CREDIT',
    NEW.amount
  );

END