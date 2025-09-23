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
      20
    OFFSET
      1000
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
  memberContribution (
    memberId,
    memberIdNumber,
    monthlyContribution,
    nonRedeemableAmount
  )
WITH RECURSIVE
  cnt (i, rand) AS (
    SELECT
      1,
      CAST(ABS(MOD(RANDOM (), 10)) AS INTEGER)
    UNION ALL
    SELECT
      i + 1,
      CAST(ABS(MOD(RANDOM (), 10)) AS INTEGER)
    FROM
      cnt
    LIMIT
      20
    OFFSET
      1000
  )
SELECT
  i,
  (
    SELECT
      memberIdNumber
    FROM
      member
    WHERE
      id = i
  ),
  (rand + 1) * 5000,
  20000
FROM
  cnt;

INSERT INTO
  memberContributionDeposit (
    memberContributionId,
    contributionCategory,
    amount
  )
WITH RECURSIVE
  cnt (i, rand) AS (
    SELECT
      1,
      CAST(ABS(MOD(RANDOM (), 10)) AS INTEGER)
    UNION ALL
    SELECT
      i + 1,
      CAST(ABS(MOD(RANDOM (), 10)) AS INTEGER)
    FROM
      cnt
    LIMIT
      20
    OFFSET
      1000
  )
SELECT
  (
    SELECT
      id
    FROM
      memberContribution
    WHERE
      memberId = i
  ),
  'Non-Redeemable Deposit',
  20000
FROM
  cnt;

INSERT INTO
  memberContributionDeposit (
    memberContributionId,
    contributionCategory,
    amount
  )
WITH RECURSIVE
  cnt (i, rand) AS (
    SELECT
      1,
      CAST(ABS(MOD(RANDOM (), 10)) AS INTEGER)
    UNION ALL
    SELECT
      i + 1,
      CAST(ABS(MOD(RANDOM (), 10)) AS INTEGER)
    FROM
      cnt
    LIMIT
      20
    OFFSET
      1000
  )
SELECT
  id,
  'Regular Deposit',
  monthlyContribution * 12.5
FROM
  cnt,
  memberContribution
WHERE
  memberId = i;