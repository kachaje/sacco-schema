PRAGMA journal_mode=WAL;

INSERT INTO
  member (
    firstName,
    lastName,
    gender,
    phoneNumber,
    title,
    maritalStatus,
    dateOfBirth,
    nationalIdentifier,
    utilityBillType,
    utilityBillNumber
  )
VALUES
  (
    "Mary",
    "Banda",
    "Female",
    "0999888777",
    "Miss",
    "Single",
    "1999-09-01",
    "KJFFJ58584",
    "ESCOM",
    "949488473"
  );

CREATE TEMP TABLE IF NOT EXISTS memberIdVar (value TEXT);

INSERT INTO
  memberIdVar
SELECT
  last_insert_rowid ();

SELECT
  value
FROM
  memberIdVar;

