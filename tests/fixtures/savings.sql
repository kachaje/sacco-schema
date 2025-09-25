INSERT INTO
  member (
    id,
    firstName,
    lastName,
    gender,
    title,
    maritalStatus,
    dateOfBirth,
    nationalIdentifier,
    utilityBillType,
    utilityBillNumber,
    phoneNumber
  )
WITH RECURSIVE
  cnt (i) AS (
    SELECT
      1
    UNION ALL
    SELECT
      i + 1
    FROM
      cnt
    LIMIT
      10
    OFFSET
      500
  )
SELECT
  i,
  CONCAT ("FirstName", i),
  CONCAT ("LastName", i),
  CASE
    WHEN MOD(i, 2) = 0 THEN "Female"
    ELSE "Male"
  END,
  CASE
    WHEN MOD(i, 2) = 0 THEN "Mrs"
    ELSE "Mr"
  END,
  'Married',
  DATE ('1979-09-01', CONCAT ('+', i, ' day')),
  CONCAT ("NATID", i),
  'ESCOM',
  CONCAT ('123456', i),
  CONCAT ('099987', i)
FROM
  cnt;

INSERT INTO
  memberSaving (
    memberId,
    savingsTypeId,
    savingsTypeName,
    withdrawPattern,
    minimumAmount
  )
WITH RECURSIVE
  cnt (i) AS (
    SELECT
      1
    UNION ALL
    SELECT
      i + 1
    FROM
      cnt
    LIMIT
      10
    OFFSET
      500
  )
SELECT
  i,
  s.id,
  s.savingsTypeName,
  s.withdrawPattern,
  s.minimumAmount
FROM
  cnt,
  savingsType s
WHERE
  s.id = (
    SELECT
      id
    FROM
      savingsType
    WHERE
      savingsTypeName = CASE
        WHEN MOD(i, 2) = 0 THEN 'Fixed Deposit'
        WHEN MOD(i, 3) = 0 THEN '30 day Call Deposit'
        ELSE 'Ordinary Deposit'
      END
  );

INSERT INTO
  memberSavingTransaction (memberSavingId, description, deposit, date)
WITH RECURSIVE
  cnt (i) AS (
    SELECT
      1
    UNION ALL
    SELECT
      i + 1
    FROM
      cnt
    LIMIT
      10
    OFFSET
      500
  )
SELECT
  s.id,
  s.savingsTypeName,
  s.minimumAmount * CASE
    WHEN s.savingsTypeName = 'Fixed Deposit' THEN 3
    WHEN s.savingsTypeName = '30 day Call Deposit' THEN 4
    ELSE 25
  END,
  '2025-09-01'
FROM
  cnt,
  memberSaving s
WHERE
  s.memberId = i;

INSERT INTO
  memberSavingTransaction (memberSavingId, description, withdrawal, date)
WITH RECURSIVE
  cnt (i) AS (
    SELECT
      1
    UNION ALL
    SELECT
      i + 1
    FROM
      cnt
    LIMIT
      10
    OFFSET
      500
  )
SELECT
  s.id,
  CONCAT (s.savingsTypeName, ' withdrawal'),
  CASE
    WHEN s.savingsTypeName = 'Fixed Deposit' THEN 2 * s.minimumAmount
    WHEN s.savingsTypeName = '30 day Call Deposit' THEN 4 * s.minimumAmount
    ELSE s.minimumAmount
  END,
  CASE
    WHEN s.savingsTypeName = 'Fixed Deposit' THEN DATE ('2025-09-01', CONCAT ('+', 3, ' month'))
    WHEN s.savingsTypeName = '30 day Call Deposit' THEN DATE ('2025-09-01', CONCAT ('+', 1, ' month'))
    ELSE DATE ('2025-09-01', CONCAT ('+', 1, ' month'))
  END
FROM
  cnt,
  memberSaving s
WHERE
  s.memberId = i;

INSERT INTO
  memberSavingTransaction (memberSavingId, description, withdrawal, date)
WITH RECURSIVE
  cnt (i) AS (
    SELECT
      1
    UNION ALL
    SELECT
      i + 1
    FROM
      cnt
    LIMIT
      10
    OFFSET
      500
  ),
  pos (x) AS (
    SELECT
      1
    UNION ALL
    SELECT
      x + 1
    FROM
      pos
    LIMIT
      10
  )
SELECT
  s.id,
  CONCAT (s.savingsTypeName, ' withdrawal'),
  s.minimumAmount * 2,
  DATE ('2025-09-01', CONCAT ('+', x, ' month'))
FROM
  cnt,
  pos,
  memberSaving s
WHERE
  s.memberId = i
  AND s.savingsTypeName = 'Ordinary Deposit';