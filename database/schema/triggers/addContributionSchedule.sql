CREATE TRIGGER IF NOT EXISTS addContributionSchedule AFTER INSERT ON memberContribution FOR EACH ROW BEGIN
INSERT INTO
  memberContributionSchedule (memberContributionId, dueDate, expectedAmount)
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
      12
  )
SELECT
  NEW.id,
  DATE (
    CURRENT_TIMESTAMP,
    CONCAT ('+', x, ' month'),
    'start of month'
  ) AS dueDate,
  NEW.monthlyContribution
FROM
  cnt;

END;