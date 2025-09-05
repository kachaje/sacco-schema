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
    "default",
    "$2a$10$Xo4x3KiCkB3xGKvaCI4Hn.Be95DEiaIT3lbvHx/kOmyx7IqGY6ILK",
    "Default User",
    "Default"
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
  OR IGNORE INTO memberSharesIdsCache (idNumber)
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
  OR IGNORE INTO memberSavingsIdsCache (idNumber)
SELECT
  CONCAT ('KSS', SUBSTR ('000000' || x, -6)) AS id
FROM
  cnt
ORDER BY
  RANDOM ();