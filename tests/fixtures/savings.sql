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