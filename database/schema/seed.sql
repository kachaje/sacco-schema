INSERT
OR IGNORE INTO userRole (name)
VALUES
  ("Default"),
  ("Member"),
  ("Admin"),
  ("Cashier"),
  ("Accountant"),
  ("Loans Officer"),
  ("Manager");

INSERT
OR IGNORE INTO user (username, password, name, userRole)
VALUES
  (
    "admin",
    "$2a$10$Xo4x3KiCkB3xGKvaCI4Hn.Be95DEiaIT3lbvHx/kOmyx7IqGY6ILK",
    "Default User",
    "Admin"
  );

WITH RECURSIVE
  cnt (x) AS (
    SELECT
      1
    UNION ALL
    SELECT
      x + 1
    FROM
      cnt
    LIMIT
      999999
  ) INSERT
  OR IGNORE INTO memberIdsCache (idNumber)
SELECT
  CONCAT ('KSM', SUBSTR ('000000' || x, -6)) AS id
FROM
  cnt
ORDER BY
  RANDOM ();

WITH RECURSIVE
  cnt (x) AS (
    SELECT
      1
    UNION ALL
    SELECT
      x + 1
    FROM
      cnt
    LIMIT
      999999
  ) INSERT
  OR IGNORE INTO contributionNumberIdsCache (idNumber)
SELECT
  CONCAT ('KSH', SUBSTR ('000000' || x, -6)) AS id
FROM
  cnt
ORDER BY
  RANDOM ();

WITH RECURSIVE
  cnt (x) AS (
    SELECT
      1
    UNION ALL
    SELECT
      x + 1
    FROM
      cnt
    LIMIT
      999999
  ) INSERT
  OR IGNORE INTO memberSavingIdsCache (idNumber)
SELECT
  CONCAT ('KSS', SUBSTR ('000000' || x, -6)) AS id
FROM
  cnt
ORDER BY
  RANDOM ();

WITH RECURSIVE
  cnt (x) AS (
    SELECT
      1
    UNION ALL
    SELECT
      x + 1
    FROM
      cnt
    LIMIT
      999999
  ) INSERT
  OR IGNORE INTO loanNumberIdsCache (loanNumber)
SELECT
  CONCAT ('KLN', SUBSTR ('000000' || x, -6)) AS id
FROM
  cnt
ORDER BY
  RANDOM ();

INSERT
OR IGNORE INTO taxRate (name, value)
VALUES
  ("VAT", 0.0);

INSERT
OR IGNORE INTO savingsType (
  savingsTypeName,
  minimumAmount,
  withdrawPattern,
  minWithdrawMonths,
  maxWithdrawMonths,
  interestRate
)
VALUES
  (
    'Fixed Deposit',
    50000,
    '3 to 12 months',
    3,
    12,
    0.1
  ),
  ('Ordinary Deposit', 2000, 'Anytime', 0, 12, 0.07),
  (
    '30 day Call Deposit',
    25000,
    '30 days',
    1,
    1,
    0.07
  );

INSERT
OR IGNORE INTO account (name, accountType, increasedBy, decreasedBy)
VALUES
  ('Assets', 'ASSET', 'DEBIT', 'CREDIT'),
  ('Liabilities', 'LIABILITY', 'CREDIT', 'DEBIT'),
  ('Expenses', 'EXPENSE', 'DEBIT', 'CREDIT'),
  ('Equity', 'EQUITY', 'CREDIT', 'DEBIT'),
  ('Revenue', 'REVENUE', 'CREDIT', 'DEBIT');