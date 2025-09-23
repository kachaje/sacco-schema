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
          FLOOR(
            (
              d.amount + COALESCE(
                (
                  SELECT
                    SUM(overflowAmount)
                  FROM
                    memberContributionSchedule
                  WHERE
                    memberContributionId = NEW.memberContributionId
                    AND COALESCE(overflowAmount, 0) > 0
                  GROUP BY
                    memberContributionId
                ),
                0
              )
            ) / c.monthlyContribution
          )
        FROM
          memberContribution c
          LEFT OUTER JOIN memberContributionDeposit d ON d.memberContributionId = c.id
        WHERE
          d.contributionCategory = 'Regular Deposit'
          AND d.id = NEW.id
      )
  ),
  schedule AS (
    SELECT
      COALESCE(
        (
          SELECT
            dueDate
          FROM
            memberContributionSchedule
          WHERE
            DATE (dueDate) > DATE (CURRENT_TIMESTAMP)
            AND memberContributionId = NEW.memberContributionId
          ORDER BY
            dueDate DESC
          LIMIT
            1
        ),
        CURRENT_TIMESTAMP
      ) AS dueDate
  ),
  contribution AS (
    SELECT
      c.id AS memberContributionId,
      d.amount + COALESCE(
        (
          SELECT
            SUM(overflowAmount)
          FROM
            memberContributionSchedule
          WHERE
            memberContributionId = NEW.memberContributionId
            AND COALESCE(overflowAmount, 0) > 0
          GROUP BY
            memberContributionId
        ),
        0
      ) AS amount,
      c.monthlyContribution
    FROM
      memberContribution c
      LEFT OUTER JOIN memberContributionDeposit d ON d.memberContributionId = c.id
    WHERE
      d.contributionCategory = 'Regular Deposit'
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

UPDATE memberContributionSchedule
SET
  overflowAmount = 0
WHERE
  memberContributionId = NEW.memberContributionId
  AND id != (
    SELECT
      LAST_INSERT_ROWID ()
  );

END;