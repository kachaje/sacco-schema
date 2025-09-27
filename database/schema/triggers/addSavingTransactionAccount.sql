CREATE TRIGGER IF NOT EXISTS addSavingTransactionAccount AFTER INSERT ON memberSavingTransaction FOR EACH ROW BEGIN
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
    CONCAT ('memberSavingTransaction:', NEW.id),
    CASE
      WHEN COALESCE(NEW.deposit, 0) > 0 THEN 'Savings Deposit'
      ELSE 'Savings Withdrawal'
    END,
    NEW.description,
    CASE
      WHEN COALESCE(NEW.deposit, 0) > 0 THEN 'DEBIT'
      ELSE 'CREDIT'
    END,
    COALESCE(NEW.deposit, NEW.withdrawal)
  ),
  (
    (
      SELECT
        id
      FROM
        account
      WHERE
        accountType = CASE
          WHEN COALESCE(NEW.deposit, 0) > 0 THEN 'REVENUE'
          ELSE 'EXPENSE'
        END
    ),
    CONCAT ('memberSavingTransaction:', NEW.id),
    CASE
      WHEN COALESCE(NEW.deposit, 0) > 0 THEN 'Savings Deposit'
      ELSE 'Savings Withdrawal'
    END,
    NEW.description,
    CASE
      WHEN COALESCE(NEW.deposit, 0) > 0 THEN 'CREDIT'
      ELSE 'DEBIT'
    END,
    COALESCE(NEW.deposit, NEW.withdrawal)
  );

END