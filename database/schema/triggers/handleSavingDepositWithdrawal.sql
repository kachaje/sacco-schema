CREATE TRIGGER IF NOT EXISTS handleSavingDeposit AFTER INSERT ON memberSavingDeposit FOR EACH ROW BEGIN
UPDATE memberSaving
SET
  balance = COALESCE(balance, 0) + COALESCE(NEW.amount, 0)
WHERE
  id = NEW.memberSavingId;

UPDATE memberSavingDeposit
SET
  balance = (
    SELECT
      balance
    FROM
      memberSaving
    WHERE
      id = NEW.memberSavingId
  )
WHERE
  id = NEW.id;

END;

CREATE TRIGGER IF NOT EXISTS handleSavingWithdrawal AFTER INSERT ON memberSavingWithdrawal FOR EACH ROW BEGIN
UPDATE memberSaving
SET
  balance = COALESCE(balance, 0) - COALESCE(NEW.amount, 0)
WHERE
  id = NEW.memberSavingId;

UPDATE memberSavingWithdrawal
SET
  balance = (
    SELECT
      balance
    FROM
      memberSaving
    WHERE
      id = NEW.memberSavingId
  )
WHERE
  id = NEW.id;

END;