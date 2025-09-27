CREATE TRIGGER IF NOT EXISTS updateAccount AFTER INSERT ON accountEntry FOR EACH ROW BEGIN
UPDATE account
SET
  balance = (
    CASE
      WHEN accountType IN ('ASSET', 'EXPENSE') THEN CASE
        WHEN NEW.debitCredit = 'DEBIT' THEN COALESCE(balance, 0) + NEW.amount
        ELSE COALESCE(balance, 0) - NEW.amount
      END
      ELSE CASE
        WHEN NEW.debitCredit = 'CREDIT' THEN COALESCE(balance, 0) + NEW.amount
        ELSE COALESCE(balance, 0) - NEW.amount
      END
    END
  )
WHERE
  id = NEW.accountId;

END;