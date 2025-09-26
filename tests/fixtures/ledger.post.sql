INSERT INTO accountEntry (
	accountId, 
	referenceNumber, 
	name, 
	description, 
	debitCredit, 
	amount
) VALUES (
	(SELECT id FROM account WHERE accountType = 'ASSET'),
	'1172', 'Some ledger entry', 'Lots of groceries', 'DEBIT', 1234
);
UPDATE account SET balance = COALESCE(balance, 0) + 1234 WHERE id = (SELECT id FROM account WHERE accountType = 'ASSET');
INSERT INTO accountEntry (
	accountId, 
	referenceNumber, 
	name, 
	description, 
	debitCredit, 
	amount
) VALUES (
	(SELECT id FROM account WHERE accountType = 'ASSET'),
	'1172', 'Some ledger entry', 'Lots of groceries', 'CREDIT', 1234
);
UPDATE account SET balance = COALESCE(balance, 0) - 1234 WHERE id = (SELECT id FROM account WHERE accountType = 'ASSET')