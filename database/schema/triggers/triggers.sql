---- START addMemberIdNumber TRIGGER ----
CREATE TRIGGER IF NOT EXISTS addMemberIdNumber AFTER INSERT ON member FOR EACH ROW BEGIN
UPDATE memberIdsCache
SET
  claimed = 1,
  memberId = NEW.id
WHERE
  id = (
    SELECT
      id
    FROM
      memberIdsCache
    WHERE
      claimed = 0
    ORDER BY
      id
    LIMIT
      1
  );

UPDATE member
SET
  memberIdNumber = (
    SELECT
      idNumber
    FROM
      memberIdsCache
    WHERE
      memberId = NEW.id
  )
WHERE
  id = NEW.id;

---- END addMemberIdNumber TRIGGER ----
END;

---- START addMemberIdNumber TRIGGER ----
CREATE TRIGGER IF NOT EXISTS addMemberIdNumber AFTER INSERT ON member FOR EACH ROW BEGIN
UPDATE memberSavingIdsCache
SET
  claimed = 1,
  memberSavingId = NEW.id
WHERE
  id = (
    SELECT
      id
    FROM
      memberSavingIdsCache
    WHERE
      claimed = 0
    ORDER BY
      id
    LIMIT
      1
  );

UPDATE member
SET
  memberSavingIdNumber = (
    SELECT
      idNumber
    FROM
      memberSavingIdsCache
    WHERE
      memberSavingId = NEW.id
  )
WHERE
  id = NEW.id;

---- END addMemberIdNumber TRIGGER ----
END;

---- START addLoanNumber TRIGGER ----
CREATE TRIGGER IF NOT EXISTS addLoanNumber AFTER INSERT ON memberLoan FOR EACH ROW BEGIN
UPDATE loanNumberIdsCache
SET
  claimed = 1,
  memberLoanId = NEW.id
WHERE
  id = (
    SELECT
      id
    FROM
      loanNumberIdsCache
    WHERE
      claimed = 0
    ORDER BY
      id
    LIMIT
      1
  );

UPDATE memberLoan
SET
  loanNumber = (
    SELECT
      loanNumber
    FROM
      loanNumberIdsCache
    WHERE
      memberLoanId = NEW.id
  )
WHERE
  id = NEW.id;

---- END addLoanNumber TRIGGER ----
END;