CREATE TRIGGER IF NOT EXISTS handleSavingTransaction AFTER INSERT ON memberSavingTransaction FOR EACH ROW BEGIN
UPDATE memberSaving
SET
  balance = COALESCE(balance, 0) + COALESCE(NEW.deposit, 0) - COALESCE(NEW.withdrawal, 0)
WHERE
  id = NEW.memberSavingId;

UPDATE memberSavingTransaction
SET
  balance = COALESCE(
    (
      SELECT
        balance
      FROM
        memberSaving
      WHERE
        id = NEW.memberSavingId
    ),
    0
  ),
  savingsTypeName = (
    SELECT
      savingsTypeName
    FROM
      memberSaving
    WHERE
      id = NEW.memberSavingId
  )
WHERE
  id = NEW.id;

END;