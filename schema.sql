CREATE TABLE
  IF NOT EXISTS notification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    message TEXT NOT NULL,
    date TEXT DEFAULT 'CURRENT_TIMESTAMP',
    messageRead TEXT CHECK (messageRead IN ('Yes', 'No')) DEFAULT 'No',
    messageDelivered TEXT CHECK (messageDelivered IN ('Yes', 'No')) DEFAULT 'No',
    active INTEGER DEFAULT 1,
    createdAt TEXT DEFAULT CURRENT_TIMESTAMP,
    updatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS notificationUpdated AFTER
UPDATE ON notification FOR EACH ROW BEGIN
UPDATE notification
SET
  updatedAt = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;