CREATE TRIGGER IF NOT EXISTS addContributionDividendAccount AFTER INSERT ON memberContributionDividend FOR EACH ROW BEGIN
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
    CONCAT ('memberContributionDividend:', NEW.id),
    'Contribution Dividend',
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
        accountType = 'EQUITY'
    ),
    CONCAT ('memberContributionDividend:', NEW.id),
    'Contribution Dividend',
    NEW.description,
    'DEBIT',
    NEW.amount
  );

END