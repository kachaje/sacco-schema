CREATE TRIGGER IF NOT EXISTS addContributionSchedule AFTER INSERT ON memberContributionDeposit FOR EACH ROW BEGIN
INSERT INTO
  memberContributionSchedule (
    memberContributionId,
    dueDate,
    expectedAmount,
    paidAmount,
    overflowAmount
  )
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
      (
        SELECT
          FLOOR(d.amount / c.monthlyContribution)
        FROM
          memberContribution c
          LEFT OUTER JOIN memberContributionDeposit d ON d.memberContributionId = c.id
        WHERE
          d.contributionCategory = "Regular Deposit"
          AND d.id = NEW.id
      )
  ),
  schedule AS (
    SELECT
      SUM(paidAmount - overflowAmount) AS balance,
      dueDate
    FROM
      memberContributionSchedule
    WHERE
      paidAmount < expectedAmount
      AND DATE (dueDate) > DATE (CURRENT_TIMESTAMP)
      AND memberContributionId = NEW.id
    ORDER BY
      dueDate ASC
    LIMIT
      1
  ),
  contribution AS (
    SELECT
      c.id AS memberContributionId,
      d.amount,
      c.monthlyContribution
    FROM
      memberContribution c
      LEFT OUTER JOIN memberContributionDeposit d ON d.memberContributionId = c.id
    WHERE
      d.contributionCategory = "Regular Deposit"
      AND d.id = NEW.id
  )
SELECT
  memberContributionId,
  CAST(
    DATE (
      COALESCE(dueDate, CURRENT_TIMESTAMP),
      CONCAT ('+', x, ' month'),
      'start of month'
    ) AS TEXT
  ) AS dueDate,
  monthlyContribution AS expectedAmount,
  CAST(
    CASE
      WHEN x * monthlyContribution <= amount THEN monthlyContribution
      ELSE 0
    END AS REAL
  ) AS paidAmount,
  CAST(
    CASE
      WHEN (x + 1) * monthlyContribution > amount THEN amount - (x * monthlyContribution)
      ELSE 0
    END AS REAL
  ) AS overflowAmount
FROM
  cnt,
  schedule,
  contribution;

END;